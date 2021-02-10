package mariadbadapter

import "github.com/jecolasurdo/tbtlarchivist/pkg/contracts"

// GetUnresearchedClipsForEpisode returns a list of clips that have not been
// researched for a given episode. If the episode doesn't exist, the metehod
// will return nil, nil. If the episode exists, but there are no unreasearched
// clips for the episode, the method will return nil, nil.
func (m *MariaDbConnection) GetUnresearchedClipsForEpisode(episodeInto contracts.EpisodeInfo) ([]contracts.ClipInfo, error) {
	panic("not implemented")
}

// UpsertEpisodeClipInfo inserts or updates the information about a clip for a
// given episode. If the episode or clip does not exist, this method will
// return nil rather than returning an error.
func (m *MariaDbConnection) UpsertEpisodeClipInfo(researchItem contracts.CompletedResearchItem) error {
	panic("not implemented")
}
