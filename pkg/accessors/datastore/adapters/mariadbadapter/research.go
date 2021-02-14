package mariadbadapter

import (
	"database/sql"
	"fmt"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// GetHighestPriorityEpisode identifies and returns the highest priority
// episode to be researched. If no episodes are available, this returns nil,
// nil.
func (m *MariaDbConnection) GetHighestPriorityEpisode() (*contracts.EpisodeInfo, error) {
	const selectStmt = `
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
			cd.initial_date_curated DESC
		LIMIT 1;
	`

	row := m.db.QueryRow(selectStmt)
	episodeInfo := contracts.EpisodeInfo{}
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

	return &episodeInfo, nil
}

// GetHighestPriorityClipsForEpisode identifies and returns the highest
// priority clips to be researched for given episode. The number of clips
// returned is limited to `clipLimit`. If no clips are available for the
// supplied episode, this returns nil, nil.
func (m *MariaDbConnection) GetHighestPriorityClipsForEpisode(episode contracts.EpisodeInfo, clipLimit int) ([]contracts.ClipInfo, error) {
	const selectStmt = `
		SELECT
			cc.initial_date_curated,
			cc.last_date_curated,
			cc.curator_info,
			cc.title,
			cc.description,
			cc.media_uri,
			cc.media_type,
			cc.priority
		FROM 
			research_backlog rb
			LEFT JOIN research_leases rl ON rb.research_id = rl.research_id
			JOIN curated_clips cc ON rb.clip_id = cc.clip_id
		WHERE
			rl.research_id IS NULL
		ORDER BY
			cc.priority DESC,
			cc.initial_date_curated DESC
		LIMIT %v;
	`
	preparedStmt := fmt.Sprintf(selectStmt, clipLimit)

	rows, err := m.db.Query(preparedStmt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	clips := make([]contracts.ClipInfo, clipLimit)
	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, err
		}

		clip := contracts.ClipInfo{}
		err = rows.Scan(
			&clip.InitialDateCurated,
			&clip.LastDateCurated,
			&clip.CuratorInformation,
			&clip.Title,
			&clip.Description,
			&clip.MediaURI,
			&clip.MediaType,
			&clip.Priority,
		)

		if err != nil {
			return nil, err
		}

		clips = append(clips, clip)
	}

	if len(clips) == 0 {
		return nil, nil
	}

	return clips, nil
}

// RecordCompletedResearch inserts a reserach item. The system currently
// presumes that research is only assigned and conducted from the backlog
// (episodes/clip pairs that have not previously been researched). Submitting
// research for an episode/clip pair that has previously been researched is not
// supported, and will result in an error (though database integrity is
// maintained if this occurs).
func (m *MariaDbConnection) RecordCompletedResearch(completedResearchItem contracts.CompletedResearchItem) error {
	found, episodeID, err := m.getEpisodeInfoID(completedResearchItem.Episode)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("episodeID not found for episode: %v", completedResearchItem.Episode)
	}

	found, clipID, err := m.getClipInfoID(completedResearchItem.Clip)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("clipID not found for clip: %v", completedResearchItem.Clip)
	}

	found, researchID, err := m.getResearchIDFromBacklog(episodeID, clipID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("researchID not found for researchItem: %v", completedResearchItem)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	const deleteStmt = `DELETE FROM research_backlog WHERE research_id = ?`
	_, err = tx.Exec(deleteStmt, researchID)
	if err != nil {
		return tryTxRollback(tx, err)
	}

	const insertStmt = `
		INSERT INTO research_complete (
			research_id,
			episode_id,
			clip_id,
			episode_duration_ns,
			research_date
		) VALUES (?,?,?,?,?);
	`
	sqlResult, err := tx.Exec(insertStmt,
		researchID,
		episodeID,
		clipID,
		completedResearchItem.EpisodeDuration,
		completedResearchItem.ClipDuration,
		completedResearchItem.ResearchDate,
	)
	if err != nil {
		return tryTxRollback(tx, err)
	}
	if err := expectOneRowAffected(sqlResult, nil); err != nil {
		return tryTxRollback(tx, err)
	}

	return tx.Commit()
}

func (m *MariaDbConnection) getResearchIDFromBacklog(episodeID, clipID int) (bool, int, error) {
	const selectStmt = `
		SELECT research_id 
		FROM research_backlog 
		WHERE episode_id = ? AND clip_id = ?;
	`
	row := m.db.QueryRow(selectStmt, episodeID, clipID)
	var researchID int
	err := row.Scan(&researchID)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, researchID, nil
}
