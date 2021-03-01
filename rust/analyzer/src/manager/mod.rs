//! The manager orchestrates high level business logic for the analysis process.

use crate::errors::AnalyzerError;

/// An `HTTPAccessor` represents something that is able to return an object from an HTTP URI.
pub trait HTTPAccessor {
    
    /// Returns the body of the response object as a byte vector.  If the request does not
    /// specifically return 200 - OK (for any reason) an `AnalyzerError` must be returned.  It is
    /// not necessarily this method's responsibility to validate the response body.  Consumers of
    /// objects that implement this trait should consider whether or not they need to validate the
    /// resulting byte vector. Implementors of this trait should be expected to ensure, as much as
    /// possible, that the resulting byte vector is a complete (not partial) response.
    fn get(&self, uri: &str) -> Result<Vec<u8>,AnalyzerError>;
}

/// An `AnalysisEngine` represents something that is able to manipulate and analyze audio data.
pub trait AnalysisEngine {

    /// Takes an byte vector representation of a raw mp3 file, and converts it to raw 16-bit mono
    /// audio data. If the supplied mp3 data cannot be converted to raw-audio for any reason, an
    /// `AnalyzerError` is returned.
    fn mp3_to_raw(&self, mp3: Vec<u8>) -> Result<Vec<i16>, AnalyzerError>;
    

    /// Takes a 16bit raw audio data and calculates its perceptual hash.  If the method is unable
    /// to proceed for any reason, it will return an `AnalyzerError`. 
    fn phash(&self, raw: Vec<i16>) -> Result<Vec<u8>, AnalyzerError>;
    
    /// Searches for any likely occurences of `candidate` within `target` and returns the position
    /// of each occurence in a vector of offsets. Any errors that result during the process of
    /// finding offsets will immediately return an `AnalyzerError`.
    fn find_offsets(&self, candidate: Vec<i16>, target: Vec<i16>) -> Result<Vec<i64>, AnalyzerError>;
}