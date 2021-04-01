//! Access to resources via http.

use crate::accessors::FromUri;
use std::io::Read;
use thiserror::Error;

/// Provides access to an item via http.
pub struct Accessor {}

const HEADER_CONTENT_LENGTH: &str = "Content-Length";

impl FromUri<Error> for Accessor {
    /// Returns the response body of a destinaction http URI. If the request fails, or if the
    /// response does not return 200 for any reason, an error is returned. This method does not
    /// validate the response body, but will ensure that the full body is returned (else an error
    /// will be returned).
    #[allow(dead_code)]
    fn get(&self, uri: String) -> Result<Vec<u8>, Error> {
        let response = ureq::get(&uri).call()?;
        if response.status() != 200 {
            return Err(Error(Box::new(ErrorKind::Non200Response {
                status: response.status(),
            })));
        }

        if !response.has(HEADER_CONTENT_LENGTH) {
            return Err(Error(Box::new(ErrorKind::NoContentLength)));
        }

        let len = response
            .header(HEADER_CONTENT_LENGTH)
            .and_then(|s| s.parse::<usize>().ok())
            .expect("Content-Length was found, but could not be parsed to a usize.");

        let mut bytes = Vec::with_capacity(len);
        response.into_reader().read_to_end(&mut bytes)?;
        Ok(bytes)
    }
}

/// A boxed error resulting from a problem accessing a resource.
#[derive(Error, Debug)]
#[error(transparent)]
pub struct Error(Box<ErrorKind>);

impl<E> From<E> for Error
where
    ErrorKind: From<E>,
{
    fn from(err: E) -> Self {
        Self(Box::new(ErrorKind::from(err)))
    }
}

/// Accessor error variants.
#[allow(missing_docs)]
#[derive(Error, Debug)]
pub enum ErrorKind {
    #[error("Ureq crate error")]
    Ureq(#[from] ureq::Error),
    #[error("IO error")]
    Io(#[from] std::io::Error),

    #[error("Received a non-200 http response {status:?}")]
    Non200Response { status: u16 },

    #[error("No Content-Length header in response")]
    NoContentLength,
}
