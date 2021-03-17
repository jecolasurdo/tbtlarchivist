//! Accessors provide access to external resources.

use std::error::Error;

pub mod file;
pub mod http;

/// `FromURI` represents something that is able to return an object from a URI.
pub trait FromURI<E>
where
    E: Error + Send + Sync,
{
    /// Returns the object as a byte vector.  If the method cannot succeed for
    /// for any reason an `AccessorError` must be returned.  It is not
    /// necessarily this method's responsibility to validate the returned object.
    /// Consumers of objects that implement this trait should consider whether or
    /// not they need to validate the  resulting byte vector.
    fn get(&self, uri: String) -> Result<Vec<u8>, E>;
}
