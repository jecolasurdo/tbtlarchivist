//! The manager orchestrates high level business logic for the analysis process.

use crate::accessors::FromHTTP;
use crate::engines::Analyzer;
use crate::errors::AnalyzerError;
use contracts::{CompletedResearchItem, PendingResearchItem};

pub struct AnalysisManager<'a> {
    http_accessor: &'a dyn FromHTTP<'a>,
    engine: dyn Analyzer,
}

impl<'a> AnalysisManager<'a> {
    fn run(&'a self, pending_work_item: &'a PendingResearchItem) -> Result<(), AnalyzerError> {
        let mut mp3_data = self
            .http_accessor
            .get(pending_work_item.get_episode().get_media_uri())?;
        let episode_raw = self.engine.mp3_to_raw(mp3_data)?;
        let episode_phash = self.engine.phash(episode_raw)?;
        for clip in pending_work_item.get_clips() {
            mp3_data = self.http_accessor.get(clip.get_media_uri())?;
            let clip_raw = self.engine.mp3_to_raw(mp3_data)?;
            let clip_phash = self.engine.phash(clip_raw)?;
            let offsets = self.engine.find_offsets(clip_raw, episode_raw)?;
            todo!("frame and push to stdout")
        }
        Ok(())
    }
}
