package archivist

import (
	"context"
)

// A CompletedResearchWorker can initialize a stream of inbound completed
// research data, and supplies a method for processing each datum.
type CompletedResearchWorker struct{}

// InitializeDataStream opens a stream of inbound research results that need to
// be processed.
func (c *CompletedResearchWorker) InitializeDataStream(ctx context.Context) (<-chan interface{}, error) {
	panic("not implemented")
}

// ProcessDatum processes a research result.
func (c *CompletedResearchWorker) ProcessDatum(ctx context.Context, datum interface{}) error {
	panic("not implemented")
}
