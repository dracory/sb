package sb

import (
	"context"
	"database/sql"
)

type BuilderInterface interface {
	// Column adds a column to the table
	Column(column Column) BuilderInterface

	// Create creates a table
	Create() string

	// CreateIfNotExists creates a table if it doesn't exist
	CreateIfNotExists() string

	// CreateIndex creates an index on the table
	CreateIndex(indexName string, columnName ...string) string

	// Delete deletes a table
	Delete() string

	// Drop drops a table
	Drop() string

	// DropIfExists drops a table if it exists
	DropIfExists() string

	// Insert inserts a row into the table
	Insert(columnValuesMap map[string]string) string

	// GroupBy groups the results by a column
	GroupBy(groupBy GroupBy) BuilderInterface

	// Limit limits the number of results
	Limit(limit int64) BuilderInterface

	// Offset offsets the results
	Offset(offset int64) BuilderInterface

	// OrderBy orders the results by a column
	OrderBy(columnName string, sortDirection string) BuilderInterface

	// Select selects the columns from the table
	Select(columns []string) string

	// Table sets the table name
	Table(name string) BuilderInterface

	// Update updates a row in the table
	Update(columnValues map[string]string) string

	// View sets the view name
	View(name string) BuilderInterface

	// ViewColumns sets the view columns
	ViewColumns(columns []string) BuilderInterface

	// ViewSQL sets the view SQL
	ViewSQL(sql string) BuilderInterface

	// Where sets the where clause
	Where(where Where) BuilderInterface

	// TableColumnAdd adds a column to the table
	TableColumnAdd(tableName string, column Column) (sqlString string, err error)

	// TableColumnChange changes a column in the table
	TableColumnChange(tableName string, column Column) (sqlString string, err error)

	// Table column drop drops a column
	TableColumnDrop(tableName string, columnName string) (sqlString string, err error)

	// TableColumnExists checks if a column exists in a table
	TableColumnExists(tableName, columnName string) (sqlString string, sqlParams []any, err error)

	// TableColumnRename renames a column in a table
	TableColumnRename(tableName, oldColumnName, newColumnName string) (sqlString string, err error)

	// TableRename renames a table
	TableRename(oldTableName string, newTableName string) (sqlString string, err error)
}

type DatabaseInterface interface {
	// DB the database connection
	DB() *sql.DB

	// Type the database type, i.e. "mssql", "mysql", "postgres", "sqlite"
	Type() string

	// BeginTransaction starts a transaction
	BeginTransaction() (err error)

	// BeginTransactionWithContext starts a transaction with context
	BeginTransactionWithContext(ctx context.Context, opts *sql.TxOptions) (err error)

	// Close closes the database
	Close() (err error)

	// CommitTransaction commits the transaction
	CommitTransaction() (err error)

	// DebugEnable enables or disables debug
	DebugEnable(debug bool)

	// ExecInTransaction executes a function in a transaction
	ExecInTransaction(fn func(d *Database) error) (err error)

	// Exec executes a query
	Exec(sqlStr string, args ...any) (sql.Result, error)

	IsMssql() bool
	IsMysql() bool
	IsPostgres() bool
	IsSqlite() bool
	SqlLog() []map[string]string
	SqlLogEmpty()
	SqlLogLen() int
	SqlLogEnable(enable bool)
	SqlLogShrink(leaveLast int)
	Open() (err error)
	Query(sqlStr string, args ...any) (*sql.Rows, error)
	RollbackTransaction() (err error)
	SelectToMapAny(sqlStr string, args ...any) ([]map[string]any, error)
	SelectToMapString(sqlStr string, args ...any) ([]map[string]string, error)

	// Tx the transaction
	Tx() *sql.Tx
}
