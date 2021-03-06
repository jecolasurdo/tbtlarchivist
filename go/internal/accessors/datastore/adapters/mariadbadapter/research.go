package mariadbadapter

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			ce.date_aired DESC,
			ce.initial_date_curated DESC
		LIMIT 1;
	`

	row := m.db.QueryRow(selectStmt)
	episodeInfo := contracts.EpisodeInfo{}
	initialDateCurated := new(time.Time)
	lastDateCurated := new(time.Time)
	dateAired := new(time.Time)
	err := row.Scan(
		&initialDateCurated,
		&lastDateCurated,
		&episodeInfo.CuratorInformation,
		&dateAired,
		&episodeInfo.Title,
		&episodeInfo.Description,
		&episodeInfo.MediaUri,
		&episodeInfo.MediaType,
		&episodeInfo.Priority,
	)

	episodeInfo.InitialDateCurated = timestamppb.New(*initialDateCurated)
	episodeInfo.LastDateCurated = timestamppb.New(*lastDateCurated)
	episodeInfo.DateAired = timestamppb.New(*dateAired)

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
func (m *MariaDbConnection) GetHighestPriorityClipsForEpisode(episode *contracts.EpisodeInfo, clipLimit int) ([]*contracts.ClipInfo, error) {
	found, episodeID, err := m.getEpisodeInfoID(episode)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("episode not found: %v", episode)
	}

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
			AND rb.episode_id = ?
		ORDER BY
			cc.priority DESC,
			cc.initial_date_curated DESC
		LIMIT ?;
	`
	rows, err := m.db.Query(selectStmt, episodeID, clipLimit)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	clips := []*contracts.ClipInfo{}
	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return nil, err
		}

		clip := new(contracts.ClipInfo)
		initialDateCurated := new(time.Time)
		lastDateCurated := new(time.Time)
		err = rows.Scan(
			initialDateCurated,
			lastDateCurated,
			&clip.CuratorInformation,
			&clip.Title,
			&clip.Description,
			&clip.MediaUri,
			&clip.MediaType,
			&clip.Priority,
		)

		clip.InitialDateCurated = timestamppb.New(*initialDateCurated)
		clip.LastDateCurated = timestamppb.New(*lastDateCurated)

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

// RecordCompletedResearch inserts a research item. The system currently
// presumes that research is only assigned and conducted from the backlog
// (episodes/clip pairs that have not previously been researched). Submitting
// research for an episode/clip pair that has previously been researched is not
// supported, and will result in an error (though database integrity is
// maintained if this occurs). Episode and clip hash calculations are presumed
// to be deterministic. Thus, episode and clip hashes are only inserted; not
// updated. This is a means to an end, and may change in the future. Hashes are
// maintained for all researached clips and episodes, but "completed research"
// is only explicitly recorded for episode/clip pairs where the clip is found
// within the episode. If research is conducted for a clip/episode pair, and
// the clip is not found in the episode, the clip/episode pair is removed from
// the backlog, and not added to the completed research table. This is
// currently done to save space in the database since the vast majority of
// clip/episode pairs are non-matches.  Researched but negative pairings can be
// inferred as pairs that are in neither the backlog table nor the completed
// table.
func (m *MariaDbConnection) RecordCompletedResearch(completedResearchItem *contracts.CompletedResearchItem) error {
	found, episodeID, err := m.getEpisodeInfoID(completedResearchItem.EpisodeInfo)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("episodeID not found for episode: %v", completedResearchItem.EpisodeInfo)
	}

	found, clipID, err := m.getClipInfoID(completedResearchItem.ClipInfo)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("clipID not found for clip: %v", completedResearchItem.ClipInfo)
	}

	found, researchID, err := m.getResearchIDFromBacklog(episodeID, clipID)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("researchID not found for researchItem: %v", completedResearchItem)
	}

	foundClipHash := true
	const selectClipHashStmt = `SELECT 1 FROM clip_hashes WHERE clip_id = ?;`
	row := m.db.QueryRow(selectClipHashStmt, clipID)
	if row == nil || row.Err() == sql.ErrNoRows {
		foundClipHash = false
	} else if row.Err() != nil {
		return row.Err()
	}

	foundEpisodeHash := true
	const selectEpisodeHashStmt = `SELECT 1 FROM episode_hashes WHERE episode_id = ?;`
	row = m.db.QueryRow(selectEpisodeHashStmt, episodeID)
	if row == nil || row.Err() == sql.ErrNoRows {
		foundEpisodeHash = false
	} else if row.Err() != nil {
		return row.Err()
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	if !foundClipHash {
		const insertClipHashStmt = `INSERT INTO clip_hashes (clip_id, hash) VALUES (?,?);`
		sqlResult, err := tx.Exec(insertClipHashStmt, clipID, completedResearchItem.ClipHash)
		if err != nil {
			return tryTxRollback(tx, err)
		}
		if err := expectOneRowAffected(sqlResult, nil); err != nil {
			return tryTxRollback(tx, err)
		}
	}

	if !foundEpisodeHash {
		const insertEpisodeHashStmt = `INSERT INTO episode_hashes (episode_id, hash) VALUES (?,?);`
		sqlResult, err := tx.Exec(insertEpisodeHashStmt, episodeID, completedResearchItem.EpisodeHash)
		if err != nil {
			return tryTxRollback(tx, err)
		}
		if err := expectOneRowAffected(sqlResult, nil); err != nil {
			return tryTxRollback(tx, err)
		}
	}

	const deleteStmt = `DELETE FROM research_backlog WHERE research_id = ?`
	_, err = tx.Exec(deleteStmt, researchID)
	if err != nil {
		return tryTxRollback(tx, err)
	}

	if len(completedResearchItem.ClipOffsets) > 0 {
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
			completedResearchItem.ResearchDate.AsTime(),
		)
		if err != nil {
			return tryTxRollback(tx, err)
		}
		if err := expectOneRowAffected(sqlResult, nil); err != nil {
			return tryTxRollback(tx, err)
		}

		const insertOffsetStmt = `
			INSERT INTO episode_clip_offsets (research_id, offset_ns)
			VALUES (?, ?);
		`
		for _, offset := range completedResearchItem.ClipOffsets {
			sqlResult, err = tx.Exec(insertOffsetStmt, researchID, offset)
			if err != nil {
				return tryTxRollback(tx, err)
			}
			if err := expectOneRowAffected(sqlResult, nil); err != nil {
				return tryTxRollback(tx, err)
			}
		}
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
