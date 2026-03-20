package sb

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type OrderBy struct {
	Column    string
	Direction string
}

type Column struct {
	Name          string
	Type          string
	Length        int
	Decimals      int
	AutoIncrement bool
	PrimaryKey    bool
	Nullable      bool
	Unique        bool
	Default       string
}

type GroupBy struct {
	Column string
}

// TruncateOptions provides database-specific options for table truncation
type TruncateOptions struct {
	Cascade       bool // For PostgreSQL: adds CASCADE to handle foreign key constraints
	ResetIdentity bool // For MSSQL: resets identity column seed value after truncation
}

// IndexColumn defines a single column in an index, with optional direction and prefix length.
type IndexColumn struct {
	Name      string // column name
	Direction string // "ASC" or "DESC" — defaults to ASC if empty
	Length    int    // prefix length for MySQL text/blob columns (ignored on other dialects)
}

// IndexOptions controls advanced index creation behaviour.
type IndexOptions struct {
	Unique      bool          // emit UNIQUE keyword
	IfNotExists bool          // emit IF NOT EXISTS where supported
	Columns     []IndexColumn // columns with optional direction / prefix length
	Using       string        // index type: BTREE, HASH, GIN, GIST, BRIN, FULLTEXT, SPATIAL
	Include     []string      // PostgreSQL / MSSQL covering-index columns (INCLUDE clause)
	Where       string        // partial index predicate (PostgreSQL, SQLite, MSSQL)
	Storage     string        // PostgreSQL WITH (...) storage parameters, e.g. "fillfactor=90"
	Comment     string        // MySQL COMMENT on index
}

// DropIndexOptions controls advanced index removal behaviour.
type DropIndexOptions struct {
	IfExists bool   // emit IF EXISTS where supported
	Schema   string // schema-qualify the index name (PostgreSQL only)
}

// JoinType represents the type of JOIN operation
type JoinType string

// Join represents a database JOIN operation
type Join struct {
	Type        JoinType
	Table       string
	Alias       string
	OnCondition string
	Columns     []string // Optional: for specific column selection
}

type Builder struct {
	Dialect            string
	sql                map[string]any
	sqlColumns         []Column
	sqlErrors          []error // Collect errors during fluent chaining
	sqlGroupBy         []GroupBy
	sqlJoins           []Join
	sqlLimit           int64
	sqlOffset          int64
	sqlOrderBy         []OrderBy
	sqlSelectColumns   []string
	sqlTableName       string
	sqlViewName        string
	sqlViewColumns     []string
	sqlViewSQL         string
	sqlWhere           []Where
	columnSQLGenerator ColumnSQLGenerator
	params             []interface{} // Track parameters for parameterized queries
	paramIndex         int           // Track parameter index for naming
	interpolatedValues bool          // Track if using interpolated values mode
}

var _ BuilderInterface = (*Builder)(nil)

func NewBuilder(dialect string) *Builder {
	var columnSQLGenerator ColumnSQLGenerator
	switch dialect {
	case DIALECT_MYSQL:
		columnSQLGenerator = MySQLColumnSQLGenerator{}
	case DIALECT_POSTGRES:
		columnSQLGenerator = PostgreSQLColumnSQLGenerator{}
	case DIALECT_SQLITE:
		columnSQLGenerator = SQLiteColumnSQLGenerator{}
	case DIALECT_MSSQL:
		columnSQLGenerator = MSSQLColumnSQLGenerator{}
	default:
		// Create builder with error collected for unsupported dialect
		return &Builder{
			Dialect:            dialect,
			sql:                map[string]any{},
			sqlColumns:         []Column{},
			sqlErrors:          []error{NewValidationError("unsupported dialect: " + dialect)}, // Collect error
			sqlGroupBy:         []GroupBy{},
			sqlJoins:           []Join{},
			sqlLimit:           0,
			sqlOffset:          0,
			sqlOrderBy:         []OrderBy{},
			sqlTableName:       "",
			sqlViewName:        "",
			sqlViewColumns:     []string{},
			sqlViewSQL:         "",
			sqlWhere:           []Where{},
			columnSQLGenerator: nil, // No generator for unsupported dialect
		}
	}

	return &Builder{
		Dialect:            dialect,
		sql:                map[string]any{},
		sqlColumns:         []Column{},
		sqlErrors:          []error{}, // Initialize error collection
		sqlGroupBy:         []GroupBy{},
		sqlJoins:           []Join{},
		sqlLimit:           0,
		sqlOffset:          0,
		sqlOrderBy:         []OrderBy{},
		sqlTableName:       "",
		sqlViewName:        "",
		sqlViewColumns:     []string{},
		sqlViewSQL:         "",
		sqlWhere:           []Where{},
		columnSQLGenerator: columnSQLGenerator,
	}
}

func (b *Builder) Table(tableName string) BuilderInterface {
	b.sqlTableName = tableName
	return b
}

func (b *Builder) View(viewName string) BuilderInterface {
	b.sqlViewName = viewName
	return b
}

func (b *Builder) ViewSQL(sql string) BuilderInterface {
	b.sqlViewSQL = sql
	return b
}

func (b *Builder) ViewColumns(columns []string) BuilderInterface {
	b.sqlViewColumns = columns
	return b
}

func (b *Builder) Column(column Column) BuilderInterface {
	if column.Name == "" {
		b.sqlErrors = append(b.sqlErrors, ErrEmptyColumnName)
		return b
	}

	if column.Type == "" {
		b.sqlErrors = append(b.sqlErrors, NewValidationError("column type is required"))
		return b
	}

	b.sqlColumns = append(b.sqlColumns, column)
	return b
}

// addError adds an error to the error collection
func (b *Builder) addError(err error) {
	if err != nil {
		b.sqlErrors = append(b.sqlErrors, err)
	}
}

// hasErrors returns true if there are collected errors
func (b *Builder) hasErrors() bool {
	return len(b.sqlErrors) > 0
}

// getErrors returns all collected errors
func (b *Builder) getErrors() []error {
	return b.sqlErrors
}

// validateAndReturnError checks for collected errors and returns them as a single error
func (b *Builder) validateAndReturnError() error {
	if len(b.sqlErrors) == 0 {
		return nil
	}

	// Return the first error for simplicity, could also combine multiple errors
	return b.sqlErrors[0]
}

