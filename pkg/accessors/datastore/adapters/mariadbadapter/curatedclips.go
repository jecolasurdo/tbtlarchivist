package mariadbadapter

import (
	"database/sql"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// UpsertClipInfo inserts or updates clip info.
func (m *MariaDbConnection) UpsertClipInfo(clipInfo contracts.ClipInfo) error {
	clipExists, clipID, err := m.getClipInfoID(clipInfo)
	if err != nil {
		return err
	}
	if clipExists {
		return m.updateClipInfo(clipID, clipInfo)
	}
	return m.insertClipInfo(clipInfo)
}

func (m *MariaDbConnection) getClipInfoID(clipInfo contracts.ClipInfo) (bool, int, error) {
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

func (m *MariaDbConnection) updateClipInfo(clipID int, clipInfo contracts.ClipInfo) error {
	const updateStmt = `
	UPDATE curated_clips
	SET date_curated = ?,
		curator_info = ?,
		title = ?,
		description = ?,
		media_uri = ?,
		media_type = ?
	WHERE clip_id = ?;
	`
	result, err := m.db.Exec(updateStmt,
		clipInfo.DateCurated,
		clipInfo.CuratorInformation,
		clipInfo.Title,
		clipInfo.Description,
		clipInfo.MediaURI,
		clipInfo.MediaType,
		clipID,
	)

	return expectOneRowAffected(result, err)
}

func (m *MariaDbConnection) insertClipInfo(clipInfo contracts.ClipInfo) error {
	const insertStmt = `
	INSERT INTO curated_clips (
		date_curated,
		curator_info,
		title,
		description,
		media_uri,
		media_type
	)
	VALUES (?,?,?,?,?,?);
	`
	result, err := m.db.Exec(insertStmt,
		clipInfo.DateCurated,
		clipInfo.CuratorInformation,
		clipInfo.Title,
		clipInfo.Description,
		clipInfo.MediaURI,
		clipInfo.MediaType,
	)
	return expectOneRowAffected(result, err)
}
