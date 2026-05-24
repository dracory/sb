package schema

import (
	"context"
	"database/sql"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

// TableCreate creates a table with the given columns.
func TableCreate(ctx context.Context, db *sql.DB, tableName string, columns []sb.Column) error {
	databaseType := database.DatabaseType(db)

	builder := sb.NewBuilder(databaseType).Table(tableName)

	for _, column := range columns {
		builder.Column(column)
	}

	sqlTable, err := builder.Create()
	if err != nil {
		return err
	}

	if ctx == nil {
		ctx = context.Background()
	}

	_, err = db.ExecContext(ctx, sqlTable)

	return err
}
