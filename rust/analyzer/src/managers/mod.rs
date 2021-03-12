//! Managers orchestrate high level business logic for the analysis process.

pub mod standard;

use crate::accessors::FromURI;
use crate::engines::Analyzer;
use cancel::Token;
use contracts::{CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::Receiver;
use std::error::Error;

/// A `Runner` is responsible for ochestrating high level analyzer logic.
pub trait Runner<'a, A, U, E>
where
    A: Analyzer<E> + Send + Sync,
    U: FromURI<'a, E> + Send + Sync,
    E: Error + Send + Sync,
{
    /// Starts the analysis process.
    fn run(
        &'a self,
        ctx: &'a Token,
        pri: &'a PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, E>>;
}
