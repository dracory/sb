package sb

import (
	"context"
	"errors"
	"strings"

	"github.com/gouniverse/base/database"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// TableColumns returns a list of columns for a given table name
func TableColumns(ctx context.Context, q database.Queryable, tableName string, commonize bool) (columns []Column, err error) {
	databaseType := database.DatabaseType(q)

	if strings.TrimSpace(tableName) == "" {
		return nil, errors.New("table name cannot be empty")
	}

	if databaseType == "" {
		return nil, errors.New("driver type cannot be empty")
	}

	if databaseType == database.DATABASE_TYPE_SQLITE {
		return tableColumnsSqlite(ctx, q, tableName, commonize)
	}

	if databaseType == database.DATABASE_TYPE_MYSQL {
		return tableColumnsMysql(ctx, q, tableName, commonize)
	}

	return columns, errors.New("not implemented for database driver: " + databaseType)
}

func tableColumnsMysql(ctx context.Context, q database.Queryable, tableName string, commonize bool) (columns []Column, err error) {
	sql := "DESCRIBE `" + tableName + "`;"

	rows, err := database.SelectToMapString(ctx, q, sql)

	if err != nil {
		return columns, err
	}

	for _, row := range rows {
		columnName := lo.ValueOr(row, "Field", "")
		columnTypeRaw := strings.ToLower(lo.ValueOr(row, "Type", ""))
		columnDefault := lo.ValueOr(row, "Default", "")
		columnNullable := lo.ValueOr(row, "Null", "")
		columnKey := lo.ValueOr(row, "Key", "")
		columnExtra := lo.ValueOr(row, "Extra", "")

		columnType, length, decimals := rawColumnProcess(columnTypeRaw)

		if commonize {
			if strings.Contains(columnType, "int") {
				columnType = COLUMN_TYPE_INTEGER
			} else if strings.Contains(columnType, "char") {
				columnType = COLUMN_TYPE_STRING
			} else if strings.Contains(columnType, "text") {
				columnType = COLUMN_TYPE_TEXT
			} else if strings.Contains(columnType, "float") {
				columnType = COLUMN_TYPE_FLOAT
			} else if strings.Contains(columnType, "blob") {
				columnType = COLUMN_TYPE_BLOB
			}
		}

		isPrimaryKey := columnKey == "PRI"
		isUnique := columnKey == "UNI"
		isNullable := columnNullable == "YES"
		isAutoIncrement := strings.Contains(columnExtra, "auto_increment")

		column := Column{
			Name:          columnName,
			Type:          columnType,
			PrimaryKey:    isPrimaryKey,
			Unique:        isUnique,
			Nullable:      isNullable,
			Default:       columnDefault,
			AutoIncrement: isAutoIncrement,
		}

		if length != "" && isNumeric(length) {
			column.Length = cast.ToInt(length)
		}

		if decimals != "" && isNumeric(decimals) {
			column.Decimals = cast.ToInt(decimals)
		}

		columns = append(columns, column)
	}

	return columns, nil
}

func tableColumnsSqlite(ctx context.Context, q database.Queryable, tableName string, commonize bool) (columns []Column, err error) {
	sql := "SELECT * FROM 'SQLITE_MASTER' WHERE type='table' ORDER BY NAME ASC;"
	sql += "PRAGMA table_info('" + tableName + "');"

	rows, err := database.SelectToMapString(ctx, q, sql)

	if err != nil {
		return columns, err
	}

	for _, row := range rows {
		columnName := lo.ValueOr(row, "name", "")
		columnTypeRaw := strings.ToLower(lo.ValueOr(row, "type", ""))
		columnDefault := lo.ValueOr(row, "dflt_value", "")
		columnNotNull := lo.ValueOr(row, "notnull", "")
		columnPrimaryKey := lo.ValueOr(row, "pk", "")

		columnType, length, decimals := rawColumnProcess(columnTypeRaw)

		if commonize {
			if strings.Contains(columnType, "int") {
				columnType = COLUMN_TYPE_INTEGER
			} else if strings.Contains(columnType, "char") {
				columnType = COLUMN_TYPE_STRING
			} else if strings.Contains(columnType, "text") {
				// sqlite would return text even for varchar,
				// which is why we check the length, to see
				// if its string field
				lengthInt := cast.ToInt(length)
				if lengthInt > 0 {
					columnType = COLUMN_TYPE_STRING
				} else {
					columnType = COLUMN_TYPE_TEXT
				}
			} else if strings.Contains(columnType, "real") {
				columnType = COLUMN_TYPE_FLOAT
			} else if strings.Contains(columnType, "blob") {
				columnType = COLUMN_TYPE_BLOB
			}
		}

		isPrimaryKey := columnPrimaryKey == "1"
		isNullable := columnNotNull == "0"
		isAutoIncrement := false

		if columnType == "INTEGER" && isPrimaryKey {
			isAutoIncrement = true
		}

		column := Column{
			Name:          columnName,
			Type:          columnType,
			PrimaryKey:    isPrimaryKey,
			Nullable:      isNullable,
			Default:       columnDefault,
			AutoIncrement: isAutoIncrement,
		}

		if length != "" && isNumeric(length) {
			column.Length = cast.ToInt(length)
		}

		if decimals != "" && isNumeric(decimals) {
			column.Decimals = cast.ToInt(decimals)
		}

		columns = append(columns, column)
	}

	return columns, nil
}

func rawColumnProcess(columnType string) (scolumnType string, length string, decimals string) {
	if !strings.Contains(columnType, "(") {
		return columnType, "", ""
	}

	splitByParen := strings.Split(columnType, "(")
	columnType = strings.TrimSpace(splitByParen[0])
	props := strings.TrimRight(splitByParen[1], ")")

	if strings.Contains(props, ",") {
		splitByComma := strings.Split(props, ",")
		length = strings.TrimSpace(splitByComma[0])
		decimals = strings.TrimSpace(splitByComma[1])
	} else {
		length = strings.TrimSpace(props)
	}

	return columnType, length, decimals
}

func isNumeric(str string) bool {
	for _, char := range str {
		if !(char >= '0' && char <= '9' || char == '.') {
			return false
		}
	}
	return true
}
