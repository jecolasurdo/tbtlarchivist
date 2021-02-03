package pendingresearcharchivist

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/datastore"
	"github.com/jecolasurdo/tbtlarchivist/pkg/accessors/messagebus"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

/*
The pending research archivist does the following: When initialized, it
immediately calculates the work surplus. The work surplus is equal to the
number of downstream researchers minus the number of work-items currently on
the pending queue. If this number is negative or zero, the service terminates.
If the number is positive, the one work item is added to the queue and the
process repeats at the work surplus calculation step. Keeping each instance
focused on only adding one work item at a time reduces the overhead of each
instance and allows the process of work-creation to more effectively scale
horizaontally.

Description about what happens down the line after this service publishes the
work-item to the message bus:

	The researcher will consume the work-item, download the episode and begin
	analyzing the clips. As each clip is analyzed, the episode-clip tuple is
	published back to the bus. The completed-work archivist will pick up that
	work, record it in the database, and extend the episode lease (not
	necessarily in that order).
*/

const (
	episodeLeaseDuration = 2 * time.Hour
)

type API struct {
	Errors <-chan error
	Done   <-chan struct{}
}

func Initialize(ctx context.Context, messageBus messagebus.Sender, db datastore.DataStorer) *API {
	errorSource := make(chan error)
	done := make(chan struct{})

	go func() {
		defer close(errorSource)
		defer close(done)

		panic("todo")
		// add ~5 second jitter here to avoid a pileup if multiple archivists
		// start at the same time.

		// then add loop until the overhead is negative or zero.

		queueInfo, err := messageBus.Inspect()
		if err != nil {
			errorSource <- fmt.Errorf("error while inspecting queue: %v", err)
			return
		}

		// If there are no consumers and there are no available pending messages
		// then we queue up one work item in hopes that a consumer will come
		// available soon.
		// If there is more work than consumers, then we do nothing.
		// If there is less work than consumers, then we queue up one work item.
		if !(queueInfo.Consumers == 0 && queueInfo.Messages == 0) {
			overhead := queueInfo.Consumers - queueInfo.Messages
			if overhead <= 0 {
				return
			}
		}

		if contextIsDone(ctx) {
			return
		}

		episode, err := db.GetMostRecentUnleasedEpisode()
		if err != nil {
			errorSource <- fmt.Errorf("error occured while finding unleased episode, %v", err)
			return
		}
		if episode == nil {
			log.Println("No unleased episodes available.")
			return
		}

		if contextIsDone(ctx) {
			return
		}

		// It's technically possible that multiple archivists try to set the
		// episode lease for the same episode at the same time. 1) This should
		// be fairly infrequent, as it would require multiple services to try
		// and take out a lease at the same moment. 2) If it does happen, there
		// is no detriment to the system aside from an inefficient use of
		// resources. 3) The possibility of this happening can be futher
		// reduced by having each service-host initialize the archivists at odd
		// or jittered intervals. This would further reduce the liklihood that
		// multiple archivist instances from competing for leases.
		err = db.SetEpisodeLease(*episode, time.Now().Add(episodeLeaseDuration).UTC())
		if err != nil {
			errorSource <- fmt.Errorf("error setting episode lease: %v\n%v", err, episode)
			return
		}

		if contextIsDone(ctx) {
			return
		}

		clips, err := db.GetUnresearchedClipsForEpisode(*episode)
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

		if contextIsDone(ctx) {
			return
		}

		err = messageBus.Send(messageBytes)
		if err != nil {
			errorSource <- fmt.Errorf("error while trying to push a pendingResearchItem to the message bus. %v %v", pendingResearchItem, err)
			return
		}

	}()

	return &API{
		Errors: errorSource,
		Done:   done,
	}
}

func contextIsDone(ctx context.Context) bool {
	return ctx.Err() != nil
}
