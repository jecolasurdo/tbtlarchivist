package mariadbadapter

import (
	"fmt"
	"strings"
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
func (m *MariaDbConnection) CreateResearchLease(newLeaseID string, episode contracts.EpisodeInfo, clips []contracts.ClipInfo, deadline time.Time) error {
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

	panic("TODO (see comment below)")
	// The way you're doing this is dumb.
	// The title concatenation is not terrible, but from there you should just
	// do some joins to idenfity the leasable set.
	// Let the database do the work.

	// verify each clip exists by getting their ids
	titles := make([]string, len(clips))
	for _, clip := range clips {
		titles = append(titles, fmt.Sprintf(`"%v"`, clip.Title))
	}

	const selectClipIDsStmt = `
		SELECT clip_id FROM curated_clips WHERE title IN [%v]
	`
	preparedInsertClipIDsStmt := fmt.Sprintf(selectClipIDsStmt, strings.Join(titles, ","))

	tx, err := m.db.Begin()
	if err != nil {
		return err
	}

	rows, err := tx.Query(preparedInsertClipIDsStmt)
	if err != nil {
		return tryTxRollback(tx, err)
	}

	clipIDs := make([]int, len(clips))
	for rows.Next() {
		err := rows.Err()
		if err != nil {
			return tryTxRollback(tx, err)
		}
		var clipID int
		err = rows.Scan(&clipID)
		if err != nil {
			return tryTxRollback(tx, err)
		}
		clipIDs = append(clipIDs, clipID)
	}

	if len(clipIDs) != len(clips) {
		err = fmt.Errorf("requested IDs for %v clips, but found IDs for %v clips", len(clips), len(clipIDs))
		return tryTxRollback(tx, err)
	}

	// verify that none of them have an active lease
	// insert a lease item for each clip/episode pair
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
