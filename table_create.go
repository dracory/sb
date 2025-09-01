package sb

import (
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

func TableCreate(db *sql.DB, tableName string, columns []Column) error {
	databaseType := database.DatabaseType(db)

	builder := NewBuilder(databaseType).Table(tableName)

	for _, column := range columns {
		builder.Column(column)
	}

	sqlTable := builder.Create()

	_, err := db.Exec(sqlTable)

	return err
}
