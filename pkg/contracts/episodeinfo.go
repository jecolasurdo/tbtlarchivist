package contracts

import (
	"encoding/json"
	"time"
)

// EpisodeInfo contains information about an episode.
type EpisodeInfo struct {
	// DateCurated represents the date that the curator service found and
	// analyzed the episode.
	DateCurated time.Time

	// CuratorInformation provides information about the utility that extracted
	// this information.
	CuratorInformation string

	// DateAired is the date that the episode was originally aired.
	DateAired time.Time

	// Duration is the length of the episode.
	Duration time.Duration

	// Title is the name of the episode.
	Title string

	// Description is the episode description.
	Description string

	// MediaURI is a URI for where the episode media can be accessed.
	MediaURI string

	// MediaType is the media type for the episode (such as mp3, etc).
	MediaType string
}

func (e EpisodeInfo) String() string {
	jsonBytes, err := json.MarshalIndent(e, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
