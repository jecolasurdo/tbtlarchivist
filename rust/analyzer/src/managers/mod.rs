//! The manager orchestrates high level business logic for the analysis process.

use crate::accessors::FromURI;
use crate::engines::Analyzer;
use crate::errors::AnalyzerError;
use contracts::{ClipInfo, CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::{unbounded, Receiver, Sender};
use protobuf::well_known_types::Timestamp;
use std::convert::TryInto;
use std::thread;
use std::time::{SystemTime, UNIX_EPOCH};

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
    pub fn run(
        &'static self,
        pri: &'static PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, AnalyzerError>> {
        let (tx, rx) = unbounded();
        thread::spawn(move || {
            if let Err(err) = self.process_episode(pri, tx) {
                tx.send(Err(err)).expect("run: unable to transmit error");
            }
        });
        rx
    }

    pub(self) fn process_episode(
        &'static self,
        pri: &'static PendingResearchItem,
        tx: Sender<Result<CompletedResearchItem, AnalyzerError>>,
    ) -> Result<(), AnalyzerError> {
        let mut mp3_data = self.uri_accessor.get(pri.get_episode().get_media_uri())?;
        let episode_raw = self.analyzer.mp3_to_raw(mp3_data)?;
        let episode_phash = self.analyzer.phash(episode_raw)?;
        for clip in pri.get_clips() {
            if let Err(err) = self.process_clip(pri, episode_raw, episode_phash, *clip, tx) {
                // errors at this level do not halt the entire process. Instead
                // we just forward them to the caller.
                tx.send(Err(err))
                    .expect("process_episode: unable to transmit error");
            }
        }
        Ok(())
    }

    pub(self) fn process_clip(
        &'static self,
        pri: &'static PendingResearchItem,
        episode_raw: Vec<i16>,
        episode_phash: Vec<u8>,
        clip: ClipInfo,
        tx: Sender<Result<CompletedResearchItem, AnalyzerError>>,
    ) -> Result<(), AnalyzerError> {
        let mp3_data = self.uri_accessor.get(clip.get_media_uri())?;
        let clip_raw = self.analyzer.mp3_to_raw(mp3_data)?;
        let clip_phash = self.analyzer.phash(clip_raw)?;
        let offsets = self.analyzer.find_offsets(clip_raw, episode_raw)?;

        let mut cri = CompletedResearchItem::new();
        cri.set_research_date(proto_now());
        // finish setting fields
        // ...
        // ,..

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
