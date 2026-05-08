package sb_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

func TestTableColumnDropSql(t *testing.T) {
	// Test SQL generation for different dialects
	testCases := []struct {
		name       string
		dialect    string
		tableName  string
		columnName string
		expected   string
	}{
		{
			name:       "MySQL",
			dialect:    sb.DIALECT_MYSQL,
			tableName:  "users",
			columnName: "old_column",
			expected:   "ALTER TABLE `users` DROP COLUMN `old_column`;",
		},
		{
			name:       "PostgreSQL",
			dialect:    sb.DIALECT_POSTGRES,
			tableName:  "users",
			columnName: "old_column",
			expected:   `ALTER TABLE "users" DROP COLUMN "old_column";`,
		},
		{
			name:       "SQLite",
			dialect:    sb.DIALECT_SQLITE,
			tableName:  "users",
			columnName: "old_column",
			expected:   `ALTER TABLE "users" DROP COLUMN "old_column";`,
		},
		{
			name:       "MSSQL",
			dialect:    sb.DIALECT_MSSQL,
			tableName:  "users",
			columnName: "old_column",
			expected:   "ALTER TABLE [users] DROP COLUMN [old_column];",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sql, err := sb.NewBuilder(tc.dialect).TableColumnDrop(tc.tableName, tc.columnName)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if sql != tc.expected {
				t.Errorf("Expected: %s, Got: %s", tc.expected, sql)
			}
		})
	}
}

func TestTableColumnDropSqlNilQueryable(t *testing.T) {
	ctx := database.Context(context.Background(), nil)
	_, err := sb.TableColumnDropSql(ctx, "users", "column")
	if err == nil {
		t.Fatal("Expected error for nil queryable")
	}
	if err.Error() != "ArgumentError: queryable cannot be nil" {
		t.Errorf("Expected 'ArgumentError: queryable cannot be nil', got: %v", err)
	}
}

func TestTableColumnDropSqlEmptyNames(t *testing.T) {
	// Test empty table name
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnDrop("", "column")
	if err == nil {
		t.Fatal("Expected error for empty table name")
	}
	if err.Error() != "ValidationError: table name cannot be empty" {
		t.Errorf("Expected 'ValidationError: table name cannot be empty', got: %v", err)
	}

	// Test empty column name
	_, err = sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnDrop("users", "")
	if err == nil {
		t.Fatal("Expected error for empty column name")
	}
	if err.Error() != "ValidationError: column name cannot be empty" {
		t.Errorf("Expected 'ValidationError: column name cannot be empty', got: %v", err)
	}

	// Test both empty
	_, err = sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnDrop("", "")
	if err == nil {
		t.Fatal("Expected error for empty names")
	}
	// Should get table name error first
	if err.Error() != "ValidationError: table name cannot be empty" {
		t.Errorf("Expected 'ValidationError: table name cannot be empty', got: %v", err)
	}
}

func TestTableColumnDropIfExistsSqlEmptyNames(t *testing.T) {
	ctx := database.Context(context.Background(), nil)

	// Test empty table name - should fail on nil queryable first
	_, err := sb.TableColumnDropIfExistsSql(ctx, "", "column")
	if err == nil {
		t.Fatal("Expected error for empty table name")
	}
	// Since queryable is nil, we get ArgumentError first
	if err.Error() != "ArgumentError: queryable cannot be nil" {
		t.Errorf("Expected 'ArgumentError: queryable cannot be nil', got: %v", err)
	}

	// Test empty column name - should fail on nil queryable first
	_, err = sb.TableColumnDropIfExistsSql(ctx, "users", "")
	if err == nil {
		t.Fatal("Expected error for empty column name")
	}
	if err.Error() != "ArgumentError: queryable cannot be nil" {
		t.Errorf("Expected 'ArgumentError: queryable cannot be nil', got: %v", err)
	}
}

func TestTableColumnDropMySQL(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnDropMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "temp_column", Type: sb.COLUMN_TYPE_STRING, Length: 50},
	}

	db, err := initMySQLWithTable("test_column_drop", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Verify column exists before drop
	exists, err := sb.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist before drop")
	}

	// Drop the column
	err = sb.TableColumnDrop(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error dropping column: ", err.Error())
	}

	// Verify column does not exist after drop
	exists, err = sb.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence after drop: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist after drop")
	}
}

func TestTableColumnDropSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "temp_column", Type: sb.COLUMN_TYPE_STRING, Length: 50},
	}

	db, err := initSQLiteWithTable("test_column_drop", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Verify column exists before drop
	exists, err := sb.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist before drop")
	}

	// Drop the column
	err = sb.TableColumnDrop(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error dropping column: ", err.Error())
	}

	// Verify column does not exist after drop
	exists, err = sb.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence after drop: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist after drop")
	}
}

func TestTableColumnDropIfExistsMySQL(t *testing.T) {
	// Skip if not in GitHub Actions
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnDropIfExistsMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initMySQLWithTable("test_column_drop_if_exists", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Try to drop a column that doesn't exist - should not error
	err = sb.TableColumnDropIfExists(ctx, "test_column_drop_if_exists", "nonexistent_column")
	if err != nil {
		t.Fatal("Dropping non-existent column should not error, got: ", err.Error())
	}
}

func TestTableColumnDropIfExistsSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initSQLiteWithTable("test_column_drop_if_exists", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Try to drop a column that doesn't exist - should not error
	err = sb.TableColumnDropIfExists(ctx, "test_column_drop_if_exists", "nonexistent_column")
	if err != nil {
		t.Fatal("Dropping non-existent column should not error, got: ", err.Error())
	}
}

func TestTableColumnDropIfExistsWithActualColumn(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "temp_column", Type: sb.COLUMN_TYPE_STRING, Length: 50},
	}

	db, err := initSQLiteWithTable("test_column_drop_if_exists_real", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() {
		_ = db.Close()
	}()

	ctx := database.Context(context.Background(), db)

	// Verify column exists before drop
	exists, err := sb.TableColumnExists(ctx, "test_column_drop_if_exists_real", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist before drop")
	}

	// Drop the column using DropIfExists
	err = sb.TableColumnDropIfExists(ctx, "test_column_drop_if_exists_real", "temp_column")
	if err != nil {
		t.Fatal("Error dropping column: ", err.Error())
	}

	// Verify column does not exist after drop
	exists, err = sb.TableColumnExists(ctx, "test_column_drop_if_exists_real", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence after drop: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist after drop")
	}
}
