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
)

// A PendingResearchArchivist determines if any research work should be done,
// and, if so, produces a pending-work-item for a downstream researcher to act
// upon.
type PendingResearchArchivist struct {
	Errors <-chan error
	Done   <-chan struct{}
}

// StartPendingResearchArchivist starts the archivist, and begins the following
// process:
// 0) The archivist pauses for a random short interval to ensure its start time
// is unlikely to be the same as some other parallel archivist instances.
// 1) Available new-research overhead is calculated. Overhead is equal to the
// number of known downstream researchers minus the number of work-items
// currently on the pending queue.
// 2) If the available overhead is negative or zero, the archivist exits, and
// no further steps are taken.
// 3) If the number is positive (or if there are both no active researchers and
// no queued work) then one work item is added to the queue.
// 4) The process repeats at step 0.
//
// An archivist's host should expect the archivist to exit when the archivist
// has determined that no overhead is available to queue more work.  It is the
// host's responsibility to initialize the archivist periodically to check if
// work needs to be queued.  This can be done via an automated cron job or
// other scheduler as desired.
//
// When queuing work, it's technically possible that multiple archivists try to
// set the episode lease for the same episode at the same time. 1) This should
// be fairly infrequent, as it would require multiple services to try and take
// out the same lease at the same moment. 2) If it does happen, there is no
// detriment to the system aside from an inefficient use of resources. 3) The
// possibility of this happening is futher reduced by having each archivist
// start work at jittered intervals. Thus, if two archivists are initialized at
// nearly the same moment, the jitter will reduce the likelihood that they try
// to access the database at the same time.
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
				errorSource <- fmt.Errorf("error occured while finding unleased episode, %v", err)
				return
			}
			if episode == nil {
				log.Println("No unleased episodes available.")
				return
			}

			err = db.SetResearchLease(*episode, time.Now().Add(episodeLeaseDuration).UTC())
			if err != nil {
				errorSource <- fmt.Errorf("error setting episode lease: %v\n%v", err, episode)
				return
			}

			clips, err := db.GetHighestPriorityClipsForEpisode(*episode)
			if err != nil {
				errorSource <- fmt.Errorf("error retrieving unresearched clips for episode: %v\n%v", err, episode)
				return
			}
			if len(clips) == 0 {
				log.Println("No unresearched clips for this episode")
				return
			}

			pendingResearchItem := contracts.PendingResearchItem{
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
