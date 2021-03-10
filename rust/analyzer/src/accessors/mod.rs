//! Accessors provide access to external resources.

pub mod http;

/// `FromURI` represents something that is able to return an object from a URI.
pub trait FromURI<'a, E> {
    /// Returns the object as a byte vector.  If the method cannot succeed for
    /// for any reason an `AccessorError` must be returned.  It is not
    /// necessarily this method's responsibility to validate the returned object.
    /// Consumers of objects that implement this trait should consider whether or
    /// not they need to validate the  resulting byte vector.
    fn get(&'a self, uri: &'a str) -> Result<Vec<u8>, E>;
}
