//! Engines implement domain-specific logic, such as the calculation of perceptual hashes, and
//! audio analysis.

pub mod cosine_similarity;

use thiserror::Error;

/// An `Analyzer` represents something that is able to manipulate and analyze audio data.
pub trait Analyzer {
    /// Takes a byte vector representation of a raw mp3 file, and converts it to raw 16-bit mono
    /// audio data. If the supplied mp3 data cannot be converted to raw-audio for any reason, an
    /// `EngineError` is returned.
    fn mp3_to_raw(&self, mp3: Vec<u8>) -> Result<Vec<i16>, EngineError>;

    /// Takes 16bit raw audio data and calculates its perceptual hash.  If the method is unable
    /// to proceed for any reason, it will return an `EngineError`.
    fn phash(&self, raw: Vec<i16>) -> Result<Vec<u8>, EngineError>;
    /// Searches for any likely occurences of `candidate` within `target` and returns the position
    /// of each occurence as a vector of offsets. Any errors that result during the process of
    /// finding offsets will immediately return an `EngineError`.
    fn find_offsets(&self, candidate: Vec<i16>, target: Vec<i16>) -> Result<Vec<i64>, EngineError>;
}

#[derive(Error, Debug)]
pub enum EngineError {
    #[error("unknown engine error")]
    Unknown,
}
