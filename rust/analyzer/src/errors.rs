//! Error types associated with the analyzer.

use crate::accessors::AccessorError;
use crate::engines::EngineError;
use thiserror::Error;

#[derive(Error, Debug)]
#[error("{0}")]
pub enum AnalyzerError {
    Accessor(#[from] AccessorError),
    Engine(#[from] EngineError),
}
