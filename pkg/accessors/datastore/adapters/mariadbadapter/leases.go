package mariadbadapter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// CreateResearchLease attempts to create a lease that is shared between the
// episode and all provided clips.  This method will panic if clips is empty or
// nil. If there are no clips to lease for an episode, that should be handled
// without attempting to call this method.
func (m *MariaDbConnection) CreateResearchLease(newLeaseID uuid.UUID, episode contracts.EpisodeInfo, clips []contracts.ClipInfo, expiration time.Time) error {
	if len(clips) == 0 {
		panic("a non-zero number of clips must be supplied to this method")
	}

	found, episodeID, err := m.getEpisodeInfoID(episode)
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("episode not found in database %v", episode)
	}

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	const insertStmt = `
		INSERT INTO research_leases (lease_id, research_id, expiration)
			SELECT "%v", bl.research_id, "%v" 
			FROM research_backlog bl
				JOIN curated_clips cc ON bl.clip_id = cc.clip_id
			WHERE bl.episode_id = ? AND cc.title = ?
	`
	// The sql engine cannot prepare placeholders in select clauses, so we
	// prepare those manually. Any placeholders that can be prepared by the sql
	// engine are prepared there.
	partiallyPreparedInsert := fmt.Sprintf(insertStmt, newLeaseID, expiration.UTC().Format(time.RFC3339Nano))

	for _, clip := range clips {
		sqlResult, err := tx.Exec(partiallyPreparedInsert, episodeID, clip.Title)
		if err != nil {
			return tryTxRollback(tx, err)
		}

		if err := expectOneRowAffected(sqlResult, nil); err != nil {
			return tryTxRollback(tx, err)
		}
	}

	return tx.Commit()
}

// RenewResearchLease updates the deadline for an existing lease. If the lease
// doesn't exist, no action is taken.
func (m *MariaDbConnection) RenewResearchLease(leaseID uuid.UUID, expiration time.Time) error {
	const selectStmt = `
		UPDATE research_leases
		SET expiration = ?
		WHERE lease_id = ?
	`
	// We ignore the returned SQLResult value since we're not concerned with
	// how many (if any) leases were renewed.
	_, err := m.db.Exec(selectStmt, expiration, leaseID)
	return err
}

// RevokeResearchLease removes the leases for all items assigned to the
// specified leaseID. If the leaseID doesn't exist, no action is taken.
func (m *MariaDbConnection) RevokeResearchLease(leaseID uuid.UUID) error {
	const deleteStmt = `
		DELETE FROM research_leases
		WHERE lease_id = ?
	`
	// We ignore the returned SQLResult value since we're not concerned with
	// how many (if any) leases were revoked.
	_, err := m.db.Exec(deleteStmt, leaseID)
	return err
}
