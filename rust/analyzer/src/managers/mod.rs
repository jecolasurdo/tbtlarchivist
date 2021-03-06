//! The manager orchestrates high level business logic for the analysis process.

use crate::accessors::FromURI;
use crate::engines::Analyzer;
use crate::errors::AnalyzerError;
use cancel::Token;
use contracts::{ClipInfo, CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::{unbounded, Receiver, Sender};
use protobuf::well_known_types::Timestamp;
use std::convert::TryInto;
use std::thread;
use std::time::{SystemTime, UNIX_EPOCH};

/// An `AnalysisManager` orchestrates the process conducing the analysis prescribed
/// by a `PendingResearchItem`.
pub struct AnalysisManager<A, U>
where
    A: Analyzer + Sync,
    U: FromURI<'static> + Sync,
{
    analyzer: A,
    uri_accessor: U,
}

impl<A, U> AnalysisManager<A, U>
where
    A: Analyzer + Sync,
    U: FromURI<'static> + Sync,
{
    /// Starts the analysis process, returning a channel on which completed
    /// research and/or errors are transmitted. This channel must be polled
    /// until it is closed. To cleanly interupt and halt the operation of a
    /// running analysis, a cancellation should be broadcast via the `ctx`
    /// object.
    pub fn run(
        &'static self,
        ctx: &'static Token,
        pri: &'static PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, AnalyzerError>> {
        let (tx, rx) = unbounded();
        thread::spawn(move || {
            if let Err(err) = self.process_episode(ctx, pri, &tx) {
                tx.send(Err(err)).expect("run: unable to transmit error");
            }
        });
        rx
    }

    pub(self) fn process_episode(
        &'static self,
        ctx: &'static Token,
        pri: &'static PendingResearchItem,
        tx: &Sender<Result<CompletedResearchItem, AnalyzerError>>,
    ) -> Result<(), AnalyzerError> {
        let mp3_data = self.uri_accessor.get(pri.get_episode().get_media_uri())?;
        let episode_raw = self.analyzer.mp3_to_raw(&mp3_data)?;
        let episode_phash = self.analyzer.phash(&episode_raw)?;
        for clip in pri.get_clips() {
            if ctx.is_canceled() {
                break;
            }
            if let Err(err) = self.process_clip(pri, &episode_raw, &episode_phash, clip, tx) {
                // errors at this level do not halt the entire process. Instead we just forward
                // them to the caller. The caller may decide to broadcast a cancellation if the
                // error rates are out of hand, at which point this method would expect
                // ctx.is_cancelled() to return true.
                tx.send(Err(err))
                    .expect("process_episode: unable to transmit error");
            }
        }
        Ok(())
    }

    pub(self) fn process_clip(
        &'static self,
        pri: &'static PendingResearchItem,
        episode_raw: &[i16],
        episode_phash: &[u8],
        clip: &'static ClipInfo,
        tx: &Sender<Result<CompletedResearchItem, AnalyzerError>>,
    ) -> Result<(), AnalyzerError> {
        let mp3_data = self.uri_accessor.get(clip.get_media_uri())?;
        let clip_raw = self.analyzer.mp3_to_raw(&mp3_data)?;
        let clip_phash = self.analyzer.phash(&clip_raw)?;
        let offsets = self.analyzer.find_offsets(&clip_raw, episode_raw)?;

        let mut cri = CompletedResearchItem::new();
        cri.set_research_date(proto_now());
        cri.set_episode_info(pri.get_episode().clone());
        cri.set_clip_info(clip.clone());
        cri.set_episode_duration(0);
        cri.set_episode_hash(episode_phash.to_vec());
        cri.set_clip_duration(0);
        cri.set_clip_hash(clip_phash);
        cri.set_clip_offsets(offsets);
        cri.set_lease_id(pri.get_lease_id().to_string());
        cri.set_revoke_lease(false);

        tx.send(Ok(cri))
            .expect("process_clip: unable to transmit completed work item");
        Ok(())
    }
}

fn proto_now() -> Timestamp {
    let n = SystemTime::now().duration_since(UNIX_EPOCH).unwrap();
    let mut t = Timestamp::new();
    t.set_seconds(n.as_secs().try_into().unwrap());
    t.set_nanos(n.subsec_nanos().try_into().unwrap());
    t
}
