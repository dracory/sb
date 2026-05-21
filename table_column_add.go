package sb

import (
	"strings"

	"github.com/dracory/database"
)

// TableColumnAddSql generates the SQL for adding a column to a table
func TableColumnAddSql(ctx database.QueryableContext, tableName string, column Column) (string, error) {
	if ctx.Queryable() == nil {
		return "", ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sql, err := NewBuilder(databaseType).TableColumnAdd(tableName, column)
	if err != nil {
		return "", err
	}
	return sql, nil
}

// ExecuteTableColumnAdd executes the SQL to add a column to a table
func ExecuteTableColumnAdd(ctx database.QueryableContext, tableName string, column Column) error {
	if tableName == "" || column.Name == "" {
		return NewValidationError("table name and column name cannot be empty")
	}

	sqlColumnAdd, err := TableColumnAddSql(ctx, tableName, column)

	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlColumnAdd)

	return err
}

// TableColumnAddIfNotExistsSql generates the SQL for adding a column to a table if it doesn't exist
// Checks if the column exists first and only generates SQL if it doesn't.
func TableColumnAddIfNotExistsSql(ctx database.QueryableContext, tableName string, column Column) (string, error) {
	if ctx.Queryable() == nil {
		return "", ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	// Check if column exists first
	exists, err := TableColumnExists(ctx, tableName, column.Name)
	if err != nil {
		return "", err
	}

	if exists {
		return "", nil // Column already exists, no error
	}

	sql, err := NewBuilder(databaseType).TableColumnAdd(tableName, column)
	if err != nil {
		return "", err
	}
	return sql, nil
}

// ExecuteTableColumnAddIfNotExists executes the SQL to add a column to a table if it doesn't exist
// Checks if the column exists first and only adds if it doesn't.
// Handles TOCTOU race condition by swallowing duplicate column errors.
func ExecuteTableColumnAddIfNotExists(ctx database.QueryableContext, tableName string, column Column) error {
	if tableName == "" || column.Name == "" {
		return NewValidationError("table name and column name cannot be empty")
	}

	sqlColumnAdd, err := TableColumnAddIfNotExistsSql(ctx, tableName, column)

	if err != nil {
		return err
	}

	if sqlColumnAdd == "" {
		return nil // Column already exists, nothing to add
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlColumnAdd)

	// If the column was added between the existence check and execution,
	// return nil to maintain the "if exists" contract
	// This handles the TOCTOU race condition
	if err != nil {
		// Check if error indicates column already exists (database-specific)
		// Common error patterns for duplicate columns:
		// MySQL: Error 1060 (42S21): Duplicate column name 'column_name'
		// PostgreSQL: ERROR: column "column_name" of relation "table_name" already exists
		// SQLite: duplicate column name: column_name
		// MSSQL: Column names in each table must be unique. Column name 'column_name' is already used.
		errStr := err.Error()
		if containsColumnDuplicateError(errStr) {
			return nil
		}
		return err
	}

	return nil
}

// containsColumnDuplicateError checks if error string indicates column already exists
func containsColumnDuplicateError(errStr string) bool {
	// Check for specific column duplicate error patterns
	columnErrorPatterns := []string{
		"Duplicate column name",
		"duplicate column name",
		"Column names in each table must be unique",
	}
	for _, pattern := range columnErrorPatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}
	// For PostgreSQL/SQLite "already exists" errors, ensure they're column-specific
	// by checking for both "column" and "already exists" together
	if strings.Contains(errStr, "column") && strings.Contains(errStr, "already exists") {
		return true
	}
	return false
}
