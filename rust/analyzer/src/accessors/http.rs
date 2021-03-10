use crate::accessors::FromURI;
use std::io::Read;
use thiserror::Error;

pub struct Accessor {}

const HEADER_CONTENT_LENGTH: &str = "Content-Length";

impl<'a> FromURI<'a, AccessorError> for Accessor {
    /// Returns the response body of a destinaction http URI. If the request fails, or if the
    /// response does not return 200 for any reason, an error is returned. This method does not
    /// validate the response body, but will ensure that the full body is returned (else an error
    /// will be returned).
    #[allow(dead_code)]
    fn get(&'a self, uri: &'a str) -> Result<Vec<u8>, AccessorError> {
        let response = ureq::get(uri).call()?;
        if response.status() != 200 {
            return Err(AccessorError::Non200Response {
                status: response.status(),
            });
        }

        if !response.has(HEADER_CONTENT_LENGTH) {
            return Err(AccessorError::NoContentLength);
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

#[derive(Error, Debug)]
pub enum AccessorError {
    #[error("Ureq crate error")]
    Ureq(#[from] ureq::Error),
    #[error("IO error")]
    Io(#[from] std::io::Error),

    #[error("Received a non-200 http response {status:?}")]
    Non200Response { status: u16 },

    #[error("No Content-Length header in response")]
    NoContentLength,
}
