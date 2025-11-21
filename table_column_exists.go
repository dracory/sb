package sb

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/dracory/database"
)

// TableColumnExists checks if a column exists in a table for various database types.
func TableColumnExists(ctx database.QueryableContext, tableName, columnName string) (exists bool, err error) {
	db := ctx.Queryable()
	if db == nil {
		return false, errors.New("queryable cannot be nil")
	}

	if tableName == "" || columnName == "" {
		return false, errors.New("table name and column name cannot be empty")
	}

	databaseType := database.DatabaseType(db)

	builder := NewBuilder(databaseType).Table(tableName)
	sqlString, sqlParams, err := builder.TableColumnExists(tableName, columnName)

	if err != nil {
		return false, fmt.Errorf("failed to build query: %w", err)
	}

	err = db.QueryRowContext(ctx, sqlString, sqlParams...).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check column existence: %w", err)
	}

	return exists, nil
}
