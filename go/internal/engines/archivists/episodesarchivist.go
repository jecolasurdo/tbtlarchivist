package archivists

import (
	"context"
	"runtime"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/utils"
	"google.golang.org/protobuf/proto"
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
			if utils.ContextIsDone(ctx) {
				return
			}

			// If we're in a position where we're getting a lot of errors or
			// nil messages from the queue, we can end up hogging resources
			// from other goroutines. So, we yield to get out of their way.
			// Though the runtime technically can yield on any function call,
			// it will only do so on non-inlined calls. Since we don't know for
			// sure if the next call is inlined, we explicitly yield to be
			// safe.
			runtime.Gosched()

			msg, err := queue.Receive()
			if err != nil {
				errorSource <- err
				continue
			}

			if msg == nil {
				continue
			}

			episodeInfo := new(contracts.EpisodeInfo)
			err = proto.Unmarshal(msg.Body, episodeInfo)
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
