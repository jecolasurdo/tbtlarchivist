package contracts

import "encoding/json"

// A PendingResearchItem represents an episode and a list of of associated
// clips to research for that episode.
type PendingResearchItem struct {
	Episode EpisodeInfo
	Clips   []ClipInfo
}

// String returns a string representation of the PendingResearchItem instance.
func (p PendingResearchItem) String() string {
	jsonBytes, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
