package sb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dracory/database"
)

// TableColumnExists checks if a column exists in a table for various database types.
func TableColumnExists(ctx database.QueryableContext, tableName string, columnName string) (bool, error) {
	db := ctx.Queryable()
	if db == nil {
		return false, ErrNilQueryable
	}

	if tableName == "" || columnName == "" {
		return false, NewValidationError("table name and column name cannot be empty")
	}

	databaseType := database.DatabaseType(db)

	builder := NewBuilder(databaseType).Table(tableName)
	sqlString, sqlParams, err := builder.TableColumnExists(tableName, columnName)

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
