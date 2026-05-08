package sb

import (
	"strings"

	"github.com/dracory/database"
)

// TableColumnDropSql generates the SQL for dropping a column from a table
func TableColumnDropSql(ctx database.QueryableContext, tableName string, columnName string) (string, error) {
	if ctx.Queryable() == nil {
		return "", ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sql, err := NewBuilder(databaseType).TableColumnDrop(tableName, columnName)
	if err != nil {
		return "", err
	}
	return sql, nil
}

// TableColumnDrop executes the SQL to drop a column from a table
func TableColumnDrop(ctx database.QueryableContext, tableName string, columnName string) error {
	sqlColumnDrop, err := TableColumnDropSql(ctx, tableName, columnName)

	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlColumnDrop)

	return err
}

// TableColumnDropIfExistsSql generates the SQL for dropping a column from a table if it exists
// Note: Not all databases support IF EXISTS for DROP COLUMN, this function will generate standard DROP COLUMN
func TableColumnDropIfExistsSql(ctx database.QueryableContext, tableName string, columnName string) (string, error) {
	if ctx.Queryable() == nil {
		return "", ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	// Check if column exists first
	exists, err := TableColumnExists(ctx, tableName, columnName)
	if err != nil {
		return "", err
	}

	if !exists {
		return "", nil // Column doesn't exist, no error
	}

	sql, err := NewBuilder(databaseType).TableColumnDrop(tableName, columnName)
	if err != nil {
		return "", err
	}
	return sql, nil
}

// TableColumnDropIfExists executes the SQL to drop a column from a table if it exists
func TableColumnDropIfExists(ctx database.QueryableContext, tableName string, columnName string) error {
	sqlColumnDrop, err := TableColumnDropIfExistsSql(ctx, tableName, columnName)

	if err != nil {
		return err
	}

	if sqlColumnDrop == "" {
		return nil // Column doesn't exist, nothing to drop
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlColumnDrop)

	// If the column was dropped between the existence check and execution,
	// return nil to maintain the "if exists" contract
	// This handles the TOCTOU race condition
	if err != nil {
		// Check if error indicates column doesn't exist (database-specific)
		// Common error patterns for non-existent columns:
		// MySQL: Error 1090 (42000): Can't DROP 'column_name'; check that column/key exists
		// PostgreSQL: ERROR: column "column_name" of relation "table_name" does not exist
		// SQLite: no such column: column_name
		// MSSQL: Invalid column name 'column_name'
		errStr := err.Error()
		if containsColumnNotFoundError(errStr) {
			return nil
		}
		return err
	}

	return nil
}

// containsColumnNotFoundError checks if error string indicates column doesn't exist
func containsColumnNotFoundError(errStr string) bool {
	columnErrorPatterns := []string{
		"Can't DROP",
		"does not exist",
		"no such column",
		"Invalid column name",
	}
	for _, pattern := range columnErrorPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	return false
}
