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

	// databaseType := database.DatabaseType(db)
	// switch databaseType {
	// case database.DATABASE_TYPE_MYSQL:
	// 	err = db.QueryRowContext(ctx, "SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = ?", tableName, columnName).Scan(&exists)
	// case database.DATABASE_TYPE_POSTGRES:
	// 	err = db.QueryRowContext(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", tableName, columnName).Scan(&exists)
	// case database.DATABASE_TYPE_SQLITE:
	// 	err = db.QueryRowContext(ctx, "SELECT 1 FROM pragma_table_info(?) WHERE name = ?", tableName, columnName).Scan(&exists)
	// default:
	// 	return false, fmt.Errorf("database type '%s' not supported", databaseType)
	// }

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check column existence: %w", err)
	}

	return exists, nil
}
