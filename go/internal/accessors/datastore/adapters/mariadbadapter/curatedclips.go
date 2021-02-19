package mariadbadapter

import (
	"database/sql"
	"fmt"

	"github.com/jecolasurdo/tbtlarchivist/go/internal/contracts"
)

// UpsertClipInfo inserts or updates clip info. If the clip already exists, it
// will be updated, but its InitialDateCurated value is ignored. If the clip
// does not already exist, both InitialDateCurated and LastDateCurated are
// evaluated, but the insert will fail and an error will be returned if
// LastDateCurated is earlier than InitialDateCurated.
func (m *MariaDbConnection) UpsertClipInfo(clipInfo *contracts.ClipInfo) error {
	clipExists, clipID, err := m.getClipInfoID(clipInfo)
	if err != nil {
		return err
	}
	if clipExists {
		return m.updateClipInfo(clipID, clipInfo)
	}
	return m.insertClipInfo(clipInfo)
}

func (m *MariaDbConnection) getClipInfoID(clipInfo *contracts.ClipInfo) (bool, int, error) {
	const selectStmt = `SELECT clip_id FROM curated_clips WHERE title = ?;`
	row := m.db.QueryRow(selectStmt, clipInfo.Title)
	var clipID int
	err := row.Scan(&clipID)
	if err == sql.ErrNoRows {
		return false, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return true, clipID, nil
}

func (m *MariaDbConnection) updateClipInfo(clipID int, clipInfo *contracts.ClipInfo) error {
	// Note that on updates, we update the `last_date_curated` field and ignore
	// the  `initial_date_curated` field.
	const updateStmt = `
	UPDATE curated_clips
	SET last_date_curated = ?,
		curator_info = ?,
		title = ?,
		description = ?,
		media_uri = ?,
		media_type = ?,
		priority = ?
	WHERE clip_id = ?;
	`
	result, err := m.db.Exec(updateStmt,
		clipInfo.LastDateCurated.AsTime(),
		clipInfo.CuratorInformation,
		clipInfo.Title,
		clipInfo.Description,
		clipInfo.MediaUri,
		clipInfo.MediaType,
		clipInfo.Priority,
		clipID,
	)

	return expectOneRowAffected(result, err)
}

func (m *MariaDbConnection) insertClipInfo(clipInfo *contracts.ClipInfo) error {
	if clipInfo.LastDateCurated.AsTime().Before(clipInfo.InitialDateCurated.AsTime()) {
		return fmt.Errorf("LastDateCurated must not be earlier than InitialDateCurated. %v", clipInfo)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return tryTxRollback(tx, err)
	}

	const insertCuratedClipStmt = `
		INSERT INTO curated_clips (
			initial_date_curated,
			last_date_curated,
			curator_info,
			title,
			description,
			media_uri,
			media_type,
			priority
		)
		VALUES (?,?,?,?,?,?,?,?);
	`
	result, err := tx.Exec(insertCuratedClipStmt,
		clipInfo.InitialDateCurated.AsTime(),
		clipInfo.LastDateCurated.AsTime(),
		clipInfo.CuratorInformation,
		clipInfo.Title,
		clipInfo.Description,
		clipInfo.MediaUri,
		clipInfo.MediaType,
		clipInfo.Priority,
	)

	if err := expectOneRowAffected(result, err); err != nil {
		return tryTxRollback(tx, err)
	}

	newClipID, err := result.LastInsertId()
	if err != nil {
		return tryTxRollback(tx, err)
	}

	insertClipBacklog := fmt.Sprintf(`
		INSERT INTO research_backlog (episode_id, clip_id)
		SELECT episode_id, %v
		FROM curated_episodes;
	`, newClipID)

	result, err = tx.Exec(insertClipBacklog)
	if err != nil {
		return tryTxRollback(tx, err)
	}

	return tx.Commit()
}
