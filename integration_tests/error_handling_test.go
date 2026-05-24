package sb_test

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"github.com/dracory/sb"
)

// TestErrorHandlingIntegration tests error handling scenarios in real database contexts
func TestErrorHandlingIntegration(t *testing.T) {
	// Test with SQLite for simplicity (modernc.org/sqlite is pure Go, no CGO required)
	db, err := sql.Open("sqlite", ":memory:")
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
