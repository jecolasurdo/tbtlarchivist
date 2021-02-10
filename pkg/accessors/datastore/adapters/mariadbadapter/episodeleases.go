package mariadbadapter

import (
	"time"

	"github.com/jecolasurdo/tbtlarchivist/pkg/contracts"
)

// GetMostRecentUnleasedEpisode returns the most recent episode that does not
// have a current lease. If there is no current unleased episode (either
// because there are no episodes or because all episodes are currently leased),
// this method will return nil, nil.
func (m *MariaDbConnection) GetMostRecentUnleasedEpisode() (*contracts.EpisodeInfo, error) {
	panic("not implemented")
}

// SetEpisodeLease creates or updates a lease for an episode.
func (m *MariaDbConnection) SetEpisodeLease(episodeInfo contracts.EpisodeInfo, leaseExpiration time.Time) error {
	panic("not implemented")
}
