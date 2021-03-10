//! Managers orchestrate high level business logic for the analysis process.

pub mod standard;

use crate::accessors::FromURI;
use crate::engines::Analyzer;
use cancel::Token;
use contracts::{CompletedResearchItem, PendingResearchItem};
use crossbeam_channel::Receiver;

/// A `Runner` is responsible for ochestrating high level analyzer logic.
pub trait Runner<A, U, E>
where
    A: Analyzer<E> + Sync,
    U: FromURI<'static, E> + Sync,
    E: Send + Sync,
{
    /// Starts the analysis process.
    fn run(
        &'static self,
        ctx: &'static Token,
        pri: &'static PendingResearchItem,
    ) -> Receiver<Result<CompletedResearchItem, E>>;
}
