package mariadbadapter

import (
	"database/sql"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// GetHighestPriorityEpisode identifies and returns the highest priority
// episode to be researched. If no episodes are available, this returns nil,
// nil.
func (m *MariaDbConnection) GetHighestPriorityEpisode() (*contracts.EpisodeInfo, error) {
	selectStmt := `
		SELECT 
			ce.initial_date_curated,
			ce.last_date_curated,
			ce.curator_info,
			ce.date_aired,
			ce.title,
			ce.description,
			ce.media_uri,
			ce.media_type,
			ce.priority
		FROM 
			research_backlog rb
			LEFT JOIN research_leases rl ON rb.research_id = rl.research_id
			JOIN curated_episodes ce ON rb.episode_id = ce.episode_id
		WHERE
			rl.research_id IS NULL
		ORDER BY
			ce.priority DESC,
			cd.last_date_curated DESC
		LIMIT 1;
	`

	row := m.db.QueryRow(selectStmt)
	episodeInfo := new(contracts.EpisodeInfo)
	err := row.Scan(
		&episodeInfo.InitialDateCurated,
		&episodeInfo.LastDateCurated,
		&episodeInfo.CuratorInformation,
		&episodeInfo.DateAired,
		&episodeInfo.Title,
		&episodeInfo.Description,
		&episodeInfo.MediaURI,
		&episodeInfo.MediaType,
		&episodeInfo.Priority,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return episodeInfo, nil
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
