package mariadbadapter

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

// Config is a configuration for a mariadb instance.
type Config struct {
	Addr                  string
	DBName                string
	User                  string
	MaxConnectionLifetime time.Duration
	MaxOpenConnections    int
	MaxIdleConnections    int
}

func (c *Config) formatDSN() string {
	dbconfig := mysql.NewConfig()
	dbconfig.Addr = c.Addr
	dbconfig.DBName = c.DBName
	dbconfig.User = c.User
	dbconfig.ParseTime = true
	dbconfig.Loc = time.UTC

	// Internal error checking requires that Update statements return the number
	// of rows matched, not just the number of rows altered.
	dbconfig.ClientFoundRows = true

	// If interpolateParams is not explicitly set via the supplied constring,
	// then we set it to true here to ensure parameters are escaped
	// client-side.  This reduces the number of TCP calls to the db server.
	// This does impose some limitations that don't currently apply to this
	// system.  See https://github.com/go-sql-driver/mysql#interpolateparams
	dbconfig.InterpolateParams = true

	return dbconfig.FormatDSN()
}
