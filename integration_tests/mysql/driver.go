package sb_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dracory/sb"
)

// MySQLDriver implements common.TestDriver for MySQL
type MySQLDriver struct{}

func (d *MySQLDriver) Name() string {
	return "MySQL"
}

func (d *MySQLDriver) Dialect() string {
	return sb.DIALECT_MYSQL
}

func (d *MySQLDriver) Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "root", "127.0.0.1", "3306", "test")
	return sql.Open("mysql", dsn)
}

func (d *MySQLDriver) SkipIfUnavailable(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
	}

	// Skip if MySQL is not available
	db, err := d.Connect()
	if err != nil {
		t.Skip("MySQL not available for integration testing:", err)
		return
	}
	db.Close()
}

func (d *MySQLDriver) Cleanup(t *testing.T, db *sql.DB) {
	_ = db.Close()
}
