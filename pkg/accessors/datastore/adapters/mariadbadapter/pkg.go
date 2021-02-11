package mariadbadapter

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// MariaDb is an adapter that plugs into a mariadb instance.
type MariaDb struct {
	constring             string
	maxConnectionLifetime time.Duration
	maxOpenConnections    int
	maxIdleConnections    int
}

// MariaDbConnection represents a successful connection to a mariadb instance.
type MariaDbConnection struct {
	db *sql.DB
}

// New returns a reference to a new MariaDb instance.
func New(constring string, maxConnectionLifetime time.Duration, maxOpenConnections, maxIdleConnections int) *MariaDb {
	// If interpolateParams is not explicitly set via the supplied constring,
	// then we set it to true here to ensure parameters are escaped
	// client-side.  This reduces the number of TCP calls to the db server.
	// This does impose some limitations that don't currently apply to this
	// system.  See https://github.com/go-sql-driver/mysql#interpolateparams
	if !strings.Contains(constring, `interpolateParams`) {
		constring = constring + `?interpolateParams=true`
	}
	return &MariaDb{
		constring:             constring,
		maxConnectionLifetime: maxConnectionLifetime,
		maxOpenConnections:    maxOpenConnections,
		maxIdleConnections:    maxIdleConnections,
	}
}

// Connect attempts to open a connection to the underlaying mariadb instance.
func (m *MariaDb) Connect() (*MariaDbConnection, error) {
	db, err := sql.Open("mysql", m.constring)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(m.maxConnectionLifetime)
	db.SetMaxOpenConns(m.maxOpenConnections)
	db.SetMaxIdleConns(m.maxIdleConnections)
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
