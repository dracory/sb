package sb_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"

	"github.com/dracory/sb"
)

// PostgreSQLDriver implements common.TestDriver for PostgreSQL
type PostgreSQLDriver struct{}

func (d *PostgreSQLDriver) Name() string {
	return "PostgreSQL"
}

func (d *PostgreSQLDriver) Dialect() string {
	return sb.DIALECT_POSTGRES
}

func (d *PostgreSQLDriver) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", "55432", "test", "test", "test")
	return sql.Open("postgres", dsn)
}

func (d *PostgreSQLDriver) SkipIfUnavailable(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("PostgreSQL integration tests only run in GitHub Actions")
	}

	// Skip if PostgreSQL is not available
	db, err := d.Connect()
	if err != nil {
		t.Skip("PostgreSQL not available for integration testing:", err)
		return
	}

	// Test connection with retry logic
	var pingErr error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}
	if pingErr != nil {
		t.Skip("PostgreSQL not available after connection attempts:", pingErr)
	}

	db.Close()
}

func (d *PostgreSQLDriver) Cleanup(t *testing.T, db *sql.DB) {
	_ = db.Close()
}
