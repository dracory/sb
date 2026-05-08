package sb_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

func TestTableColumnRenameSql(t *testing.T) {
	// Test SQL generation for different dialects
	testCases := []struct {
		name          string
		dialect       string
		tableName     string
		oldColumnName string
		newColumnName string
		expected      string
	}{
		{
			name:          "MySQL",
			dialect:       sb.DIALECT_MYSQL,
			tableName:     "users",
			oldColumnName: "email",
			newColumnName: "new_email",
			expected:      "ALTER TABLE `users` RENAME COLUMN `email` TO `new_email`;",
		},
		{
			name:          "PostgreSQL",
			dialect:       sb.DIALECT_POSTGRES,
			tableName:     "users",
			oldColumnName: "email",
			newColumnName: "new_email",
			expected:      `ALTER TABLE "users" RENAME COLUMN "email" TO "new_email";`,
		},
		{
			name:          "SQLite",
			dialect:       sb.DIALECT_SQLITE,
			tableName:     "users",
			oldColumnName: "email",
			newColumnName: "new_email",
			expected:      `ALTER TABLE "users" RENAME COLUMN "email" TO "new_email";`,
		},
		{
			name:          "MSSQL",
			dialect:       sb.DIALECT_MSSQL,
			tableName:     "users",
			oldColumnName: "email",
			newColumnName: "new_email",
			expected:      "EXEC sp_rename 'users.email', 'new_email', 'COLUMN';",
		},
		{
			name:          "MSSQL with quotes in name",
			dialect:       sb.DIALECT_MSSQL,
			tableName:     "user's",
			oldColumnName: "email",
			newColumnName: "new_email",
			expected:      "EXEC sp_rename 'user''s.email', 'new_email', 'COLUMN';",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sql, err := sb.NewBuilder(tc.dialect).TableColumnRename(tc.tableName, tc.oldColumnName, tc.newColumnName)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if sql != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, sql)
			}
		})
	}
}

func TestTableColumnRenameSqlNilQueryable(t *testing.T) {
	ctx := database.Context(context.Background(), nil)
	_, err := sb.TableColumnRenameSql(ctx, "users", "old_name", "new_name")
	if err == nil {
		t.Fatal("Expected error for nil queryable")
	}
	if err.Error() != "ArgumentError: queryable cannot be nil" {
		t.Errorf("Expected 'ArgumentError: queryable cannot be nil', got: %v", err)
	}
}

func TestTableColumnRenameSqlEmptyNames(t *testing.T) {
	// Test empty table name
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnRename("", "old_name", "new_name")
	if err == nil {
		t.Fatal("Expected error for empty table name")
	}
	if err.Error() != "ValidationError: table name cannot be empty" {
		t.Errorf("Expected 'ValidationError: table name cannot be empty', got: %v", err)
	}

	// Test empty old column name
	_, err = sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnRename("users", "", "new_name")
	if err == nil {
		t.Fatal("Expected error for empty old column name")
	}
	if err.Error() != "ValidationError: column name cannot be empty" {
		t.Errorf("Expected 'ValidationError: column name cannot be empty', got: %v", err)
	}

	// Test empty new column name
	_, err = sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnRename("users", "old_name", "")
	if err == nil {
		t.Fatal("Expected error for empty new column name")
	}
	if err.Error() != "ValidationError: new column name cannot be empty" {
		t.Errorf("Expected 'ValidationError: new column name cannot be empty', got: %v", err)
	}

	// Test all empty
	_, err = sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnRename("", "", "")
	if err == nil {
		t.Fatal("Expected error for empty names")
	}
	// Should get table name error first
	if err.Error() != "ValidationError: table name cannot be empty" {
		t.Errorf("Expected 'ValidationError: table name cannot be empty', got: %v", err)
	}
}

func TestTableColumnRenameSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initSQLiteWithTable("test_column_rename", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Verify old column exists before rename
	exists, err := sb.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Old column should exist before rename")
	}

	// Rename the column
	err = sb.TableColumnRename(ctx, "test_column_rename", "email", "new_email")
	if err != nil {
		t.Fatal("Error renaming column: ", err.Error())
	}

	// Verify old column does not exist after rename
	exists, err = sb.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking old column existence after rename: ", err.Error())
	}
	if exists {
		t.Fatal("Old column should not exist after rename")
	}

	// Verify new column exists after rename
	exists, err = sb.TableColumnExists(ctx, "test_column_rename", "new_email")
	if err != nil {
		t.Fatal("Error checking new column existence after rename: ", err.Error())
	}
	if !exists {
		t.Fatal("New column should exist after rename")
	}
}

func TestTableColumnRenameMySQL(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnRenameMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initMySQLWithTable("test_column_rename", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Verify old column exists before rename
	exists, err := sb.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Old column should exist before rename")
	}

	// Rename the column
	err = sb.TableColumnRename(ctx, "test_column_rename", "email", "new_email")
	if err != nil {
		t.Fatal("Error renaming column: ", err.Error())
	}

	// Verify old column does not exist after rename
	exists, err = sb.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking old column existence after rename: ", err.Error())
	}
	if exists {
		t.Fatal("Old column should not exist after rename")
	}

	// Verify new column exists after rename
	exists, err = sb.TableColumnExists(ctx, "test_column_rename", "new_email")
	if err != nil {
		t.Fatal("Error checking new column existence after rename: ", err.Error())
	}
	if !exists {
		t.Fatal("New column should exist after rename")
	}
}
