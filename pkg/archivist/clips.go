package archivist

import (
	"context"
)

// A CuratedClipWorker can initialize a stream of inbound curated clips, and
// supplies a method for processing each clip.
type CuratedClipWorker struct{}

// InitializeDataStream opens a stream of inbound clip metadata that needs to
// be processed.
func (c *CuratedClipWorker) InitializeDataStream(ctx context.Context) (<-chan interface{}, error) {
	// clips are unique by name only
	panic("not implemented")
}

// ProcessDatum processes a clip, determing whether or not it should be added
// to the underlaying datastore.
func (c *CuratedClipWorker) ProcessDatum(ctx context.Context, datum interface{}) error {
	panic("not implemented")
}
