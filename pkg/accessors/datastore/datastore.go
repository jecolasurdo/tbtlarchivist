package datastore

import (
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

type DataStorer interface {
	UpsertClipInfo(contracts.ClipInfo) error
	UpsertEpisodeInfo(contracts.EpisodeInfo) error
	GetMostRecentUnleasedEpisode() (*contracts.EpisodeInfo, error)
	SetEpisodeLease(contracts.EpisodeInfo, time.Time) error
	GetUnresearchedClipsForEpisode(contracts.EpisodeInfo) ([]contracts.ClipInfo, error)
}
