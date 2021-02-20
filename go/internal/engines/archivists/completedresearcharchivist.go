package archivists

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/google/uuid"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"github.com/jecolasurdo/tbtlarchivist/go/internal/utils"
	"google.golang.org/protobuf/proto"
)

// A CompletedResearchArchivist determines if any upstream researchers have
// reported any completed work, and, if so, records thwat work in the datastore
// and renews the lease on the associated episode.
type CompletedResearchArchivist struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartCompletedResearchArchivist starts the archivist, which begins polling
// for completed work. When completed work is found, it is recorded in the
// datastore and the associated episode's lease is renewed. This archivist will
// continue to poll for completed work until it determines that no work is
// available, at which point the archivist will exit and no further work will
// be done. Thus, it is the responsibility of the host system to periodically
// start an archivist via a cron job or some other desired scheduler.
func StartCompletedResearchArchivist(ctx context.Context, messageBus messagebus.SenderReceiver, db datastore.DataStorer) *CompletedResearchArchivist {
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

			queueInfo, err := messageBus.Inspect()
			if err != nil {
				errorSource <- fmt.Errorf("an error occured while inspecting the completed research queue %v", err)
				continue
			}

			if queueInfo.Messages == 0 {
				log.Println("No completed work to process at the moment.")
				return
			}

			rawMessage, err := messageBus.Receive()
			if err != nil {
				errorSource <- fmt.Errorf("an error occured while trying to consume a message from the completed research queue %v", err)
				continue
			}

			if len(rawMessage.Body) == 0 {
				continue
			}

			var completedResearchItem *contracts.CompletedResearchItem
			err = proto.Unmarshal(rawMessage.Body, completedResearchItem)
			if err != nil {
				errorSource <- fmt.Errorf("an error occured while unmarshalling a completed research item. %v %v", rawMessage.Body, err)
				err = rawMessage.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- fmt.Errorf("an error occured while trying to send a negative achnowledgement to the message bus %v", err)
				}
				continue
			}

			var operationType string
			if completedResearchItem.RevokeLease {
				err = db.RevokeResearchLease(uuid.MustParse(completedResearchItem.LeaseId))
				operationType = "revoke"
			} else {
				err = db.RenewResearchLease(uuid.MustParse(completedResearchItem.LeaseId), time.Now().Add(episodeLeaseDuration).UTC())
				operationType = "renew"
			}

			if err != nil {
				errorSource <- fmt.Errorf("an error occured trying to %v a lease. %v %v", operationType, rawMessage.Body, err)
				err = rawMessage.Acknowledger.Nack(false)
				if err != nil {
					errorSource <- fmt.Errorf("an error occured while trying to send a negative achnowledgement to the message bus %v", err)
				}
				continue
			}

			// A number of actions (such as updating hashes, choosing whether
			// or not to insert a new completd research item, etc.) are
			// deferred to the data layer so the operations can be done
			// atomically without having to expose transaction awareness to the
			// archivist. This does bleed some business logic to the data
			// layer, so be careful if refactoring.
			err = db.RecordCompletedResearch(completedResearchItem)
			if err != nil {
				errorSource <- fmt.Errorf("an error occured recording completed research to the datastore. %v %v", rawMessage.Body, err)
				err = rawMessage.Acknowledger.Nack(true)
				if err != nil {
					errorSource <- fmt.Errorf("an error occured while trying to send a negative achnowledgement to the message bus %v", err)
				}
				continue
			}

			err = rawMessage.Acknowledger.Ack()
			if err != nil {
				errorSource <- fmt.Errorf("an error occured while trying to acknowledge receipt of a message %v", err)
				continue
			}
		}
	}()

	return &CompletedResearchArchivist{
		Errors: errorSource,
		Done:   done,
	}
}
