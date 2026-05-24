package sb_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dracory/sb"
	"github.com/dracory/sb/integration_tests/common"
)

// TestMySQLIntegration tests MySQL database integration with zero-panic error handling
// Only runs in GitHub Actions environment
func TestMySQLIntegration(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("Integration tests only run in GitHub Actions")
		return
	}

	// Skip if MySQL is not available
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", "root", "root", "127.0.0.1", "3306", "test")
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Skip("MySQL not available for integration testing:", err)
		return
	}
	defer db.Close()

	// Test table creation with unique name
	tableName := "test_users_mysql"
	err = common.CreateTestTable(db, tableName, sb.DIALECT_MYSQL)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer common.DropTestTable(db, tableName) // Clean up after test

	// Test successful SQL generation and execution
	sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table(tableName).
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		Select([]string{"name", "email"})

	if err != nil {
		t.Fatalf("Failed to generate SQL: %v", err)
	}

	// Execute the generated SQL with parameters
	_, err = db.Exec(sql, params...)
	if err != nil {
		t.Fatalf("Failed to execute SQL: %v\nSQL: %s", err, sql)
	}

	t.Logf("Successfully executed MySQL integration test")
}
