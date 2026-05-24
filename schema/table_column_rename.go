package schema

import (
	"github.com/dracory/database"
	"github.com/dracory/sb"
)

// TableColumnRename renames a column in a table.
func TableColumnRename(ctx database.QueryableContext, tableName string, oldColumnName string, newColumnName string) error {
	if ctx.Queryable() == nil {
		return sb.ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sqlStr, err := sb.NewBuilder(databaseType).TableColumnRename(tableName, oldColumnName, newColumnName)
	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlStr)

	return err
}