// addParam adds a parameter value and returns the appropriate placeholder for the current dialect
func (b *Builder) addParam(value interface{}) string {
	b.params = append(b.params, value)
	placeholder := b.getParameterPlaceholder()
	b.paramIndex++
	return placeholder
}

// getParameterPlaceholder returns the parameter placeholder syntax for the current dialect
func (b *Builder) getParameterPlaceholder() string {
	switch b.Dialect {
	case DIALECT_MYSQL:
		return "?"
	case DIALECT_POSTGRES:
		return fmt.Sprintf("$%d", b.paramIndex+1)
	case DIALECT_SQLITE:
		return "?"
	case DIALECT_MSSQL:
		return fmt.Sprintf("@p%d", b.paramIndex+1)
	default:
		return "?"
	}
}

// resetParams resets the parameter tracking for a new query
func (b *Builder) resetParams() {
	b.params = []interface{}{}
	b.paramIndex = 0
}

// WithInterpolatedValues enables legacy mode with interpolated values instead of parameterized queries.
// WARNING: This mode is insecure and vulnerable to SQL injection. Use only for backward compatibility.
func (b *Builder) WithInterpolatedValues() BuilderInterface {
	b.interpolatedValues = true
	return b
}

/**
 * The create method creates new database or table.
 * If the database or table can not be created it will return false.
 * False will be returned if the database or table already exist.
 * <code>
 * // Creating a new database
 * $database->create();
 *
 * // Creating a new table
 * $database->table("STATES")
 *     ->column("STATE_NAME","STRING")
 *     ->create();
 * </code>
 * @return boolean true, on success, false, otherwise
 * @access public
 */
func (b *Builder) Create() (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	if !isTable && !isView {
		return "", ErrMissingTable
	}

	sql := ""

	if isTable {
		if len(b.sqlColumns) == 0 {
			return "", ErrEmptyColumns
		}
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = `CREATE TABLE ` + b.quoteTable(b.sqlTableName) + `(` + b.columnsToSQL(b.sqlColumns) + `);`
		}
		if b.Dialect == DIALECT_MSSQL {
			sql = `CREATE TABLE [` + b.sqlTableName + `] (` + b.columnsToSQL(b.sqlColumns) + `);`
		}
	}

	if isView {
		if b.sqlViewSQL == "" {
			return "", NewValidationError("view SQL cannot be empty")
		}
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			viewColumnsToSQL := strings.Join(lo.Map(b.sqlViewColumns, func(columnName string, _ int) string {
				return b.quoteColumn(columnName)
			}), ", ")
			viewColumns := lo.If(len(b.sqlViewColumns) > 0, ` (`+viewColumnsToSQL+`)`).Else(``)

			sql = `CREATE VIEW ` + b.quoteTable(b.sqlViewName) + viewColumns + " AS " + b.sqlViewSQL
		}
	}

	return sql, nil
}

func (b *Builder) CreateIfNotExists() (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	if !isTable && !isView {
		return "", ErrMissingTable
	}

	sql := ""

	if isTable {
		if len(b.sqlColumns) == 0 {
			return "", ErrEmptyColumns
		}
		if b.Dialect == DIALECT_MYSQL {
			sql = "CREATE TABLE IF NOT EXISTS " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
		}
		if b.Dialect == DIALECT_POSTGRES {
			sql = `CREATE TABLE IF NOT EXISTS ` + b.quoteTable(b.sqlTableName) + `(` + b.columnsToSQL(b.sqlColumns) + `);`
		}
		if b.Dialect == DIALECT_SQLITE {
			sql = "CREATE TABLE IF NOT EXISTS " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
		}
	}

	if isView {
		if b.sqlViewSQL == "" {
			return "", NewValidationError("view SQL cannot be empty")
		}
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			viewColumnsToSQL := strings.Join(lo.Map(b.sqlViewColumns, func(columnName string, _ int) string {
				return b.quoteColumn(columnName)
			}), ", ")
			viewColumns := lo.If(len(b.sqlViewColumns) > 0, ` (`+viewColumnsToSQL+`)`).Else(``)

			sqlStart := "CREATE VIEW IF NOT EXISTS"
			if b.Dialect == DIALECT_MYSQL {
				sqlStart = "CREATE OR REPLACE VIEW"
			}

			sql = sqlStart + ` ` + b.quoteTable(b.sqlViewName) + viewColumns + " AS " + b.sqlViewSQL
		}
	}

	return sql, nil
}

func (b *Builder) CreateIndex(indexName string, columnName ...string) (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	if indexName == "" {
		return "", ErrEmptyIndexName
	}
	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}

	columns := lo.Map(columnName, func(columnName string, i int) string {
		return b.quoteColumn(columnName)
	})

	if len(columns) == 0 {
		return "", ErrEmptyColumns
	}

	sql := `CREATE INDEX ` + b.quoteTable(indexName) + ` ON ` + b.quoteTable(b.sqlTableName) + ` (` + strings.Join(columns, `,`) + `);`

	return sql, nil
}

