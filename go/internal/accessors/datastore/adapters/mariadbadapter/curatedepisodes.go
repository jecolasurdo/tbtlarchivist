package mariadbadapter

import (
	"database/sql"
	"fmt"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
)

// UpsertEpisodeInfo inserts or updates episode info. If the episode already exists, it
// will be updated, but its InitialDateCurated value is ignored. If the episode
// does not already exist, both InitialDateCurated and LastDateCurated are
// evaluated, but the insert will fail and an error will be returned if
// LastDateCurated is earlier than InitialDateCurated.
func (m *MariaDbConnection) UpsertEpisodeInfo(episodeInfo *contracts.EpisodeInfo) error {
	episodeExists, episodeID, err := m.getEpisodeInfoID(episodeInfo)
	if err != nil {
		return err
	}
	if episodeExists {
		return m.updateEpisodeInfo(episodeID, episodeInfo)
	}
	return m.insertEpisodeInfo(episodeInfo)
}

func (m *MariaDbConnection) getEpisodeInfoID(episodeInfo *contracts.EpisodeInfo) (bool, int, error) {
	const selectStmt = `
		SELECT episode_id 
		FROM curated_episodes 
		WHERE title = ?	AND date_aired = ?;`
	row := m.db.QueryRow(selectStmt, episodeInfo.Title, episodeInfo.DateAired.AsTime())
	var episodeID int
	err := row.Scan(&episodeID)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, episodeID, nil
}

func (m *MariaDbConnection) updateEpisodeInfo(episodeID int, episodeInfo *contracts.EpisodeInfo) error {
	// Note that on updates, we update the `last_date_curated` field and ignore
	// the  `initial_date_curated` field.
	const updateStmt = `
	UPDATE curated_episodes
	SET last_date_curated = ?,
		curator_info = ?,
		date_aired = ?,
		title = ?,
		description = ?,
		media_uri = ?,
		media_type = ?,
		priority = ?
	WHERE episode_id = ?;
	`
	result, err := m.db.Exec(updateStmt,
		episodeInfo.LastDateCurated.AsTime(),
		episodeInfo.CuratorInformation,
		episodeInfo.DateAired.AsTime(),
		episodeInfo.Title,
		episodeInfo.Description,
		episodeInfo.MediaUri,
		episodeInfo.MediaType,
		episodeInfo.Priority,
		episodeID,
	)

	return expectOneRowAffected(result, err)
}

func (m *MariaDbConnection) insertEpisodeInfo(episodeInfo *contracts.EpisodeInfo) error {
	if episodeInfo.LastDateCurated.AsTime().Before(episodeInfo.InitialDateCurated.AsTime()) {
		return fmt.Errorf("LastDateCurated must not be earlier than InitialDateCurated. %v", episodeInfo)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return tryTxRollback(tx, err)
	}

	const insertCuratedEpisodeStmt = `
		INSERT INTO curated_episodes (
			initial_date_curated,
			last_date_curated,
			curator_info,
			date_aired,
			title,
			description,
			media_uri,
			media_type,
			priority
		)
		VALUES (?,?,?,?,?,?,?,?,?);
	`
	result, err := tx.Exec(insertCuratedEpisodeStmt,
		episodeInfo.InitialDateCurated.AsTime(),
		episodeInfo.LastDateCurated.AsTime(),
		episodeInfo.CuratorInformation,
		episodeInfo.DateAired.AsTime(),
		episodeInfo.Title,
		episodeInfo.Description,
		episodeInfo.MediaUri,
		episodeInfo.MediaType,
		episodeInfo.Priority,
	)

	if err := expectOneRowAffected(result, err); err != nil {
		return tryTxRollback(tx, err)
	}

	newEpisodeID, err := result.LastInsertId()
	if err != nil {
		return tryTxRollback(tx, err)
	}

	insertEpisodeBacklog := fmt.Sprintf(`
		INSERT INTO research_backlog (episode_id, clip_id)
		SELECT %v, clip_id
		FROM curated_clips;
	`, newEpisodeID)

	result, err = tx.Exec(insertEpisodeBacklog)
	if err != nil {
		return tryTxRollback(tx, err)
	}

	return tx.Commit()
}
