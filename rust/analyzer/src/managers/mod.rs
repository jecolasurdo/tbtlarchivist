//! The manager orchestrates high level business logic for the analysis process.

use crate::accessors::FromURI;
use crate::engines::Analyzer;
use crate::errors::AnalyzerError;
use contracts::{CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::{unbounded, Receiver};
use std::thread;

pub struct AnalysisManager {
    uri_accessor: &'static (dyn FromURI<'static> + Sync),
    engine: dyn Analyzer + Sync,
}

impl AnalysisManager {
    fn run(
        &'static self,
        pending_work_item: &'static PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, AnalyzerError>> {
        let (tx, rx) = unbounded();
        thread::spawn(move || {
            // temporarily unwrapping results. These will be handled
            // more properly in the future by sending errors to `rx`.
            let mut mp3_data = self
                .uri_accessor
                .get(pending_work_item.get_episode().get_media_uri())
                .unwrap();
            let episode_raw = self.engine.mp3_to_raw(mp3_data).unwrap();
            let episode_phash = self.engine.phash(episode_raw).unwrap();
            for clip in pending_work_item.get_clips() {
                mp3_data = self.uri_accessor.get(clip.get_media_uri()).unwrap();
                let clip_raw = self.engine.mp3_to_raw(mp3_data).unwrap();
                let clip_phash = self.engine.phash(clip_raw).unwrap();
                let offsets = self.engine.find_offsets(clip_raw, episode_raw).unwrap();

                let mut cri = CompletedResearchItem::new();
                cri.set_lease_id(pending_work_item.get_lease_id().to_string());
                // finish populating cri

                tx.send(Ok(cri))
                    .expect("unable to transmit completd researech item");
            }
        });
        rx
    }
}
