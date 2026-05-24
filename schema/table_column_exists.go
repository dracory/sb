package schema

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

// TableColumnExists checks if a column exists in a table.
func TableColumnExists(ctx database.QueryableContext, tableName string, columnName string) (bool, error) {
	db := ctx.Queryable()
	if db == nil {
		return false, sb.ErrNilQueryable
	}

	if tableName == "" || columnName == "" {
		return false, sb.NewValidationError("table name and column name cannot be empty")
	}

	databaseType := database.DatabaseType(db)

	sqlString, sqlParams, err := sb.NewBuilder(databaseType).Table(tableName).ColumnExists(columnName)
	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	var exists bool
	err = db.QueryRowContext(ctx, sqlString, sqlParams...).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check column existence: %w", err)
	}

	return exists, nil
}
