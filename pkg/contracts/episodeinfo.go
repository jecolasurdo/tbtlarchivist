package contracts

import (
	"encoding/json"
	"time"
)

// EpisodeInfo contains information about an episode.
type EpisodeInfo struct {
	// InitialDateCurated represents the date that a curator service first
	// discovered this episode.
	InitialDateCurated time.Time

	// LastDateCurated respresents the most recent date that a curator
	// service found this episode.
	LastDateCurated time.Time

	// CuratorInformation provides information about the utility that extracted
	// this information.
	CuratorInformation string

	// DateAired is the date that the episode was originally aired.
	DateAired time.Time

	// Title is the name of the episode.
	Title string

	// Description is the episode description.
	Description string

	// MediaURI is a URI for where the episode media can be accessed.
	MediaURI string

	// MediaType is the media type for the episode (such as mp3, etc).
	MediaType string

	// Priority is a value used by the curators to help the archivists
	// prioritize how episodes are assigned to researchers.
	Priority int
}

// String returns a string representation of the EpisodeInfo instance.
func (e EpisodeInfo) String() string {
	jsonBytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
