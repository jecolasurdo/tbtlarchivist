package contracts

import (
	"encoding/json"
	"time"
)

// ClipInfo contains information about an audio clip.
type ClipInfo struct {
	// InitialDateCurated represents the date that a curator service first
	// discovered this clip.
	InitialDateCurated time.Time

	// LastDateCurated respresents the most recent date that a curator
	// service found this clip.
	LastDateCurated time.Time

	// CuratorInformation provides information about the utility that extracted
	// this information.
	CuratorInformation string

	// Title is the name of the clip.
	Title string

	// Description is the clip description.
	Description string

	// MediaURI is a URI for where the clip media can be accessed.
	MediaURI string

	// MediaType is the media type for the episode (such as mp3, etc.).
	MediaType string

	// Priority is a value used by the curators to help the archivists prioritize
	// how clips are assigned to researchers.
	Priority int
}

// String returns a string representation of the ClipInfo instance.
func (c ClipInfo) String() string {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
