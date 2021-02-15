package mariadbadapter

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// CreateResearchLease attempts to create a lease that is shared between the
// episode and all provided clips. This operation is atomic. If the episode or
// any of the provided clips either a) doesn't exist or b) already has an
// active lease, this operation will fail and return an error. This method also
// presumes the supplied `newLeaseID` value is sufficiently unique (i.e.  a
// UUID), and will return an error if the supplied ID is already in use. This
// method will panic if clips is empty or nil. If there are no clips to lease
// for an episode, that should be handled without attempting to call this
// method.
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
func (m *MariaDbConnection) RenewResearchLease(leaseID uuid.UUID, deadline time.Time) error {
	panic("not implemented")
}

// RevokeResearchLease removes the leases for all items assigned to the
// specified leaseID. If the leaseID doesn't exist, no action is taken.
func (m *MariaDbConnection) RevokeResearchLease(leaseID uuid.UUID) error {
	panic("not implemented")
}
