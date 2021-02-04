package archivists

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/utils"
)

// A ClipsArchivist looks for clips that have been supplied by an upstream
// clip curator, and places new clips into the collection.
type ClipsArchivist struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartClipsArchivist initializes a clips archivist. The archvist will begin polling the
// supplied queue for new clips, and will place those clips in the supplied
// datastore. The clips archivist operates indefinitely, or until its parent
// context signals that it is done. Once the archivist is initialized, the
// resulting API.Errors and API.Done channels can be monitored. The caller may
// safely exit only when the Errors and Done channels have closed.
func StartClipsArchivist(ctx context.Context, queue messagebus.Receiver, db datastore.DataStorer) *ClipsArchivist {
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

			var clipInfo contracts.ClipInfo
			err = json.Unmarshal(msg.Body, &clipInfo)
			if err != nil {
				errorSource <- err
				err := msg.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- err
				}
				continue
			}

			err = db.UpsertClipInfo(clipInfo)
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

	return &ClipsArchivist{
		Errors: errorSource,
		Done:   done,
	}

}
