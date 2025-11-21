package sb

import (
	"context"
	"database/sql"

	"github.com/dracory/database"
)

func TableCreateSql(db *sql.DB, tableName string, columns []Column) string {
	databaseType := database.DatabaseType(db)

	builder := NewBuilder(databaseType).Table(tableName)

	for _, column := range columns {
		builder.Column(column)
	}

	return builder.Create()
}

func TableCreate(ctx context.Context, db *sql.DB, tableName string, columns []Column) error {
	databaseType := database.DatabaseType(db)

	builder := NewBuilder(databaseType).Table(tableName)

	for _, column := range columns {
		builder.Column(column)
	}

	sqlTable := builder.Create()

	if ctx == nil {
		ctx = context.Background()
	}

	_, err := db.ExecContext(ctx, sqlTable)

	return err
}
