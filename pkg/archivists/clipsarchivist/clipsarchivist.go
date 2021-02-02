package clipsarchivist

import (
	"context"
	"encoding/json"
	"runtime"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// API provides information about the state of a clipsarchivist.
type API struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartWork initializes a clips archivist. The archvist will begin polling the
// supplied queue for new clips, and will place those clips in the supplied
// datastore. The clips archivist operates indefinitely, or until its parent
// context signals that it is done. Once the archivist is initialized, the
// resulting API.Errors and API.Done channels can be monitored. The caller may
// safely exit only when the Errors and Done channels have closed.
func StartWork(ctx context.Context, queue messagebus.Receiver, db datastore.DataStorer) *API {
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

			var clipInfo contracts.ClipInfo
			err := json.Unmarshal(msg.Body, &clipInfo)
			if err != nil {
				errorSource <- err
				err := msg.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- err
				}
			}

			err = db.UpsertClipInfo(clipInfo)
			if err != nil {
				errorSource <- err
				err := msg.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- err
				}
			}

			err = msg.Acknowledger.Ack()
			if err != nil {
				errorSource <- err
			}
		}
	}()

	return &API{
		Errors: errorSource,
		Done:   done,
	}

}