package datastore

import (
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// A DataStorer is anything that is able to store and retrieve data from a
// persistant data store.
type DataStorer interface {
	UpsertClipInfo(contracts.ClipInfo) error
	UpsertEpisodeInfo(contracts.EpisodeInfo) error
	GetMostRecentUnleasedEpisode() (*contracts.EpisodeInfo, error)
	SetEpisodeLease(contracts.EpisodeInfo, time.Time) error
	GetUnresearchedClipsForEpisode(contracts.EpisodeInfo) ([]contracts.ClipInfo, error)
	UpsertEpisodeClipInfo(contracts.CompletedResearchItem) error
}
