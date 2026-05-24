package schema

import (
	"strings"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

// TableColumnAdd adds a column to a table.
func TableColumnAdd(ctx database.QueryableContext, tableName string, column sb.Column) error {
	if tableName == "" || column.Name == "" {
		return sb.NewValidationError("table name and column name cannot be empty")
	}

	if ctx.Queryable() == nil {
		return sb.ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).TableColumnAdd(tableName, column)
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	return err
}

// TableColumnAddIfNotExists adds a column to a table only if it does not already exist.
// Handles TOCTOU race condition by swallowing duplicate column errors.
func TableColumnAddIfNotExists(ctx database.QueryableContext, tableName string, column sb.Column) error {
	if tableName == "" || column.Name == "" {
		return sb.NewValidationError("table name and column name cannot be empty")
	}

	if ctx.Queryable() == nil {
		return sb.ErrNilQueryable
	}

	exists, err := TableColumnExists(ctx, tableName, column.Name)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).TableColumnAdd(tableName, column)
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	if err != nil {
		if containsColumnDuplicateError(err.Error()) {
			return nil
		}
		return err
	}

	return nil
}

// containsColumnDuplicateError checks if error string indicates column already exists.
func containsColumnDuplicateError(errStr string) bool {
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
	if strings.Contains(errStr, "column") && strings.Contains(errStr, "already exists") {
		return true
	}
	return false
}
