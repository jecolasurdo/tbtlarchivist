//! Engines implement domain-specific logic, such as the calculation of perceptual hashes, and
//! audio analysis.

use std::error::Error;

pub mod cosine_similarity;

/// An `Analyzer` represents something that is able to manipulate and analyze audio data.
pub trait Analyzer<E>
where
    E: Error + Send + Sync,
{
    /// Takes a byte vector representation of a raw mp3 file, and converts it to raw 16-bit mono
    /// audio data.
    fn mp3_to_raw(&self, mp3: &[u8]) -> Result<Vec<i16>, E>;
    /// Takes 16bit raw audio data and calculates its perceptual hash.
    fn phash(&self, raw: &[i16]) -> Result<Vec<u8>, E>;
    /// Searches for any likely occurences of `candidate` within `target` and returns the position
    /// of each occurence as a vector of offsets.
    fn find_offsets(&self, candidate: &[i16], target: &[i16]) -> Result<Vec<i64>, E>;
}
