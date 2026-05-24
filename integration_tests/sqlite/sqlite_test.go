package sb_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/sb"
	"github.com/dracory/sb/integration_tests/common"
)

// TestSQLiteIntegration tests SQLite database integration with zero-panic error handling
func TestSQLiteIntegration(t *testing.T) {
	// Use in-memory SQLite database (modernc.org/sqlite is pure Go, no CGO required)
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Test table creation with unique name
	tableName := "test_users_sqlite"
	err = common.CreateTestTable(db, tableName, sb.DIALECT_SQLITE)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer common.DropTestTable(db, tableName) // Clean up after test

	// Test successful SQL generation and execution
	sql, params, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	t.Logf("Successfully executed SQLite integration test")
}
