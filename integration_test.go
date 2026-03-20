package sb_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/dracory/sb"
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
	err = createTestTable(db, tableName, sb.DIALECT_MYSQL)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer dropTestTable(db, tableName) // Clean up after test

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

// TestSQLiteIntegration tests SQLite database integration with zero-panic error handling
// Only runs in GitHub Actions environment
func TestSQLiteIntegration(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("Integration tests only run in GitHub Actions")
		return
	}

	// Use in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Test table creation with unique name
	tableName := "test_users_sqlite"
	err = createTestTable(db, tableName, sb.DIALECT_SQLITE)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}
	defer dropTestTable(db, tableName) // Clean up after test

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

// TestErrorHandlingIntegration tests error handling scenarios in real database contexts
// Only runs in GitHub Actions environment
func TestErrorHandlingIntegration(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("Integration tests only run in GitHub Actions")
		return
	}

	// Test with SQLite for simplicity
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	// Test 1: Missing table error
	sql, params, err := sb.NewBuilder(sb.DIALECT_SQLITE).Select([]string{"*"})
	_ = params // params not used in error test

	if err == nil {
		t.Fatal("Expected error for missing table but got none")
	}

	expectedError := "ValidationError: no table specified"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

	// Test 2: Empty JOIN condition error
	sql, params, err = sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("test_users").
		InnerJoin("orders", "").
		Select([]string{"*"})
	_ = params // params not used in error test

	if err == nil {
		t.Fatal("Expected error for empty JOIN condition but got none")
	}

	expectedError = "ValidationError: ON condition cannot be empty"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

	// Test 3: Unsupported dialect error (should be collected during creation)
	sql, params, err = sb.NewBuilder("unknown").
		Table("test_users").
		Select([]string{"*"})
	_ = params // params not used in error test

	if err == nil {
		t.Fatal("Expected error for unsupported dialect but got none")
	}

	expectedError = "ValidationError: unsupported dialect: unknown"
	if err.Error() != expectedError {
		t.Fatalf("Expected error '%s' but got '%s'", expectedError, err.Error())
	}

	// Test 4: Successful case after fixing errors
	sql, params, err = sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("test_users").
		Select([]string{"name", "email"})

	if err != nil {
		t.Fatalf("Unexpected error in valid case: %v", err)
	}

	// Create table and execute valid SQL
	err = createTestTable(db, "test_users", sb.DIALECT_SQLITE)
	if err != nil {
		t.Fatalf("Failed to create test table: %v", err)
	}

	_, err = db.Exec(sql)
	if err != nil {
		t.Fatalf("Failed to execute valid SQL: %v\nSQL: %s", err, sql)
	}

	t.Logf("Successfully executed error handling integration test")
}

// createTestTable creates a simple test table for integration testing using SB library
func createTestTable(db *sql.DB, tableName string, dialect string) error {
	// Create table using SB library with proper dialect-specific syntax
	createSQL, err := sb.NewBuilder(dialect).
		Table(tableName).
		Column(sb.Column{
			Name:          "id",
			Type:          sb.COLUMN_TYPE_INTEGER,
			AutoIncrement: true,
			PrimaryKey:    true,
		}).
		Column(sb.Column{
			Name:     "name",
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   100,
			Nullable: false, // NOT NULL
		}).
		Column(sb.Column{
			Name:     "email",
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   100,
			Nullable: false, // NOT NULL
		}).
		Column(sb.Column{
			Name:    "status",
			Type:    sb.COLUMN_TYPE_STRING,
			Length:  20,
			Default: "active",
		}).
		Column(sb.Column{
			Name:    "created_at",
			Type:    sb.COLUMN_TYPE_DATETIME,
			Default: "CURRENT_TIMESTAMP",
		}).
		CreateIfNotExists()

	if err != nil {
		return fmt.Errorf("failed to generate CREATE TABLE SQL: %w", err)
	}

	// Log the SQL for debugging
	fmt.Printf("DEBUG: Generated SQL for %s:\n%s\n", dialect, createSQL)

	_, err = db.Exec(createSQL)
	if err != nil {
		return fmt.Errorf("failed to execute CREATE TABLE: %w", err)
	}
	return nil
}

// dropTestTable drops a test table after testing
func dropTestTable(db *sql.DB, tableName string) error {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := db.Exec(dropSQL)
	if err != nil {
		return fmt.Errorf("failed to drop test table %s: %w", tableName, err)
	}
	return nil
}
