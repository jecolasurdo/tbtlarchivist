package archivist

import (
	"context"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

//TODO: Set Qos for channels to control how much work is buffered for
// each consumer instance.

func (a *API) getCuratedEpisodeSource(ctx context.Context) (<-chan contracts.EpisodeInfo, error) {

	// episodes are unique by name + date aired
	// check to see if the episode exists
	//	if it does not: add it
	//	else:
	//		check if any of its details of changed
	//		if so, update the details
	//	etc...
	panic("not implemented")
}

func (a *API) processCuratedEpisode(ctx context.Context, episode contracts.EpisodeInfo) error {
	panic("not implemented")
}
