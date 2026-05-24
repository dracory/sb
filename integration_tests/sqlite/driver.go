package sb_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/sb"
)

// SQLiteDriver implements common.TestDriver for SQLite
type SQLiteDriver struct{}

func (d *SQLiteDriver) Name() string {
	return "SQLite"
}

func (d *SQLiteDriver) Dialect() string {
	return sb.DIALECT_SQLITE
}

func (d *SQLiteDriver) Connect() (*sql.DB, error) {
	return sql.Open("sqlite", ":memory:")
}

func (d *SQLiteDriver) SkipIfUnavailable(t *testing.T) {
	// SQLite is always available (pure Go driver)
}

func (d *SQLiteDriver) Cleanup(t *testing.T, db *sql.DB) {
	// In-memory database doesn't need cleanup
	_ = db.Close()
}
