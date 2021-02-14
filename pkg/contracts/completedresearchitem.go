package contracts

import (
	"encoding/json"
	"time"
)

// A CompletedResearchItem represents the results of researching a clip for an
// episode.
type CompletedResearchItem struct {
	ResearchDate    time.Time
	Episode         EpisodeInfo
	EpisodeDuration Nanosecond
	Clip            ClipInfo
	ClipDuration    Nanosecond
	ClipOffsets     []Nanosecond
}

// String returns a string representation of the CompletedResearchItem instance.
func (c CompletedResearchItem) String() string {
	jsonBytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

// A Nanosecond is a duration equal to 1e9 seconds.
type Nanosecond int64