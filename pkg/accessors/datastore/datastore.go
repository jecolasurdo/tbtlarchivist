package datastore

import (
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// A DataStorer is anything that is able to store and retrieve data from a
// persistant data store. Each of the methods of this interface must be
// implemented as atomic actions. Thus, if an error is returned by any of these
// methods, the consumer should trust that the underlaying data layer has made
// all possible attempts to leave the datastore in a safe state. If such a
// condition cannot be guaranteed, the implementor should expose error types
// that describe the severity of the issue.
type DataStorer interface {
	UpsertClipInfo(contracts.ClipInfo) error
	UpsertEpisodeInfo(contracts.EpisodeInfo) error
	GetMostRecentUnleasedEpisode() (*contracts.EpisodeInfo, error)
	SetEpisodeLease(contracts.EpisodeInfo, time.Time) error
	GetUnresearchedClipsForEpisode(contracts.EpisodeInfo) ([]contracts.ClipInfo, error)
	UpsertEpisodeClipInfo(contracts.CompletedResearchItem) error
}
