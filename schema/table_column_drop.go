package schema

import (
	"strings"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

// TableColumnDrop drops a column from a table.
func TableColumnDrop(ctx database.QueryableContext, tableName string, columnName string) error {
	if ctx.Queryable() == nil {
		return sb.ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).Table(tableName).ColumnDrop(columnName)
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	return err
}

// TableColumnDropIfExists drops a column from a table only if it exists.
// Handles TOCTOU race condition by swallowing column-not-found errors.
func TableColumnDropIfExists(ctx database.QueryableContext, tableName string, columnName string) error {
	if ctx.Queryable() == nil {
		return sb.ErrNilQueryable
	}

	exists, err := TableColumnExists(ctx, tableName, columnName)
	if err != nil {
		return err
	}

	if !exists {
		return nil
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).Table(tableName).ColumnDrop(columnName)
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	if err != nil {
		if containsColumnNotFoundError(err.Error()) {
			return nil
		}
		return err
	}

	return nil
}

// containsColumnNotFoundError checks if error string indicates column doesn't exist.
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
