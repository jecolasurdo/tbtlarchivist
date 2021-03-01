//! The analyzer provides means for identifying if one audio clip exists within another.

#![warn(
    missing_docs,
    missing_doc_code_examples,
    broken_intra_doc_links,
    clippy::all,
    clippy::pedantic,
    clippy::nursery,
    clippy::as_conversions,
    clippy::todo,
    clippy::print_stdout,
    clippy::use_debug
)]
#![allow(
    clippy::must_use_candidate,
    clippy::float_cmp,
    clippy::similar_names,
    clippy::missing_errors_doc,
    clippy::missing_const_for_fn
)]

pub mod manager;
pub(crate) mod accessor;
pub(crate) mod engine;
pub mod errors;