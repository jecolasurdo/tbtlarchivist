package mariadbadapter

import (
	"fmt"
	"time"

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
func (m *MariaDbConnection) CreateResearchLease(newLeaseID string, episode contracts.EpisodeInfo, clips []contracts.ClipInfo, expiration time.Time) error {
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
			WHERE cc.title = ?
	`
	// Preparing the statement with newLeaseID and expiration here, but leaving
	// the title value to be prepared by the sql package. Placeholders can't be
	// used in the select clause, but we prefer to rely on the sql engine to
	// escape/convert datatypes as much as possible.
	partiallyPreparedInsert := fmt.Sprintf(insertStmt, newLeaseID, expiration)

	for _, clip := range clips {
		panic("TODO (see below)")
		// range over the clips
		// prepare each insert statement
		// if the insert fails, or alters 0 records rollback and return an error
	}

	panic("not implemented")
}

// RenewResearchLease updates the deadline for an existing lease. If the lease
// doesn't exist, no action is taken.
func (m *MariaDbConnection) RenewResearchLease(leaseID string, deadline time.Time) error {
	panic("not implemented")
}

// RevokeResearchLease removes the leases for all items assigned to the
// specified leaseID. If the leaseID doesn't exist, no action is taken.
func (m *MariaDbConnection) RevokeResearchLease(leaseID string) error {
	panic("not implemented")
}