// CreateIndexWithOptions generates a CREATE INDEX statement with full dialect-aware options.
//
// Example (PostgreSQL partial unique index):
//
//	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
//	    Table("users").
//	    CreateIndexWithOptions("idx_users_active_email", sb.IndexOptions{
//	        Unique: true,
//	        Columns: []sb.IndexColumn{{Name: "email"}},
//	        Where:  "deleted_at IS NULL",
//	    })
//	// CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_active_email"
//	//   ON "users" ("email") WHERE deleted_at IS NULL;
func (b *Builder) CreateIndexWithOptions(name string, opts IndexOptions) (string, error) {
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}
	if name == "" {
		return "", ErrEmptyIndexName
	}
	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}
	if len(opts.Columns) == 0 {
		return "", ErrEmptyColumns
	}

	var sb strings.Builder

	sb.WriteString("CREATE ")
	if opts.Unique {
		sb.WriteString("UNIQUE ")
	}
	sb.WriteString("INDEX ")

	// IF NOT EXISTS — supported by PostgreSQL, SQLite; not by MySQL or MSSQL
	if opts.IfNotExists && (b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE) {
		sb.WriteString("IF NOT EXISTS ")
	}

	sb.WriteString(b.quoteTable(name))
	sb.WriteString(" ON ")
	sb.WriteString(b.quoteTable(b.sqlTableName))

	// USING clause — PostgreSQL only for standard index types; MySQL uses it for FULLTEXT/SPATIAL
	if opts.Using != "" {
		switch b.Dialect {
		case DIALECT_POSTGRES:
			sb.WriteString(" USING ")
			sb.WriteString(opts.Using)
		case DIALECT_MYSQL:
			// MySQL uses USING inside the column list for BTREE/HASH,
			// but FULLTEXT/SPATIAL are keywords before the column list.
			if opts.Using == INDEX_TYPE_FULLTEXT || opts.Using == INDEX_TYPE_SPATIAL {
				// handled below — rewrite the CREATE line
			}
		}
	}

	// MySQL FULLTEXT / SPATIAL require the keyword before the column list,
	// replacing "INDEX" entirely. Rebuild from scratch for these cases.
	if b.Dialect == DIALECT_MYSQL &&
		(opts.Using == INDEX_TYPE_FULLTEXT || opts.Using == INDEX_TYPE_SPATIAL) {
		sb.Reset()
		sb.WriteString("CREATE ")
		if opts.Unique {
			sb.WriteString("UNIQUE ")
		}
		sb.WriteString(opts.Using) // FULLTEXT or SPATIAL
		sb.WriteString(" INDEX ")
		sb.WriteString(b.quoteTable(name))
		sb.WriteString(" ON ")
		sb.WriteString(b.quoteTable(b.sqlTableName))
	}

	// Column list
	sb.WriteString(" (")
	for i, col := range opts.Columns {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(b.quoteColumn(col.Name))

		// MySQL prefix length must come immediately after the column name, before direction
		if b.Dialect == DIALECT_MYSQL && col.Length > 0 {
			sb.WriteString("(")
			sb.WriteString(strconv.Itoa(col.Length))
			sb.WriteString(")")
		}

		// Direction (not meaningful for FULLTEXT/SPATIAL but harmless to omit)
		// Only output direction if explicitly set (not empty and not default ASC)
		if col.Direction != "" && col.Direction != "ASC" &&
			!(b.Dialect == DIALECT_MYSQL &&
				(opts.Using == INDEX_TYPE_FULLTEXT || opts.Using == INDEX_TYPE_SPATIAL)) {
			sb.WriteString(" ")
			sb.WriteString(strings.ToUpper(col.Direction))
		}
	}

	// MySQL BTREE/HASH USING goes inside the column list parentheses
	if b.Dialect == DIALECT_MYSQL && opts.Using != "" &&
		opts.Using != INDEX_TYPE_FULLTEXT && opts.Using != INDEX_TYPE_SPATIAL {
		sb.WriteString(" USING ")
		sb.WriteString(opts.Using)
	}

	sb.WriteString(")")

	// INCLUDE clause — PostgreSQL 11+ and MSSQL
	if len(opts.Include) > 0 &&
		(b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_MSSQL) {
		sb.WriteString(" INCLUDE (")
		for i, col := range opts.Include {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(b.quoteColumn(col))
		}
		sb.WriteString(")")
	}

	// Partial index WHERE clause — PostgreSQL, SQLite, MSSQL
	if opts.Where != "" &&
		(b.Dialect == DIALECT_POSTGRES ||
			b.Dialect == DIALECT_SQLITE ||
			b.Dialect == DIALECT_MSSQL) {
		sb.WriteString(" WHERE ")
		sb.WriteString(opts.Where)
	}

	// PostgreSQL storage parameters
	if b.Dialect == DIALECT_POSTGRES && opts.Storage != "" {
		sb.WriteString(" WITH (")
		sb.WriteString(opts.Storage)
		sb.WriteString(")")
	}

	// MySQL index comment
	if b.Dialect == DIALECT_MYSQL && opts.Comment != "" {
		sb.WriteString(" COMMENT '")
		sb.WriteString(strings.ReplaceAll(opts.Comment, "'", "''"))
		sb.WriteString("'")
	}

	sb.WriteString(";")
	return sb.String(), nil
}

// CreateUniqueIndex creates a UNIQUE index on one or more columns.
func (b *Builder) CreateUniqueIndex(name string, columns ...string) (string, error) {
	return b.CreateIndexWithOptions(name, IndexOptions{
		Unique:  true,
		Columns: indexColumnsFromNames(columns),
	})
}

// CreateCompositeIndex creates an index on multiple columns with explicit ordering.
func (b *Builder) CreateCompositeIndex(name string, columns []IndexColumn) (string, error) {
	return b.CreateIndexWithOptions(name, IndexOptions{Columns: columns})
}

// CreatePartialIndex creates an index with a WHERE predicate (PostgreSQL, SQLite, MSSQL).
func (b *Builder) CreatePartialIndex(name string, where string, columns ...string) (string, error) {
	return b.CreateIndexWithOptions(name, IndexOptions{
		Columns: indexColumnsFromNames(columns),
		Where:   where,
	})
}

// CreateCoveringIndex creates a covering index using the INCLUDE clause (PostgreSQL, MSSQL).
func (b *Builder) CreateCoveringIndex(name string, include []string, columns ...string) (string, error) {
	return b.CreateIndexWithOptions(name, IndexOptions{
		Columns: indexColumnsFromNames(columns),
		Include: include,
	})
}

// indexColumnsFromNames is a package-level helper that converts plain column names
// to []IndexColumn with default ASC direction (empty string means default ASC).
func indexColumnsFromNames(names []string) []IndexColumn {
	cols := make([]IndexColumn, len(names))
	for i, n := range names {
		cols[i] = IndexColumn{Name: n, Direction: ""} // Empty string means default ASC
	}
	return cols
}

