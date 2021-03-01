//! Error types associated with the analyzer.

#[derive(Debug, Clone, PartialEq)]
/// A general error that has occurred during an analysis operation.
pub struct AnalyzerError {
    msg: String,
}

impl<'a> AnalyzerError {
    /// Instantiates a new `AnalyzerError` with a message.
    pub fn new(msg: String) -> Self {
        Self { msg }
    }

    /// A message associated with this error.
    pub fn message(&self) -> String {
        self.msg.clone()
    }
}