//! Error types associated with the analyzer.

use crate::accessors::AccessorError;
use crate::engines::EngineError;
use thiserror::Error;

#[derive(Error, Debug)]
#[error("{0}")]
/// An error that might occur while an Analyzer is running.
pub enum AnalyzerError {
    /// An error returned by an accessor.
    Accessor(#[from] AccessorError),
    /// An error returned by an engine.
    Engine(#[from] EngineError),
}