// DropIndex removes an index from a table.
// The method generates database-specific SQL:
//   - MySQL: DROP INDEX `index_name` ON `table_name`;
//   - PostgreSQL: DROP INDEX "index_name";
//   - SQLite: DROP INDEX "index_name";
//   - MSSQL: DROP INDEX [index_name] ON [table_name];
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").DropIndex("idx_users_email")
//	// Returns: "DROP INDEX `idx_users_email` ON `users`;"
func (b *Builder) DropIndex(indexName string) (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	if indexName == "" {
		return "", ErrEmptyIndexName
	}
	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}

	switch b.Dialect {
	case DIALECT_MYSQL:
		return "DROP INDEX " + b.quoteTable(indexName) + " ON " + b.quoteTable(b.sqlTableName) + ";", nil
	case DIALECT_POSTGRES:
		return "DROP INDEX " + b.quoteTable(indexName) + ";", nil
	case DIALECT_SQLITE:
		return "DROP INDEX " + b.quoteTable(indexName) + ";", nil
	case DIALECT_MSSQL:
		return "DROP INDEX " + b.quoteTable(indexName) + " ON " + b.quoteTable(b.sqlTableName) + ";", nil
	default:
		return "", ErrInvalidDialect
	}
}

// DropIndexIfExists removes an index from a table if it exists.
// Uses IF EXISTS clause where supported by the database dialect.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").DropIndexIfExists("idx_users_email")
//	// Returns: "DROP INDEX IF EXISTS \"idx_users_email\";"
func (b *Builder) DropIndexIfExists(indexName string) (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	if indexName == "" {
		return "", ErrEmptyIndexName
	}
	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}

	switch b.Dialect {
	case DIALECT_MYSQL:
		// MySQL doesn't support IF EXISTS for indexes, use regular DROP
		return "DROP INDEX " + b.quoteTable(indexName) + " ON " + b.quoteTable(b.sqlTableName) + ";", nil
	case DIALECT_POSTGRES:
		return "DROP INDEX IF EXISTS " + b.quoteTable(indexName) + ";", nil
	case DIALECT_SQLITE:
		return "DROP INDEX IF EXISTS " + b.quoteTable(indexName) + ";", nil
	case DIALECT_MSSQL:
		return "DROP INDEX IF EXISTS " + b.quoteTable(indexName) + " ON " + b.quoteTable(b.sqlTableName) + ";", nil
	default:
		return "", ErrInvalidDialect
	}
}

// DropIndexWithSchema removes an index from a specific schema (PostgreSQL only).
// For other dialects, falls back to regular DropIndex behavior.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").
//	  DropIndexWithSchema("idx_users_email", "public")
//	// Returns: "DROP INDEX IF EXISTS \"public\".\"idx_users_email\";"
func (b *Builder) DropIndexWithSchema(indexName string, schema string) (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	if indexName == "" {
		return "", ErrEmptyIndexName
	}
	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}

	switch b.Dialect {
	case DIALECT_POSTGRES:
		if schema != "" {
			return "DROP INDEX IF EXISTS " + b.quoteTable(schema) + "." + b.quoteTable(indexName) + ";", nil
		}
		return "DROP INDEX IF EXISTS " + b.quoteTable(indexName) + ";", nil
	default:
		// Other dialects don't support schema-qualified index names
		return b.DropIndex(indexName)
	}
}

// DropIndexWithOptions generates a DROP INDEX statement with full dialect-aware options.
func (b *Builder) DropIndexWithOptions(name string, opts DropIndexOptions) (string, error) {
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}
	if name == "" {
		return "", ErrEmptyIndexName
	}

	var sb strings.Builder
	sb.WriteString("DROP INDEX ")

	switch b.Dialect {
	case DIALECT_POSTGRES:
		if opts.IfExists {
			sb.WriteString("IF EXISTS ")
		}
		if opts.Schema != "" {
			sb.WriteString(b.quoteTable(opts.Schema))
			sb.WriteString(".")
		}
		sb.WriteString(b.quoteTable(name))

	case DIALECT_SQLITE:
		if opts.IfExists {
			sb.WriteString("IF EXISTS ")
		}
		sb.WriteString(b.quoteTable(name))

	case DIALECT_MSSQL:
		if opts.IfExists {
			sb.WriteString("IF EXISTS ")
		}
		sb.WriteString(b.quoteTable(name))
		if b.sqlTableName != "" {
			sb.WriteString(" ON ")
			sb.WriteString(b.quoteTable(b.sqlTableName))
		}

	case DIALECT_MYSQL:
		// MySQL has no IF EXISTS for DROP INDEX
		sb.WriteString(b.quoteTable(name))
		if b.sqlTableName != "" {
			sb.WriteString(" ON ")
			sb.WriteString(b.quoteTable(b.sqlTableName))
		}

	default:
		return "", ErrInvalidDialect
	}

	sb.WriteString(";")
	return sb.String(), nil
}

// Join adds a JOIN clause to the query.
// The joinType parameter specifies the type of join (INNER, LEFT, RIGHT, FULL, CROSS).
// The table parameter specifies the table to join.
// The onCondition parameter specifies the join condition.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).Table("orders").
//	  Join(sb.JOIN_INNER, "users", "orders.user_id = users.id").
//	  Select([]string{"orders.*", "users.name"})
//	// Returns: "SELECT orders.*, users.name FROM orders INNER JOIN users ON orders.user_id = users.id;"
func (b *Builder) Join(joinType JoinType, table string, onCondition string) BuilderInterface {
	if onCondition == "" {
		b.sqlErrors = append(b.sqlErrors, ErrEmptyOnCondition)
		return b
	}

	join := Join{
		Type:        joinType,
		Table:       table,
		OnCondition: onCondition,
	}

	b.sqlJoins = append(b.sqlJoins, join)
	return b
}

// JoinWithAlias adds a JOIN clause with table alias to the query.
// The joinType parameter specifies the type of join (INNER, LEFT, RIGHT, FULL, CROSS).
// The table parameter specifies the table to join.
// The alias parameter specifies the alias for the joined table.
// The onCondition parameter specifies the join condition.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("orders").
//	  JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
//	  Select([]string{"orders.*", "p.avatar"})
//	// Returns: "SELECT orders.*, p.avatar FROM orders LEFT JOIN profiles AS p ON orders.user_id = p.user_id;"
func (b *Builder) JoinWithAlias(joinType JoinType, table string, alias string, onCondition string) BuilderInterface {
	if onCondition == "" {
		b.sqlErrors = append(b.sqlErrors, ErrEmptyOnCondition)
		return b
	}

	join := Join{
		Type:        joinType,
		Table:       table,
		Alias:       alias,
		OnCondition: onCondition,
	}

	b.sqlJoins = append(b.sqlJoins, join)
	return b
}

