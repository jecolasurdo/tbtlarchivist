package archivist

import (
	"sync"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// API is an instance of an archivist. An archivist is responsible for the
// following tasks:
//  - Recording episode and clip metadata as reported from the curators.
//  - Creating work for researchers.
//  - Checking in work returned from the researchers.
type API struct {
	ErrorSource <-chan error
}

type worker struct {
	source   <-chan interface{}
	delegate func(interface{}) error
}

// event loop
// 	check for curated episode
//  check for curated clip
//	check if there's work to be done for the researchers
//	check for work returned from researchers

func Initialize() (*API, error) {
	workerSource := []<-chan worker{}
	getCuratedEpisodeSource()

}

func (a *API) poll(workerSource []<-chan worker) <-chan error {
	backlog := make(chan worker)
	go func() {
		defer close(backlog)
		wg := new(sync.WaitGroup)
		for _, w := range workerSource {
			wg.Add(1)
			go func(w <-chan worker) {
				for workItem := range w {
					backlog <- workItem
				}
				wg.Done()
			}(w)
		}
		wg.Wait()
	}()

	errorSource := make(chan error)
	go func() {
		defer close(errorSource)
		for work := range backlog {
			errorSource <- work.delegate(work.source)
		}
	}()

	return errorSource
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

func processPendingResearch(pendingResearch contracts.ResearchPending) error {
	panic("not implemented")
}

func (a *API) getCompletedResearchSource() (<-chan contracts.ResearchComplete, error) {
	// upsert research and update leases if applicable
	panic("not implemented")
}

func processCompletedResearch(completedResearch contracts.ResearchComplete) error {
	panic("not implemented")
}
