package archivists

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
	"github.com/jecolasurdo/tbtlarchivist/pkg/utils"
)

const (
	episodeLeaseDuration = 2 * time.Hour
	lowerPacingBound     = 0.0
	upperPacingBound     = 2000.0
	pacingBasis          = time.Millisecond
	clipLimit            = 100
)

// A PendingResearchArchivist determines if any research work should be done,
// and, if so, produces a pending-work-item for a downstream researcher to act
// upon.
type PendingResearchArchivist struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartPendingResearchArchivist starts the archivist, which attempts to create
// pending work-items for downstream researchers to consume.
//
// An archivist's host should expect the archivist to exit when the archivist
// has determined that no overhead is available to queue more work.  It is the
// host's responsibility to initialize the archivist periodically to check if
// work needs to be queued.  This can be done via an automated cron job or
// other scheduler as desired.
func StartPendingResearchArchivist(ctx context.Context, messageBus messagebus.Sender, db datastore.DataStorer) *PendingResearchArchivist {
	errorSource := make(chan error)
	done := make(chan struct{})

	go func() {
		defer close(errorSource)
		defer close(done)

		pace := utils.SetUniformPace(lowerPacingBound, upperPacingBound, pacingBasis)
		for {
			if utils.ContextIsDone(ctx) {
				break
			}

			pace.Wait()

			queueInfo, err := messageBus.Inspect()
			if err != nil {
				errorSource <- fmt.Errorf("error while inspecting queue: %v", err)
				return
			}

			if !(queueInfo.Consumers == 0 && queueInfo.Messages == 0) {
				overhead := queueInfo.Consumers - queueInfo.Messages
				if overhead <= 0 {
					break
				}
			}

			episode, err := db.GetHighestPriorityEpisode()
			if err != nil {
				errorSource <- fmt.Errorf("error occured finding highest priority episode, %v", err)
				return
			}
			if episode == nil {
				log.Println("No episodes available to assign for research.")
				return
			}

			clips, err := db.GetHighestPriorityClipsForEpisode(*episode, clipLimit)
			if err != nil {
				errorSource <- fmt.Errorf("error retrieving clips for episode: %v\n%v", err, episode)
				return
			}
			if len(clips) == 0 {
				log.Println("No clips available to assign for research for this episode.")
				return
			}

			leaseID, err := db.CreateResearchLease(*episode, clips, time.Now().Add(episodeLeaseDuration).UTC())
			if err != nil {
				errorSource <- fmt.Errorf("error creating lease: %v\n%v", err, episode)
				return
			}

			pendingResearchItem := contracts.PendingResearchItem{
				LeaseID: leaseID,
				Episode: *episode,
				Clips:   clips,
			}
			messageBytes, err := json.MarshalIndent(pendingResearchItem, "", "  ")
			if err != nil {
				errorSource <- fmt.Errorf("error marshalling pendingResearchItem to json. %v %v", pendingResearchItem, err)
				return
			}

			err = messageBus.Send(messageBytes)
			if err != nil {
				errorSource <- fmt.Errorf("error while trying to push a pendingResearchItem to the message bus. %v %v", pendingResearchItem, err)
				return
			}
		}
	}()

	return &PendingResearchArchivist{
		Errors: errorSource,
		Done:   done,
	}
}
