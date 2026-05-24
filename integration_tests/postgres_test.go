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

// TestPostgreSQLIntegration tests PostgreSQL database integration with zero-panic error handling
// Only runs in GitHub Actions environment
func TestPostgreSQLIntegration(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("Integration tests only run in GitHub Actions")
		return
	}

	// Skip if PostgreSQL is not available
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", "127.0.0.1", "55432", "test", "test", "test")
	t.Logf("PostgreSQL DSN: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Skip("PostgreSQL not available for integration testing:", err)
		return
	}
	defer db.Close()

	// Test connection with retry logic (PostgreSQL service might need time to be ready)
	var pingErr error
	maxRetries := 5
	for i := 0; i < maxRetries; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			t.Logf("Successfully connected to PostgreSQL on attempt %d", i+1)
			break
		}
		t.Logf("Attempt %d: Failed to ping PostgreSQL: %v", i+1, pingErr)
		if i < maxRetries-1 {
			time.Sleep(2 * time.Second)
		}
	}
	if pingErr != nil {
		t.Fatalf("Failed to connect to PostgreSQL after %d attempts: %v", maxRetries, pingErr)
	}

	// Test table creation with unique name
	tableName := "test_users_postgres"
	err = createTestTable(db, tableName, sb.DIALECT_POSTGRES)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer dropTestTable(db, tableName) // Clean up after test

	// Test successful SQL generation and execution
	sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
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

	t.Logf("Successfully executed PostgreSQL integration test")
}
