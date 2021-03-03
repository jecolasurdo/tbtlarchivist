//! The manager orchestrates high level business logic for the analysis process.

use crate::accessors::FromHTTP;
use crate::engines::Analyzer;
use crate::errors::AnalyzerError;
use contracts::{PendingResearchItem, CompletedResearchItem};


pub struct AnalysisManager<'a>{
    http_accessor: &'a dyn FromHTTP<'a>,
    engine: dyn Analyzer
}

impl<'a> AnalysisManager<'a> {
    fn run(&'a self, pending_work_item: &'a PendingResearchItem) -> Result<CompletedResearchItem, AnalyzerError> {
        let mut mp3_data =  self.http_accessor.get(pending_work_item.get_episode().get_media_uri())?; 
        let episode_raw = self.engine.mp3_to_raw(mp3_data)?;
        let episode_phash = self.engine.phash(episode_raw)?;
        for clip in pending_work_item.get_clips() {
            mp3_data = self.http_accessor.get(clip.get_media_uri())?;
            let clip_raw = self.engine.mp3_to_raw(mp3_data)?;
            let clip_phash = self.engine.phash(clip_raw)?;
            let offsets = self.engine.find_offsets(clip_raw, episode_raw)?;
            
            // Not clear on whose responsibility it is to frame a completed work item and send it
            // to stdout
            //
            // At the moment I feel like the most straight forward approach would be to have the
            // manager return a crossbeam channel with results that the caller can poll and process
            // however it wants.
            //
            // That does add some responsibility to the host, but I don't think the manager should
            // care about stdout or serialization protocol.
            //
            // Need to ponder this a little.
            todo!("push the offset to a channel that is being polled?")
        }
        
    }

}
