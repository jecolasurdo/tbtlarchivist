//! Access to resources on the local filesystem.

use crate::accessors::FromUri;
use std::fs::File;
use std::io::Read;
use thiserror::Error;

/// Provides access to an item on the local filesystem.
pub struct Accessor {}

impl FromUri<Error> for Accessor {
    #[allow(dead_code)]
    fn get(&self, uri: String) -> Result<Vec<u8>, Error> {
        let mut file = File::open(uri)?;
        let mut data = Vec::new();
        file.read_to_end(&mut data)?;
        Ok(data)
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
    #[error("IO error")]
    Io(#[from] std::io::Error),
}
