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
        match self.http_accessor.get(pending_work_item.get_episode().get_media_uri()) {
            Ok(mp3_data) => {
                match self.engine.mp3_to_raw(mp3_data) {
                    Ok(raw) => {
                        match self.engine.phash(raw) {
                            Ok(episode_phash) => {
                                for clip in pending_work_item.get_clips() {
                                    todo!("ok, this is why rust has a more fluent way of handling errors... need to implement that")        
                                }
                            },
                            Err(_) => {
                                todo!()
                            }
                        }
                    },
                    Err(_) => {
                        todo!()
                    }
                }
            },
            Err(_) => {
                todo!()
            }
        }
    }

}
