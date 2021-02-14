package mariadbadapter

import (
	"database/sql"
	"fmt"
	"time"
)

// MariaDb is an adapter that plugs into a mariadb instance.
type MariaDb struct {
	config                *Config
	maxConnectionLifetime time.Duration
	maxOpenConnections    int
	maxIdleConnections    int
}

// MariaDbConnection represents a successful connection to a mariadb instance.
type MariaDbConnection struct {
	db *sql.DB
}

// New returns a reference to a new MariaDb instance.
func New(config *Config) *MariaDb {
	return &MariaDb{
		config: config,
	}
}

// Connect attempts to open a connection to the underlaying mariadb instance.
func (m *MariaDb) Connect() (*MariaDbConnection, error) {
	db, err := sql.Open("mysql", m.config.formatDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(m.config.MaxConnectionLifetime)
	db.SetMaxOpenConns(m.config.MaxOpenConnections)
	db.SetMaxIdleConns(m.config.MaxIdleConnections)
	return &MariaDbConnection{
		db: db,
	}, nil
}

// expectOneRowAffected evaluates a sql.Result and an error. If err is not nil,
// the function immediately returns err. If err is not nil, then the function
// evaluates sql.Result. If sql.Result.Error is not nil, that error is
// returned.  If sql.Result.Error is nil, then sql.Result.RowsAffected is
// checked. If the value is not 1, then an error is returned. Else, nil is
// returned.
func expectOneRowAffected(result sql.Result, err error) error {
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("expected one row to be affected, but %v were affected", rowsAffected)
	}

	return nil
}

// tryTxRollback attempts to roll back the supplied transaction. If tx is nil,
// the function will return previousErr. If tx is not nil, tx.Rollback is
// called. If tx.Rollback returns an error, the rollback error and previousErr
// are joined and returned as a single error. If the rollback succeeds without
// error, previousErr is returned. previousErr is permitted to be nil.
func tryTxRollback(tx *sql.Tx, previousErr error) error {
	if tx == nil {
		return previousErr
	}

	err := tx.Rollback()
	if err != nil && previousErr != nil {
		err = fmt.Errorf("%v\n%v", previousErr, err)
	}

	return err
}
