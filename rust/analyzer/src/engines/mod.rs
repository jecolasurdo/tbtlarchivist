//! Engines implement domain-specific logic, such as the calculation of perceptual hashes, and
//! audio analysis.

use std::error::Error;

pub mod cosine_similarity;

/// A container for raw decoded audio data and metadata associated with the
/// audio.
pub struct Raw {
    /// The raw monaural audio samples.
    pub data: Vec<i16>,
    /// The duration of the audio in `data` expressed in nanoseconds.
    pub duration_ns: i64,
}

/// An `Analyzer` represents something that is able to manipulate and analyze audio data.
pub trait Analyzer<E>
where
    E: Error + Send + Sync,
{
    /// Takes a byte vector representation of a raw mp3 file, and converts it to raw 16-bit mono
    /// audio data.
    fn mp3_to_raw(&self, mp3: &[u8]) -> Result<Raw, E>;
    /// Takes 16bit raw audio data and calculates its perceptual hash.
    fn fingerprint(&self, raw: &[i16]) -> Result<String, E>;
    /// Searches for any likely occurences of `candidate` within `target` and returns the position
    /// of each occurence as a vector of nanosecond offsets.
    fn find_offsets(&self, candidate: &[i16], target: &[i16]) -> Result<Vec<i64>, E>;
}
