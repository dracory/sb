package sb

import (
	"errors"

	"github.com/dracory/database"
)

func TableDropSql(ctx database.QueryableContext, tableName string) (string, error) {
	if ctx.Queryable() == nil {
		return "", errors.New("queryable cannot be nil")
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	return NewBuilder(databaseType).Table(tableName).Drop(), nil
}

func TableDropIfExistsSql(ctx database.QueryableContext, tableName string) (string, error) {
	if ctx.Queryable() == nil {
		return "", errors.New("queryable cannot be nil")
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	return NewBuilder(databaseType).Table(tableName).DropIfExists(), nil
}

func TableDrop(ctx database.QueryableContext, tableName string) error {
	sqlTableDrop, err := TableDropSql(ctx, tableName)

	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlTableDrop)

	return err
}

func TableDropIfExists(ctx database.QueryableContext, tableName string) error {
	if ctx.Queryable() == nil {
		return errors.New("queryable cannot be nil")
	}

	sqlTableDrop, err := TableDropIfExistsSql(ctx, tableName)

	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlTableDrop)

	return err
}
