package archivist

import (
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// API is an instance of an archivist. An archivist is responsible for the
// following tasks:
//  - Recording episode and clip metadata as reported from the curators.
//  - Creating work for researchers.
//  - Checking in work returned from the researchers.
type API struct {
	errorSource <-chan error
}

type worker struct {
	source   <-chan interface{}
	delegate func(interface{}) error
}

// Initialize activates a set of workers. If there is an error activating any
// of the workers, this method will immediately return an error. Each worker
// produces a stream of pending-work-items. The places all pending-work-items
// from all workers into a common queue, and polls this queue.  As each item is
// dequeued, the item is executed by the worker's delegate function. If the
// delegate function returns an error, that error is passed into the APIs error
// channel.
func Initialize() (*API, error) {
	
	panic("TODO: Finish documenting the theory of operation for this method and
	figure out how to wind down the workers safely in the event of an error
	during initialization, during a channel closure, or other context
	cancelation")

	a := new(API)

	ces, err := a.getCuratedEpisodeSource()
	if err != nil {
		return nil, err
	}

	ccs, err := a.getCuratedClipSource()
	if err != nil {
		return nil, err
	}

	prs, err := a.getPendingResearchSource()
	if err != nil {
		return nil, err
	}

	crs, err := a.getCompletedResearchSource()
	if err != nil {
		return nil, err
	}

	errSrc := make(chan error)
	go func() {
		allChannelsOpen := true
		for allChannelsOpen {
			select {
			case ce, ok := <-ces:
				if !ok {
					allChannelsOpen = false
				} else {
					catch(a.processCuratedEpisode(ce), errSrc)
				}
			case cc, ok := <-ccs:
				if !ok {
					allChannelsOpen = false
				} else {
					catch(a.processCuratedClip(cc), errSrc)
				}
			case pr, ok := <-prs:
				if !ok {
					allChannelsOpen = false
				} else {
					catch(a.processPendingResearch(pr), errSrc)
				}
			case cr, ok := <-crs:
				if !ok {
					allChannelsOpen = false
				} else {
					catch(a.processCompletedResearch(cr), errSrc)

				}
			}
		}
	}()

	a.errorSource = errSrc
	return a, nil
}

func catch(err error, ch chan<- error) {
	if err != nil {
		ch <- err
	}
}

//TODO: Set Qos for channels to control how much work is buffered for
// each consumer instance.

func (a *API) getCuratedEpisodeSource() (<-chan contracts.EpisodeInfo, error) {
	// episodes are unique by name + date aired
	// check to see if the episode exists
	//	if it does not: add it
	//	else:
	//		check if any of its details of changed
	//		if so, update the details
	//	etc...
	panic("not implemented")
}

func (a *API) processCuratedEpisode(episode contracts.EpisodeInfo) error {
	panic("not implemented")
}

func (a *API) getCuratedClipSource() (<-chan contracts.ClipInfo, error) {
	// similar process to episode handling, except clips are unique by name only
	panic("not implemented")
}

func (a *API) processCuratedClip(clip contracts.ClipInfo) error {
	panic("not implemented")
}

func (a *API) getPendingResearchSource() (<-chan contracts.ResearchPending, error) {
	// check to see how many consumers there are for a queue
	// compare the consumer count to the message count
	// Then determine how much work to create, ie consumerCount - messageCount
	// Create that much work (including leases) and send it to the queue
	panic("not implemented")
}

func (a *API) processPendingResearch(pendingResearch contracts.ResearchPending) error {
	panic("not implemented")
}

func (a *API) getCompletedResearchSource() (<-chan contracts.ResearchComplete, error) {
	// upsert research and update leases if applicable
	panic("not implemented")
}

func (a *API) processCompletedResearch(completedResearch contracts.ResearchComplete) error {
	panic("not implemented")
}
