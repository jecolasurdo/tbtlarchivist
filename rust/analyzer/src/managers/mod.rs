//! The manager orchestrates high level business logic for the analysis process.

use crate::accessors::FromURI;
use crate::engines::Analyzer;
use crate::errors::AnalyzerError;
use contracts::{ClipInfo, CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::{unbounded, Receiver, Sender};
use std::thread;

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
    fn run(
        &'static self,
        pri: &'static PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, AnalyzerError>> {
        let (tx, rx) = unbounded();
        thread::spawn(move || {
            if let Err(err) = self.process_episode(pri, tx) {
                tx.send(Err(err)).unwrap();
            }
        });
        rx
    }

    fn process_episode(
        &'static self,
        pri: &'static PendingResearchItem,
        tx: Sender<Result<CompletedResearchItem, AnalyzerError>>,
    ) -> Result<(), AnalyzerError> {
        let mut mp3_data = self.uri_accessor.get(pri.get_episode().get_media_uri())?;
        let episode_raw = self.analyzer.mp3_to_raw(mp3_data)?;
        let episode_phash = self.analyzer.phash(episode_raw)?;
        for clip in pri.get_clips() {
            if let Err(err) = self.process_clip(pri, episode_raw, *clip, tx) {
                tx.send(Err(err)).unwrap();
            }
        }
        Ok(())
    }

    fn process_clip(
        &'static self,
        pri: &'static PendingResearchItem,
        episode_raw: Vec<i16>,
        clip: ClipInfo,
        tx: Sender<Result<CompletedResearchItem, AnalyzerError>>,
    ) -> Result<(), AnalyzerError> {
        let mp3_data = self.uri_accessor.get(clip.get_media_uri())?;
        let clip_raw = self.analyzer.mp3_to_raw(mp3_data)?;
        let clip_phash = self.analyzer.phash(clip_raw)?;
        let offsets = self.analyzer.find_offsets(clip_raw, episode_raw)?;

        let mut cri = CompletedResearchItem::new();
        cri.set_lease_id(pri.get_lease_id().to_string());
        // finish populating cri

        tx.send(Ok(cri)).unwrap();

        Ok(())
    }
}