// LeftJoin adds a LEFT JOIN clause to the query.
// The table parameter specifies the table to join.
// The onCondition parameter specifies the join condition.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_SQLITE).Table("orders").
//	  LeftJoin("users", "orders.user_id = users.id").
//	  Select([]string{"orders.*", "users.name"})
//	// Returns: "SELECT orders.*, users.name FROM orders LEFT JOIN users ON orders.user_id = users.id;"
func (b *Builder) LeftJoin(table string, onCondition string) BuilderInterface {
	return b.Join(JOIN_LEFT, table, onCondition)
}

// RightJoin adds a RIGHT JOIN clause to the query.
// The table parameter specifies the table to join.
// The onCondition parameter specifies the join condition.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_MSSQL).Table("orders").
//	  RightJoin("users", "orders.user_id = users.id").
//	  Select([]string{"orders.*", "users.name"})
//	// Returns: "SELECT orders.*, users.name FROM orders RIGHT JOIN users ON orders.user_id = users.id;"
func (b *Builder) RightJoin(table string, onCondition string) BuilderInterface {
	return b.Join(JOIN_RIGHT, table, onCondition)
}

// InnerJoin adds an INNER JOIN clause to the query.
// The table parameter specifies the table to join.
// The onCondition parameter specifies the join condition.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).Table("orders").
//	  InnerJoin("users", "orders.user_id = users.id").
//	  Select([]string{"orders.*", "users.name"})
//	// Returns: "SELECT orders.*, users.name FROM orders INNER JOIN users ON orders.user_id = users.id;"
func (b *Builder) InnerJoin(table string, onCondition string) BuilderInterface {
	return b.Join(JOIN_INNER, table, onCondition)
}

/**
 * The delete method deletes a row in a table. For deleting a database
 * or table use the drop method.
 * <code>
 * // Deleting a row
 * sql := builder.Table("STATES").Where("STATE_NAME","=","Alabama").Delete();
 * </code>
 * @return string
 * @access public
 */
// Delete deletes rows from a table
func (b *Builder) Delete() (string, []interface{}, error) {
	// Reset parameters for new query
	b.resetParams()

	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", nil, err
	}

	if b.sqlTableName == "" {
		return "", nil, ErrMissingTable
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		var err error
		where, err = b.whereToSql(b.sqlWhere)
		if err != nil {
			return "", nil, err
		}
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	sql := ""
	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
		sql = "DELETE FROM " + b.quoteTable(b.sqlTableName) + where + orderBy + limit + offset + ";"
	}
	return sql, b.params, nil
}

// Drop deletes a table or a view
func (b *Builder) Drop() (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	if !isTable && !isView {
		return "", ErrMissingTable
	}

	sql := ""

	if isTable {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP TABLE " + b.quoteTable(b.sqlTableName) + ";"
		}
	}

	if isView {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP VIEW " + b.quoteTable(b.sqlViewName) + ";"
		}
	}

	return sql, nil
}

func (b *Builder) DropIfExists() (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	if !isTable && !isView {
		return "", ErrMissingTable
	}

	sql := ""

	if isTable {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP TABLE IF EXISTS " + b.quoteTable(b.sqlTableName) + ";"
		}
	}

	if isView {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP VIEW IF EXISTS " + b.quoteTable(b.sqlViewName) + ";"
		}
	}

	return sql, nil
}

func (b *Builder) Limit(limit int64) BuilderInterface {
	b.sqlLimit = limit
	return b
}

func (b *Builder) Offset(offset int64) BuilderInterface {
	b.sqlOffset = offset
	return b
}

func (b *Builder) GroupBy(groupBy GroupBy) BuilderInterface {
	b.sqlGroupBy = append(b.sqlGroupBy, groupBy)
	return b
}

func (b *Builder) OrderBy(columnName, direction string) BuilderInterface {
	if strings.EqualFold(direction, "desc") || strings.EqualFold(direction, "descending") {
		direction = "DESC"
	} else {
		direction = "ASC"
	}

	b.sqlOrderBy = append(b.sqlOrderBy, OrderBy{
		Column:    columnName,
		Direction: direction,
	})

	return b
}

