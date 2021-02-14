package mariadbadapter

import "github.com/jecolasurdo/tbtlarchivist/pkg/contracts"

// GetHighestPriorityEpisode identifies and returns the highest priority
// episode to be researched. If no episodes are available, this returns nil,
// nil.
func (m *MariaDbConnection) GetHighestPriorityEpisode() (*contracts.EpisodeInfo, error) {
	panic("not implemented")
}

// GetHighestPriorityClipsForEpisode identifies and returns the highest
// priority clips to be researched for given episode. The number of clips
// returned is limited to `clipLimit`. If no clips are available for the
// supplied episode, this returns nil, nil.
func (m *MariaDbConnection) GetHighestPriorityClipsForEpisode(episode contracts.EpisodeInfo, clipLimit int) ([]contracts.ClipInfo, error) {
	panic("not implemented")
}

// UpsertCompletedResearch inserts or updates a reserach item.
func (m *MariaDbConnection) UpsertCompletedResearch(completedResearchItem contracts.CompletedResearchItem) error {
	panic("not implemented")
}
