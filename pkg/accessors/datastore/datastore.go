package datastore

import (
	"fmt"
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

// FakeDataStorer is a minimal fake data store implementation that does little
// more than print the name of the method called. This is mostly just a a
// placeholder used in development for use in sitations where a full mock
// package might be overkill at the moment.
type FakeDataStorer struct{}

// UpsertClipInfo .
func (f *FakeDataStorer) UpsertClipInfo(clipInfo contracts.ClipInfo) error {
	fmt.Println("UpsertClipInfo", clipInfo)
	return nil
}

// UpsertEpisodeInfo .
func (f *FakeDataStorer) UpsertEpisodeInfo(episodeInfo contracts.EpisodeInfo) error {
	fmt.Println("UpsertEpisodeInfo", episodeInfo)
	return nil
}

// GetMostRecentUnleasedEpisode .
func (f *FakeDataStorer) GetMostRecentUnleasedEpisode() (*contracts.EpisodeInfo, error) {
	fmt.Println("GetMostRecentUnleasedEpisode")
	return new(contracts.EpisodeInfo), nil
}

// SetEpisodeLease .
func (f *FakeDataStorer) SetEpisodeLease(episodeInfo contracts.EpisodeInfo, deadline time.Time) error {
	fmt.Println("SetEpisodeLease", episodeInfo, deadline)
	return nil
}

// GetUnresearchedClipsForEpisode .
func (f *FakeDataStorer) GetUnresearchedClipsForEpisode(episodeInfo contracts.EpisodeInfo) ([]contracts.ClipInfo, error) {
	fmt.Println("GetUnresearchedClipsForEpisode", episodeInfo)
	return nil, nil
}

// UpsertEpisodeClipInfo .
func (f *FakeDataStorer) UpsertEpisodeClipInfo(episodeClipInfo contracts.CompletedResearchItem) error {
	fmt.Println("UpsertEpisodeClipInfo", episodeClipInfo)
	return nil
}
