package schema

import (
	"errors"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

// TableDrop drops a table from the database.
func TableDrop(ctx database.QueryableContext, tableName string) error {
	if ctx.Queryable() == nil {
		return errors.New("queryable cannot be nil")
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).Table(tableName).Drop()
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	return err
}

// TableDropIfExists drops a table from the database if it exists.
func TableDropIfExists(ctx database.QueryableContext, tableName string) error {
	if ctx.Queryable() == nil {
		return errors.New("queryable cannot be nil")
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).Table(tableName).DropIfExists()
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	return err
}
