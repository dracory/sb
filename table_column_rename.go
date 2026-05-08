package sb

import (
	"github.com/dracory/database"
)

// TableColumnRenameSql generates the SQL for renaming a column in a table
func TableColumnRenameSql(ctx database.QueryableContext, tableName string, oldColumnName string, newColumnName string) (string, error) {
	if ctx.Queryable() == nil {
		return "", ErrNilQueryable
	}

	databaseType := database.DatabaseType(ctx.Queryable())

	sql, err := NewBuilder(databaseType).TableColumnRename(tableName, oldColumnName, newColumnName)
	if err != nil {
		return "", err
	}
	return sql, nil
}

// TableColumnRename executes the SQL to rename a column in a table
func TableColumnRename(ctx database.QueryableContext, tableName string, oldColumnName string, newColumnName string) error {
	sqlColumnRename, err := TableColumnRenameSql(ctx, tableName, oldColumnName, newColumnName)

	if err != nil {
		return err
	}

	_, err = ctx.Queryable().ExecContext(ctx, sqlColumnRename)

	return err
}