// Rename renames a table or a view
func (b *Builder) TableRename(oldTableName, newTableName string) (sql string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sql = "EXEC sp_rename " + b.quoteTable(oldTableName) + ", " + b.quoteTable(newTableName) + ", 'OBJECT';"
		return sql, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sql = "ALTER TABLE " + b.quoteTable(oldTableName) + " RENAME TO " + b.quoteTable(newTableName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sql = "ALTER TABLE " + b.quoteTable(oldTableName) + " RENAME " + b.quoteTable(newTableName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sql = "ALTER TABLE " + b.quoteTable(oldTableName) + " RENAME TO " + b.quoteTable(newTableName) + ";"
		return sql, nil
	}

	return "", errors.New("renaming a table is not supported for driver " + b.Dialect + "")
}

// TableColumnAdd adds a column to the table
func (b *Builder) TableColumnAdd(tableName string, column Column) (sql string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	return "", errors.New("adding a column is not supported for driver " + b.Dialect + "")
}

// TableColumnChange changes a column in the table
func (b *Builder) TableColumnChange(tableName string, column Column) (sqlString string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " ALTER COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " ALTER COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " MODIFY COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " ALTER COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	return "", errors.New("modifying a column is not supported for driver " + b.Dialect + "")
}

// TableColumnDrop drops a column from the table
func (b *Builder) TableColumnDrop(tableName, columnName string) (sqlString string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	return "", errors.New("dropping a column is not supported for driver " + b.Dialect + "")
}

// TableColumnExists checks if a column exists in a table for various database types
//
//	Example:
//	b := NewBuilder(DIALECT_MYSQL)
//	sqlString, sqlParams, err := b.TableColumnExists("test_table", "test_column")
//
// Params:
// - tableName: The name of the table to check.
// - columnName: The name of the column to check.
//
// Returns:
// - sql: The SQL statement to check for the existence of the column.
// - params: An array of parameters to be bound to the statement.
// - err: An error object, if any.
func (b *Builder) TableColumnExists(tableName, columnName string) (sql string, params []interface{}, err error) {
	switch b.Dialect {
	case DIALECT_MYSQL:
		return "SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = ?", []interface{}{tableName, columnName}, nil
	case DIALECT_POSTGRES:
		return "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", []interface{}{tableName, columnName}, nil
	case DIALECT_SQLITE:
		return "SELECT 1 FROM pragma_table_info(?) WHERE name = ?", []interface{}{tableName, columnName}, nil
	default:
		return "", nil, fmt.Errorf("database type '%s' not supported", b.Dialect)
	}
}

func (b *Builder) TableColumnRename(tableName, oldColumnName, newColumnName string) (sql string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sql = "EXEC sp_rename " + b.quoteTable(tableName) + "." + b.quoteTable(oldColumnName) + ", " + b.quoteTable(newColumnName) + ", 'COLUMN';"
		return sql, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " RENAME COLUMN " + b.quoteTable(oldColumnName) + " TO " + b.quoteTable(newColumnName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " RENAME COLUMN " + b.quoteTable(oldColumnName) + " TO " + b.quoteTable(newColumnName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " RENAME COLUMN " + b.quoteTable(oldColumnName) + " TO " + b.quoteTable(newColumnName) + ";"
		return sql, nil
	}

	return "", errors.New("renaming a column is not supported for driver " + b.Dialect + "")
}

/** The <b>select</b> method selects rows from a table, based on criteria.
 * <code>
 * // Selects all the rows from the table
 * $db->table("USERS")->select();
 *
 * // Selects the rows where the column NAME is different from Peter, in descending order
 * $db->table("USERS")
 *     ->where("NAME","!=","Peter")
 *     ->orderby("NAME","desc")
 *     ->select();
 * </code>
 * @return mixed rows as associative array, false on error
 * @access public
 */
func (b *Builder) Select(columns []string) (string, []interface{}, error) {
	// Reset parameters for new query
	b.resetParams()

	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", nil, err
	}

	if b.sqlTableName == "" {
		return "", nil, ErrMissingTable
	}

	// Store the select columns for subquery use
	b.sqlSelectColumns = columns

	join := b.joinToSQL()

	groupBy := ""
	if len(b.sqlGroupBy) > 0 {
		groupBy = b.groupByToSql(b.sqlGroupBy)
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		var err error
		where, err = b.whereToSql(b.sqlWhere)
		if err != nil {
			return "", nil, err
		}
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	columnsStr := "*"

	if len(columns) > 0 {
		for index, column := range columns {
			if strings.Contains(column, "(") {
				columns[index] = column // Do not quote function calls
			} else {
				columns[index] = b.quoteColumn(column)
			}
		}
		columnsStr = strings.Join(columns, ", ")
	}

	sql := ""

	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
		sql = "SELECT " + columnsStr + " FROM " + b.quoteTable(b.sqlTableName) + join + where + groupBy + orderBy + limit + offset + ";"
	} else if b.Dialect == DIALECT_MSSQL {
		sql = "SELECT " + columnsStr + " FROM " + b.quoteTable(b.sqlTableName) + join + where + groupBy + orderBy + ";"
	}

	return sql, b.params, nil
}

/**
 * The <b>update</b> method updates the values of a row in a table.
 * <code>
 * $updated_user = array("USER_MANE"=>"Mike");
 * $database->table("USERS")->where("USER_NAME","==","Peter")->update($updated_user);
 * </code>
 * @param Array an associative array, where keys are the column names of the table
 * @return int 0 or 1, on success, false, otherwise
 * @access public
 */
func (b *Builder) Insert(columnValuesMap map[string]string) (string, []interface{}, error) {
	// Reset parameters for new query
	b.resetParams()

	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", nil, err
	}

	if b.sqlTableName == "" {
		return "", nil, ErrMissingTable
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	columnNames := []string{}
	columnValues := []string{}

	// Order keys
	keys := make([]string, 0, len(columnValuesMap))
	for k := range columnValuesMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, columnName := range keys {
		columnValue := columnValuesMap[columnName]
		columnNames = append(columnNames, b.quoteColumn(columnName))

		// Use parameterized queries by default, unless interpolatedValues mode is enabled
		if b.interpolatedValues {
			columnValues = append(columnValues, b.quoteValue(columnValue))
		} else {
			columnValues = append(columnValues, b.addParam(columnValue))
		}
	}

	return "INSERT INTO " + b.quoteTable(b.sqlTableName) + " (" + strings.Join(columnNames, ", ") + ") VALUES (" + strings.Join(columnValues, ", ") + ")" + limit + offset + ";", b.params, nil
}

// Truncate removes all data from a table.
// The method generates database-specific SQL:
//   - MySQL: TRUNCATE TABLE `table_name`;
//   - PostgreSQL: TRUNCATE TABLE "table_name";
//   - SQLite: DELETE FROM "table_name";
//   - MSSQL: TRUNCATE TABLE [table_name];
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Truncate()
//	// Returns: "TRUNCATE TABLE `users`;"
func (b *Builder) Truncate() (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}

	switch b.Dialect {
	case DIALECT_MYSQL:
		return "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName) + ";", nil
	case DIALECT_POSTGRES:
		return "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName) + ";", nil
	case DIALECT_SQLITE:
		return "DELETE FROM " + b.quoteTable(b.sqlTableName) + ";", nil
	case DIALECT_MSSQL:
		return "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName) + ";", nil
	default:
		return "", ErrInvalidDialect
	}
}

// TruncateWithOptions removes all data from a table with additional options.
// The opts parameter provides database-specific behavior:
//   - Cascade (PostgreSQL): adds CASCADE to handle foreign key constraints
//   - ResetIdentity (MSSQL): resets identity column seed value after truncation
//
// Examples:
//
//	// PostgreSQL with CASCADE for foreign key constraints
//	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("orders").
//	  TruncateWithOptions(sb.TruncateOptions{Cascade: true})
//	// Returns: "TRUNCATE TABLE \"orders\" CASCADE;"
//
//	// MSSQL with identity reset
//	sql := sb.NewBuilder(sb.DIALECT_MSSQL).Table("users").
//	  TruncateWithOptions(sb.TruncateOptions{ResetIdentity: true})
//	// Returns: "TRUNCATE TABLE [users]; DBCC CHECKIDENT ('users', RESEED, 0)"
func (b *Builder) TruncateWithOptions(opts TruncateOptions) (string, error) {
	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", err
	}

	if b.sqlTableName == "" {
		return "", ErrMissingTable
	}

	switch b.Dialect {
	case DIALECT_MYSQL:
		return "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName) + ";", nil

	case DIALECT_POSTGRES:
		sql := "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName)
		if opts.Cascade {
			sql += " CASCADE"
		}
		return sql, nil

	case DIALECT_SQLITE:
		return "DELETE FROM " + b.quoteTable(b.sqlTableName) + ";", nil

	case DIALECT_MSSQL:
		if opts.ResetIdentity {
			return "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName) + "; DBCC CHECKIDENT ('" + b.sqlTableName + "', RESEED, 0)", nil
		}
		return "TRUNCATE TABLE " + b.quoteTable(b.sqlTableName) + ";", nil

	default:
		return "", ErrInvalidDialect
	}
}

