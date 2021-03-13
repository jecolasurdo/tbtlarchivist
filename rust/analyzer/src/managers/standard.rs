//! contains the typical concrete Runner implementation.

use crate::{accessors::FromURI, engines::Analyzer};
use cancel::Token;
use contracts::{CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::{unbounded, Receiver, Sender};
use crossbeam_utils::thread;
use protobuf::well_known_types::Timestamp;
use std::error::Error;
use std::marker::PhantomData;
use std::{
    convert::TryInto,
    time::{SystemTime, UNIX_EPOCH},
};

const RAW_SAMPLE_RATE: usize = 44_100;

// Basis of 1e9 represents a nanosecond duration resolution,
const RAW_DURATION_BASIS: usize = 1_000_000_000;

/// Returns a new `AnalysisManager`.
pub fn new<A, U, AE, UE, E>(analyzer_engine: A, uri_accessor: U) -> AnalysisManager<A, U, AE, UE, E>
where
    A: Analyzer<AE> + Send + Sync,
    U: FromURI<UE> + Send + Sync,
    AE: Error + Send + Sync,
    UE: Error + Send + Sync,
    E: From<AE> + From<UE> + Send + Sync,
{
    AnalysisManager {
        analyzer_engine,
        uri_accessor,
        _phantom_ae: PhantomData,
        _phantom_ue: PhantomData,
        _phantom_e: PhantomData,
    }
}

/// An `AnalysisManager` orchestrates the process conducing the analysis prescribed
/// by a `PendingResearchItem`.
pub struct AnalysisManager<A, U, AE, UE, E>
where
    A: Analyzer<AE> + Send + Sync,
    U: FromURI<UE> + Send + Sync,
    AE: Error + Send + Sync,
    UE: Error + Send + Sync,
    E: From<AE> + From<UE> + Send + Sync,
{
    analyzer_engine: A,
    uri_accessor: U,
    _phantom_ae: PhantomData<AE>,
    _phantom_ue: PhantomData<UE>,
    _phantom_e: PhantomData<E>,
}

impl<A, U, AE, UE, E> AnalysisManager<A, U, AE, UE, E>
where
    A: Analyzer<AE> + Send + Sync,
    U: FromURI<UE> + Send + Sync,
    AE: Error + Send + Sync,
    UE: Error + Send + Sync,
    E: From<AE> + From<UE> + Send + Sync,
{
    /// Starts the analysis process, returning a channel on which completed
    /// research and/or errors are transmitted. This channel must be polled
    /// until it is closed. To cleanly interupt and halt the operation of a
    /// running analysis, a cancellation should be broadcast via the `ctx`
    /// object.
    pub fn run(
        &self,
        ctx: &Token,
        pri: &PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, E>> {
        let (tx, rx) = unbounded();

        thread::scope(|s| {
            s.spawn(|_| {
                if let Err(err) = self.process_episode(ctx, pri, &tx) {
                    tx.send(Err(err)).expect("run: Unable to transmit error.");
                }
            });
        })
        .expect("run: An unrecoverable error occured.");
        rx
    }

    fn process_episode(
        &self,
        ctx: &Token,
        pri: &contracts::PendingResearchItem,
        tx: &Sender<Result<CompletedResearchItem, E>>,
    ) -> Result<(), E> {
        let mp3_data = self
            .uri_accessor
            .get(pri.get_episode().get_media_uri().to_string())?;
        let episode_raw = self.analyzer_engine.mp3_to_raw(&mp3_data)?;
        let episode_phash = self.analyzer_engine.phash(&episode_raw)?;
        for clip in pri.get_clips() {
            if ctx.is_canceled() {
                break;
            }
            if let Err(err) = self.process_clip(pri, &episode_raw, &episode_phash, clip, tx) {
                // errors at this level do not halt the entire process. Instead
                // we just forward them to the caller. The caller may decide to
                // broadcast a cancellation if the error rates are out of hand,
                // at which point this method would expect ctx.is_cancelled()
                // to return true.
                tx.send(Err(err))
                    .expect("process_episode: Unable to transmit error.");
            }
        }
        Ok(())
    }

    fn process_clip(
        &self,
        pri: &PendingResearchItem,
        episode_raw: &[i16],
        episode_phash: &[u8],
        clip: &contracts::ClipInfo,
        tx: &Sender<Result<CompletedResearchItem, E>>,
    ) -> Result<(), E> {
        let mp3_data = self.uri_accessor.get(clip.get_media_uri().to_string())?;
        let clip_raw = self.analyzer_engine.mp3_to_raw(&mp3_data)?;
        let clip_phash = self.analyzer_engine.phash(&clip_raw)?;
        let offsets = self.analyzer_engine.find_offsets(&clip_raw, episode_raw)?;

        let mut cri = CompletedResearchItem::new();
        cri.set_research_date(proto_now());
        cri.set_episode_info(pri.get_episode().clone());
        cri.set_clip_info(clip.clone());
        cri.set_episode_duration(duration(episode_raw.len()));
        cri.set_episode_hash(episode_phash.to_vec());
        cri.set_clip_duration(duration(clip_raw.len()));
        cri.set_clip_hash(clip_phash);
        cri.set_clip_offsets(offsets);
        cri.set_lease_id(pri.get_lease_id().to_string());
        cri.set_revoke_lease(false);

        tx.send(Ok(cri))
            .expect("process_clip: Unable to transmit completed-work-item.");
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

#[allow(clippy::as_conversions)]
fn duration(samples: usize) -> i64 {
    (samples / RAW_SAMPLE_RATE * RAW_DURATION_BASIS) as i64
}
