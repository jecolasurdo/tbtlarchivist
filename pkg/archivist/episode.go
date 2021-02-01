package archivist

import (
	"context"
)

// A CuratedEpisodeWorker can initialize a stream of inbound curated episodes,
// and supplies a method for processing each episode.
type CuratedEpisodeWorker struct{}

// InitializeDataStream opens a stream of inbound episode metadata that needs
// to be processed.
func (c *CuratedEpisodeWorker) InitializeDataStream(ctx context.Context) (<-chan interface{}, error) {
	// episodes are unique by name + date aired
	// check to see if the episode exists
	//	if it does not: add it
	//	else:
	//		check if any of its details of changed
	//		if so, update the details
	//	etc...
	panic("not implemented")
}

// ProcessDatum processes an episode, determing whether or not it should be
// added to the underlaying datastore.
func (c *CuratedEpisodeWorker) ProcessDatum(ctx context.Context, datum interface{}) error {
	panic("not implemented")
}
