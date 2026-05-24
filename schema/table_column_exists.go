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
		return false, errors.New("queryable cannot be nil")
	}

	if tableName == "" || columnName == "" {
		return false, fmt.Errorf("table name and column name cannot be empty")
	}

	databaseType := database.DatabaseType(db)

	sqlString, sqlParams, err := tableColumnExistsSql(databaseType, tableName, columnName)
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

// tableColumnExistsSql returns the SQL query for checking if a column exists.
func tableColumnExistsSql(dialect, tableName, columnName string) (string, []interface{}, error) {
	switch dialect {
	case sb.DIALECT_MYSQL:
		return "SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = ?", []interface{}{tableName, columnName}, nil
	case sb.DIALECT_POSTGRES:
		return "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", []interface{}{tableName, columnName}, nil
	case sb.DIALECT_SQLITE:
		return "SELECT 1 FROM pragma_table_info(?) WHERE name = ?", []interface{}{tableName, columnName}, nil
	default:
		return "", nil, fmt.Errorf("database type '%s' not supported", dialect)
	}
}
