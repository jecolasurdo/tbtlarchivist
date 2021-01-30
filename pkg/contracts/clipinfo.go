package contracts

import (
	"encoding/json"
	"time"
)

// ClipInfo contains information about an audio clip.
type ClipInfo struct {
	// DateCurated respresents the date that the curator service found and
	// analyzed the clip.
	DateCurated time.Time

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
}

// String returns a string representation of the ClipInfo instance.
func (c ClipInfo) String() string {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
