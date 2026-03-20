package sb

import (
	"github.com/samber/lo"
)

type ColumnSQLGenerator interface {
	GenerateSQL(column Column) string
}

type MySQLColumnSQLGenerator struct{}

func (g MySQLColumnSQLGenerator) GenerateSQL(column Column) string {
	columnLength := column.Length
	columnType := lo.
		IfF(column.Type == COLUMN_TYPE_STRING, func() string {
			columnLength = lo.Ternary(columnLength == 0, 255, columnLength)
			return "VARCHAR"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_INTEGER, func() string {
			columnLength = lo.Ternary(columnLength == 0, 20, columnLength)
			return "BIGINT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_FLOAT, func() string {
			return "DOUBLE"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_TEXT, func() string {
			return "LONGTEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_LONGTEXT, func() string {
			return "LONGTEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_BLOB, func() string {
			return "LONGBLOB"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATE, func() string {
			return "DATE"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATETIME, func() string {
			return "DATETIME"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DECIMAL, func() string {
			return "DECIMAL"
		}).
		Else(column.Type)

	sql := "`" + column.Name + "` " + columnType

	// Column length
	if columnType == "DECIMAL" {
		if columnLength == 0 {
			columnLength = 10
		}
		if column.Decimals == 0 {
			column.Decimals = 2
		}
		sql += "(" + toString(columnLength) + "," + toString(column.Decimals) + ")"

	} else if columnLength != 0 {
		sql += "(" + toString(columnLength) + ")"
	}

	// Auto increment
	if column.AutoIncrement {
		sql += " AUTO_INCREMENT"
	}

	// Primary key
	if column.PrimaryKey {
		sql += " PRIMARY KEY"
	}

	// Non Nullable / Required
	if !column.Nullable {
		sql += " NOT NULL"
	}

	if column.Unique {
		sql += " UNIQUE"
	}
	return sql
}

type PostgreSQLColumnSQLGenerator struct{}

func (g PostgreSQLColumnSQLGenerator) GenerateSQL(column Column) string {
	// For PostgreSQL, SERIAL replaces INTEGER when AutoIncrement is true
	columnType := lo.
		IfF(column.Type == COLUMN_TYPE_STRING, func() string {
			return "TEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_INTEGER, func() string {
			if column.AutoIncrement {
				return "SERIAL"
			}
			return "INTEGER"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_FLOAT, func() string {
			return "REAL"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_TEXT, func() string {
			return "TEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_LONGTEXT, func() string {
			return "TEXT" // PostgreSQL only has TEXT (which defaults to LONGTEXT)
		}).
		ElseIfF(column.Type == COLUMN_TYPE_BLOB, func() string {
			return "BYTEA"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATE, func() string {
			return "DATE"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATETIME, func() string {
			return "TIMESTAMP"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DECIMAL, func() string {
			return "DECIMAL"
		}).
		Else(column.Type)

	sql := `"` + column.Name + `" ` + columnType + ``

	// Column length
	if columnType == "DECIMAL" {
		if column.Length == 0 {
			column.Length = 10
		}
		if column.Decimals == 0 {
			column.Decimals = 2
		}
		sql += "(" + toString(column.Length) + "," + toString(column.Decimals) + ")"

	} else if column.Length != 0 {
		supportsLength := columnType == "VARCHAR" || columnType == "CHAR" || columnType == "BIT VARYING" || columnType == "BIT"
		if supportsLength {
			sql += "(" + toString(column.Length) + ")"
		}
	}

	// Auto increment is already handled in columnType for SERIAL
	// No need to append " SERIAL" here

	// Primary key
	if column.PrimaryKey {
		sql += " PRIMARY KEY"
	}

	// Non Nullable / Required
	if !column.Nullable {
		sql += " NOT NULL"
	}

	if column.Unique {
		sql += " UNIQUE"
	}
	return sql
}

type SQLiteColumnSQLGenerator struct{}

func (g SQLiteColumnSQLGenerator) GenerateSQL(column Column) string {
	columnType := lo.
		IfF(column.Type == COLUMN_TYPE_STRING, func() string {
			return "TEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_INTEGER, func() string {
			return "INTEGER"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_FLOAT, func() string {
			return "REAL"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_TEXT, func() string {
			return "TEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_LONGTEXT, func() string {
			return "TEXT" // SQLite only has TEXT (which defaults to LONGTEXT)
		}).
		ElseIfF(column.Type == COLUMN_TYPE_BLOB, func() string {
			return "BLOB"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATE, func() string {
			return "DATE"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATETIME, func() string {
			return "DATETIME"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DECIMAL, func() string {
			return "DECIMAL"
		}).
		Else(column.Type)

	sql := `"` + column.Name + `" ` + columnType + ``

	// Column length
	if columnType == "DECIMAL" {
		if column.Length == 0 {
			column.Length = 10
		}
		if column.Decimals == 0 {
			column.Decimals = 2
		}
		sql += "(" + toString(column.Length) + "," + toString(column.Decimals) + ")"

	} else if column.Length != 0 {
		sql += "(" + toString(column.Length) + ")"
	}

	// Primary key
	if column.PrimaryKey {
		sql += " PRIMARY KEY"
	}

	// Auto increment
	if column.AutoIncrement {
		sql += " AUTOINCREMENT"
	}

	// Non Nullable / Required
	if !column.Nullable {
		sql += " NOT NULL"
	}

	if column.Unique {
		sql += " UNIQUE"
	}

	return sql
}

type MSSQLColumnSQLGenerator struct{}

func (g MSSQLColumnSQLGenerator) GenerateSQL(column Column) string {
	columnType := lo.
		IfF(column.Type == COLUMN_TYPE_STRING, func() string {
			return "NVARCHAR"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_INTEGER, func() string {
			return "INTEGER"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_FLOAT, func() string {
			return "FLOAT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_TEXT, func() string {
			return "TEXT"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_BLOB, func() string {
			return "VARBINARY"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATE, func() string {
			return "DATE"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DATETIME, func() string {
			return "DATETIME2"
		}).
		ElseIfF(column.Type == COLUMN_TYPE_DECIMAL, func() string {
			return "DECIMAL"
		}).
		Else(column.Type)

	sql := `"` + column.Name + `" ` + columnType + ``

	// Column length
	if columnType == "DECIMAL" {
		if column.Length == 0 {
			column.Length = 10
		}
		if column.Decimals == 0 {
			column.Decimals = 2
		}
		sql += "(" + toString(column.Length) + "," + toString(column.Decimals) + ")"

	} else if columnType == "VARBINARY" {
		sql += "(MAX)"

	} else if column.Length != 0 {
		sql += "(" + toString(column.Length) + ")"
	}

	// Auto increment
	if column.AutoIncrement {
		sql += " AUTOINCREMENT"
	}

	// Primary key
	if column.PrimaryKey {
		sql += " PRIMARY KEY"
	}

	// Non Nullable / Required
	if !column.Nullable {
		sql += " NOT NULL"
	}

	if column.Unique {
		sql += " UNIQUE"
	}

	return sql
}
