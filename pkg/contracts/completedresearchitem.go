package contracts

import (
	"encoding/json"
	"time"
)

// A CompletedResearchItem represents the results of researching a clip for an
// episode.
type CompletedResearchItem struct {
	Episode     EpisodeInfo
	Clip        ClipInfo
	ClipOffsets []time.Duration
}

// String returns a string representation of the CompletedResearchItem instance.
func (c CompletedResearchItem) String() string {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
