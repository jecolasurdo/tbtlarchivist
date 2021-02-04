package archivists

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// An EpisodesArchivist looks for episodes that have been supplied by an
// upstream episode curator, and places new episodes into the collection.
type EpisodesArchivist struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartEpisodesArchivist initializes an episode archivist. The archvist will begin polling
// the supplied queue for new episodes, and will place those episodes in the
// supplied datastore. The archivist operates indefinitely, or until its parent
// context signals that it is done. Once the archivist is initialized, the
// resulting API.Errors and API.Done channels can be monitored. The caller may
// safely exit only when the Errors and Done channels have closed.
func StartEpisodesArchivist(ctx context.Context, queue messagebus.Receiver, db datastore.DataStorer) *EpisodesArchivist {
	errorSource := make(chan error)
	done := make(chan struct{})
	go func() {
		defer close(errorSource)
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			msg := queue.Receive()
			if msg == nil {
				runtime.Gosched()
				continue
			}

			var episodeInfo contracts.EpisodeInfo
			err := json.Unmarshal(msg.Body, &episodeInfo)
			if err != nil {
				errorSource <- err
				err := msg.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- err
				}
				continue
			}

			err = db.UpsertEpisodeInfo(episodeInfo)
			if err != nil {
				errorSource <- err
				err := msg.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- err
				}
				continue
			}

			err = msg.Acknowledger.Ack()
			if err != nil {
				errorSource <- err
			}
		}
	}()

	return &EpisodesArchivist{
		Errors: errorSource,
		Done:   done,
	}

}
