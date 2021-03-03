//! Accessors provide access to external resources.

use crate::errors::AnalyzerError;

/// `FromHTTP` represents something that is able to return an object from an HTTP URI.
pub trait FromHTTP {
    
    /// Returns the body of the response object as a byte vector.  If the request does not
    /// specifically return 200 - OK (for any reason) an `AnalyzerError` must be returned.  It is
    /// not necessarily this method's responsibility to validate the response body.  Consumers of
    /// objects that implement this trait should consider whether or not they need to validate the
    /// resulting byte vector. Implementors of this trait should be expected to ensure, as much as
    /// possible, that the resulting byte vector is a complete (not partial) response.
    fn get(&self, uri: &str) -> Result<Vec<u8>,AnalyzerError>;
}


// pub fn read_wav(filename: String) -> Result<Vec<i16>, std::io::Error> {
//     let mut file = File::open(filename)?;
//     let (_, bit_depth) = wav::read(&mut file)?;
//     match bit_depth {
//         wav::BitDepth::Sixteen(x) => Ok(x),
//         _ => Err(std::io::Error::new(
//             ErrorKind::Other,
//             "Unexpected bit depth. Only support 16bit",
//         )),
//     }
// }