// Update updates the values of rows in a table.
//
// Example:
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").
//	  Where(sb.Where{Column: "id", Operator: "=", Value: "1"}).
//	  Update(map[string]string{"name": "John", "email": "john@example.com"})
func (b *Builder) Update(columnValues map[string]string) (string, []interface{}, error) {
	// Reset parameters for new query
	b.resetParams()

	// First validate any collected errors from fluent chaining
	if err := b.validateAndReturnError(); err != nil {
		return "", nil, err
	}

	if b.sqlTableName == "" {
		return "", nil, ErrMissingTable
	}

	join := b.joinToSQL()

	groupBy := ""
	if len(b.sqlGroupBy) > 0 {
		groupBy = b.groupByToSql(b.sqlGroupBy)
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		var err error
		where, err = b.whereToSql(b.sqlWhere)
		if err != nil {
			return "", nil, err
		}
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	// Order keys
	keys := make([]string, 0, len(columnValues))
	for k := range columnValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	updateSql := []string{}
	for _, columnName := range keys {
		columnValue := columnValues[columnName]

		// Use parameterized queries by default, unless interpolatedValues mode is enabled
		if b.interpolatedValues {
			updateSql = append(updateSql, b.quoteColumn(columnName)+"="+b.quoteValue(columnValue))
		} else {
			updateSql = append(updateSql, b.quoteColumn(columnName)+"="+b.addParam(columnValue))
		}
	}

	return "UPDATE " + b.quoteTable(b.sqlTableName) + " SET " + strings.Join(updateSql, ", ") + join + where + groupBy + orderBy + limit + offset + ";", b.params, nil
}

// Where adds a WHERE clause to the query.
//
// Security Note: Values are properly quoted to prevent SQL injection.
// Avoid using the Raw field with user input. Use Column/Operator/Value instead.
//
// Example:
//
//	// Safe - values are properly quoted
//	sql := builder.Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail})
//
//	// Dangerous - avoid using Raw with user input
//	sql := builder.Where(&sb.Where{Raw: "email = '" + userEmail + "'"})
func (b *Builder) Where(where *Where) BuilderInterface {
	if where == nil {
		return b
	}

	b.sqlWhere = append(b.sqlWhere, *where)
	return b
}

// Subquery creates a subquery builder for use in WHERE clauses.
// Returns a new Builder instance with the same dialect for building subqueries.
//
// Example:
//
//	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("orders").
//	  Select([]string{"user_id"}).
//	  Where(&sb.Where{Column: "total", Operator: ">", Value: 1000})
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("users").
//	  InSubquery(subquery).
//	  Select([]string{"name"})
func (b *Builder) Subquery() BuilderInterface {
	return NewBuilder(b.Dialect)
}

// assertBuilder safely converts a BuilderInterface to *Builder with proper error handling.
// This prevents unsafe type assertions that could panic at runtime.
func (b *Builder) assertBuilder(subquery BuilderInterface) (*Builder, error) {
	if subquery == nil {
		return nil, ErrNilSubquery
	}

	builder, ok := subquery.(*Builder)
	if !ok {
		return nil, NewValidationError("subquery must be a *Builder instance")
	}

	return builder, nil
}

// Exists adds an EXISTS subquery condition to the query.
// The subquery parameter is the subquery to check for existence.
//
// Example:
//
//	activeOrders := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("orders").
//	  Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("users").
//	  Exists(activeOrders).
//	  Select([]string{"name"})
func (b *Builder) Exists(subquery BuilderInterface) (BuilderInterface, error) {
	if subquery == nil {
		return nil, ErrNilSubquery
	}

	where := Where{
		Operator: "EXISTS",
		Subquery: func() *Builder {
			builder, err := b.assertBuilder(subquery)
			if err != nil {
				// This should not happen as we already validated subquery above
				return nil
			}
			return builder
		}(),
		IsNot: false,
	}
	b.sqlWhere = append(b.sqlWhere, where)
	return b, nil
}

// NotExists adds a NOT EXISTS subquery condition to the query.
// The subquery parameter is the subquery to check for non-existence.
//
// Example:
//
//	activeOrders := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("orders").
//	  Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("users").
//	  NotExists(activeOrders).
//	  Select([]string{"name"})
func (b *Builder) NotExists(subquery BuilderInterface) (BuilderInterface, error) {
	if subquery == nil {
		return nil, ErrNilSubquery
	}

	where := Where{
		Operator: "EXISTS",
		Subquery: func() *Builder {
			builder, err := b.assertBuilder(subquery)
			if err != nil {
				// This should not happen as we already validated subquery above
				return nil
			}
			return builder
		}(),
		IsNot: true,
	}
	b.sqlWhere = append(b.sqlWhere, where)
	return b, nil
}

// InSubquery adds an IN subquery condition to the query.
// The subquery parameter is the subquery to use for the IN clause.
//
// Example:
//
//	highValueUsers := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("orders").
//	  Select([]string{"user_id"}).
//	  Where(&sb.Where{Column: "total", Operator: ">", Value: 1000})
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("users").
//	  InSubquery(highValueUsers).
//	  Select([]string{"name"})
func (b *Builder) InSubquery(subquery BuilderInterface) (BuilderInterface, error) {
	if subquery == nil {
		return nil, ErrNilSubquery
	}

	where := Where{
		Operator: "IN",
		Subquery: func() *Builder {
			builder, err := b.assertBuilder(subquery)
			if err != nil {
				// This should not happen as we already validated subquery above
				return nil
			}
			return builder
		}(),
		IsNot: false,
	}
	b.sqlWhere = append(b.sqlWhere, where)
	return b, nil
}

// NotInSubquery adds a NOT IN subquery condition to the query.
// The subquery parameter is the subquery to use for the NOT IN clause.
//
// Example:
//
//	inactiveUsers := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("orders").
//	  Select([]string{"user_id"}).
//	  Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"})
//
//	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
//	  Table("users").
//	  NotInSubquery(inactiveUsers).
//	  Select([]string{"name"})
func (b *Builder) NotInSubquery(subquery BuilderInterface) (BuilderInterface, error) {
	if subquery == nil {
		return nil, ErrNilSubquery
	}

	where := Where{
		Operator: "IN",
		Subquery: func() *Builder {
			builder, err := b.assertBuilder(subquery)
			if err != nil {
				// This should not happen as we already validated subquery above
				return nil
			}
			return builder
		}(),
		IsNot: true,
	}
	b.sqlWhere = append(b.sqlWhere, where)
	return b, nil
}

// columnsToSQL converts the columns statements to SQL.
func (b *Builder) columnsToSQL(columns []Column) string {
	columnSQLs := []string{}

	for i := 0; i < len(columns); i++ {
		column := columns[i]
		columnSQLs = append(columnSQLs, b.columnSQLGenerator.GenerateSQL(column))
	}

	return strings.Join(columnSQLs, ", ")
}

// joinToSQL converts the joins slice to SQL string
func (b *Builder) joinToSQL() string {
	if len(b.sqlJoins) == 0 {
		return ""
	}

	var joins []string
	for _, join := range b.sqlJoins {
		tableExpr := b.quoteTable(join.Table)
		if join.Alias != "" {
			tableExpr += " AS " + b.quoteTable(join.Alias)
		}

		joinSQL := fmt.Sprintf(" %s JOIN %s ON %s",
			string(join.Type),
			tableExpr,
			join.OnCondition)
		joins = append(joins, joinSQL)
	}

	return " " + strings.Join(joins, " ")
}

func (b *Builder) groupByToSql(groupBys []GroupBy) string {
	sql := []string{}
	for _, groupBy := range groupBys {
		sql = append(sql, b.quoteColumn(groupBy.Column))
	}

	if len(sql) > 0 {
		return " GROUP BY " + strings.Join(sql, ",")
	}

	return ""
}

// /**
//      * Joins tables to SQL.
//      * @return String the join SQL string
//      * @access private
//      */
// 	 private function join_to_sql($join, $table_name)
// 	 {
// 		 $sql = '';
// 		 // MySQL
// 		 if ($this->database_type == 'mysql') {
// 			 foreach ($join as $what) {
// 				 $type = $what[3] ?? '';
// 				 $alias = $what[4] ?? '';
// 				 $sql .= ' ' . $type . ' JOIN `' . $what[0] . '`';
// 				 if ($alias != "") {
// 					 $sql .= ' AS ' . $alias . '';
// 					 $what[0] = $alias;
// 				 }
// 				 if ($what[1] == $what[2]) {
// 					 $sql .= ' USING (`' . $what[1] . '`)';
// 				 } else {
// 					 $sql .= ' ON ' . $table_name . '.' . $what[1] . '=' . $what[0] . '.' . $what[2];
// 				 }
// 			 }
// 		 }
// 		 // SQLite
// 		 if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
// 			 foreach ($join as $what) {
// 				 $type = $what[3] ?? '';
// 				 $alias = $what[4] ?? '';
// 				 $sql .= " $type JOIN '" . $what[0] . "'";
// 				 if ($alias != "") {
// 					 $sql .= " AS '$alias'";
// 					 $what[0] = $alias;
// 				 }
// 				 $sql .= ' ON ' . $table_name . '.' . $what[1] . '=' . $what[0] . '.' . $what[2];
// 			 }
// 		 }

// 		 return $sql;
// 	 }

func (b *Builder) orderByToSql(orderBys []OrderBy) string {
	sql := []string{}

	if b.Dialect == DIALECT_MYSQL {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if b.Dialect == DIALECT_POSTGRES {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if b.Dialect == DIALECT_SQLITE {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if len(sql) > 0 {
		return ` ORDER BY ` + strings.Join(sql, `,`)
	}

	return ""
}

func (b *Builder) quoteColumn(columnName string) string {
	columnSplit := strings.Split(columnName, ".")
	columnQuoted := []string{}

	for _, columnPart := range columnSplit {
		if columnPart == "*" {
			columnQuoted = append(columnQuoted, columnPart)
			continue
		}

		if strings.Contains(columnPart, "(") {
			columnQuoted = append(columnQuoted, columnPart)
		}

		columnQuoted = append(columnQuoted, b.quote(columnPart, "column"))
	}

	return strings.Join(columnQuoted, ".")
}

func (b *Builder) quoteTable(tableName string) string {
	tableSplit := strings.Split(tableName, ".")
	tableQuoted := []string{}

	for _, tablePart := range tableSplit {
		tableQuoted = append(tableQuoted, b.quote(tablePart, "table"))
	}

	return strings.Join(tableQuoted, ".")
}

/**
 * The <b>tables</b> method returns the names of all the tables, that
 * exist in the database.
 * <code>
 * foreach($database->tables() as $table){
 *     echo $table;
 * }
 * </code>
 * @param String the name of the table
 * @return array the names of the tables
 * @access public
 */
//  func (b *Builder) Tables(value string)
//  {
// 	 $tables = array();

// 	 if ($this->database_type == 'mysql') {
// 		 //$sql = "SHOW TABLES";
// 		 $sql = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_SCHEMA='" . $this->database_name . "'";
// 		 $result = $this->executeQuery($sql);
// 		 if ($result === false)
// 			 return false;
// 		 foreach ($result as $row) {
// 			 $tables[] = $row['TABLE_NAME'];
// 		 }
// 		 return $tables;
// 	 }

// 	 if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
// 		 $sql = "SELECT * FROM 'SQLITE_MASTER' WHERE type='table' ORDER BY NAME ASC";
// 		 $result = $this->executeQuery($sql);
// 		 if ($result === false) {
// 			 return false;
// 		 }
// 		 foreach ($result as $row) {
// 			 $tables[] = $row['name'];
// 		 }
// 		 return $tables;
// 	 }
// 	 return false;
//  }
