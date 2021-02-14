package mariadbadapter

import (
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// CreateResearchLease attempts to create a lease that is shared between the
// episode and all provided clips. This operation is atomic. If the episode or
// any of the provided clips either a) don't exist or b) already have an active
// lease, this operation will fail and return an error.
func (m *MariaDbConnection) CreateResearchLease(episode contracts.EpisodeInfo, clips []contracts.ClipInfo, deadline time.Time) (string, error) {
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
