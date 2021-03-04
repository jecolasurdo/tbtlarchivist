//! Accessors provide access to external resources.

use thiserror::Error;

/// `FromHTTP` represents something that is able to return an object from an HTTP URI.
pub trait FromHTTP<'a> {
    /// Returns the body of the response object as a byte vector.  If the request does not
    /// specifically return 200 - OK (for any reason) an `AccessorError` must be returned.  It is
    /// not necessarily this method's responsibility to validate the response body.  Consumers of
    /// objects that implement this trait should consider whether or not they need to validate the
    /// resulting byte vector. Implementors of this trait should be expected to ensure, as much as
    /// possible, that the resulting byte vector is a complete (not partial) response.
    fn get(&'a self, uri: &'a str) -> Result<Vec<u8>, AccessorError>;
}

#[derive(Error, Debug)]
pub enum AccessorError {
    #[error("unknown engine error")]
    Unknown,
}
