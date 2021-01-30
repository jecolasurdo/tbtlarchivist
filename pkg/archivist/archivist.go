package archivist

import "github.com/jecolasurdo/tbtlarchivist/pkg/contracts"

// API is an instance of an archivist. An archivist is responsible for the
// following tasks:
//  - Recording episode and clip metadata as reported from the curators.
//  - Creating work for researchers.
//  - Checking in work returned from the researchers.
type API struct{}

// event loop
// 	check for curated episode
//  check for curated clip
//	check if there's work to be done for the researchers
//	check for work returned from researchers

func Start() error {
	api := new(API)
	errorTrap := make(chan error)

	curatedEpisodeSource, errorSource := api.getCuratedEpisodeSource()
	curatedClipSource, errorSource := api.getCuratedClipSource()
	pendingResearchSource, errorSource := api.getPendingResearchSource()
	completedResearchSource, errorSource := api.getCompletedResearchSource()

	for {
		select {

		//TODO: Set Qos for channels to control how much work is buffered for
		// each consumer instance.

		case curatedEpisode := <-curatedEpisodeSource:
			err := api.processCuratedEpisode(curatedEpisode)
			if err != nil {
				return err
			}
		case curatedClip := <-curatedClipSource:
			err := api.processCuratedClip(curatedClip)
			if err != nil {
				return err
			}
		case pendingResearch := <-pendingResearchSource:
			err := processPendingResearch(pendingResearch)
			if err != nil {
				return err
			}
		case completedResearch := <-completedResearchSource:
			err := processCompletedResearch(completedResearch)
			if err != nil {
				return err
			}
		}
	}

}

func (a *API) getCuratedEpisodeSource() (<-chan contracts.EpisodeInfo, <-chan error) {
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

func (a *API) getCuratedClipSource() (<-chan contracts.ClipInfo, <-chan error) {
	// similar process to episode handling, except clips are unique by name only
	panic("not implemented")
}

func (a *API) processCuratedClip(clip contracts.ClipInfo) error {
	panic("not implemented")
}

func (a *API) getPendingResearchSource() (<-chan contracts.ResearchPending, <-chan error) {
	// check to see how many consumers there are for a queue
	// compare the consumer count to the message count
	// Then determine how much work to create, ie consumerCount - messageCount
	// Create that much work (including leases) and send it to the queue
	panic("not implemented")
}

func processPendingResearch(pendingResearch contracts.ResearchPending) error {
	panic("not implemented")
}

func (a *API) getCompletedResearchSource() (<-chan contracts.ResearchComplete, <-chan error) {
	// upsert research and update leases if applicable
	panic("not implemented")
}

func processCompletedResearch(completedResearch contracts.ResearchComplete) error {
	panic("not implemented")
}
