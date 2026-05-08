---
path: api_reference.md
page-type: reference
summary: Complete API reference documentation for all SB SQL Builder methods and types.
tags: [api, reference, methods, types]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# API Reference

Complete API reference for the SB SQL Builder library.

## Core Types

### Builder

Main entry point for constructing SQL queries.

```go
type Builder struct {
    dialect           string
    sqlTable          string
    sqlColumns        []Column
    sqlWhere          []Where
    sqlJoins          []Join
    sqlOrderBy        []OrderBy
    sqlLimit          int
    sqlOffset         int
    sqlParameters     []interface{}
    parameterIndex    int
    sqlErrors         []error
    sqlSelectColumns  []string
}
```

### Column

Represents a table column with its properties.

```go
type Column struct {
    Name         string
    Type         string
    Length       int
    PrimaryKey   bool
    AutoIncrement bool
    Nullable     bool
    Unique       bool
    DefaultValue string
    Comment      string
}
```

### Where

Represents a WHERE condition with support for subqueries.

```go
type Where struct {
    Raw       string
    Column    string
    Operator  string
    Type      string
    Value     string
    Subquery  *Builder
    IsNot     bool
    Children  []Where
}
```

### Join

Represents a JOIN operation.

```go
type Join struct {
    Type        JoinType
    Table       string
    Alias       string
    OnCondition string
    Columns     []string
}

type JoinType string

const (
    JOIN_INNER JoinType = "INNER"
    JOIN_LEFT  JoinType = "LEFT"
    JOIN_RIGHT JoinType = "RIGHT"
    JOIN_FULL  JoinType = "FULL"
    JOIN_CROSS JoinType = "CROSS"
)
```

### IndexColumn

Represents a column in an index with direction and prefix length.

```go
type IndexColumn struct {
    Name        string
    Direction   string // ASC or DESC
    PrefixLength int   // For MySQL text columns
}
```

### IndexOptions

Options for advanced index creation.

```go
type IndexOptions struct {
    Unique      bool
    IfNotExists bool
    Using       string
    Columns     []IndexColumn
    Include     []string
    Where       string
    Storage     string
    Comment     string
}
```

## Constructor

### NewBuilder

Creates a new Builder instance with the specified dialect.

```go
func NewBuilder(dialect string) *Builder
```

**Parameters:**
- `dialect`: Database dialect (DIALECT_MYSQL, DIALECT_POSTGRES, DIALECT_SQLITE, DIALECT_MSSQL)

**Returns:**
- `*Builder`: New Builder instance

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL)
```

## Table Methods

### Table

Sets the table name for the query.

```go
func (b *Builder) Table(name string) BuilderInterface
```

**Parameters:**
- `name`: Table name

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users")
```

## Column Methods

### Column

Adds a column definition for CREATE TABLE operations.

```go
func (b *Builder) Column(column Column) BuilderInterface
```

**Parameters:**
- `column`: Column definition

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Column(sb.Column{
        Name:       "id",
        Type:       sb.COLUMN_TYPE_STRING,
        Length:     40,
        PrimaryKey: true,
    })
```

## WHERE Methods

### Where

Adds a WHERE condition to the query.

```go
func (b *Builder) Where(where *Where) BuilderInterface
```

**Parameters:**
- `where`: WHERE condition

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    })
```

### WhereRaw

Adds a raw WHERE clause string.

```go
func (b *Builder) WhereRaw(raw string) BuilderInterface
```

**Parameters:**
- `raw`: Raw WHERE clause

**Returns:**
- `BuilderInterface`: Builder instance for chaining

## JOIN Methods

### Join

Adds a JOIN clause to the query.

```go
func (b *Builder) Join(joinType JoinType, table string, onCondition string) BuilderInterface
```

**Parameters:**
- `joinType`: Type of join (JOIN_INNER, JOIN_LEFT, JOIN_RIGHT, JOIN_FULL, JOIN_CROSS)
- `table`: Table to join
- `onCondition`: Join condition

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "orders.user_id = users.id")
```

### JoinWithAlias

Adds a JOIN clause with table alias.

```go
func (b *Builder) JoinWithAlias(joinType JoinType, table string, alias string, onCondition string) BuilderInterface
```

**Parameters:**
- `joinType`: Type of join
- `table`: Table to join
- `alias`: Table alias
- `onCondition`: Join condition

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id")
```

### InnerJoin

Convenience method for INNER JOIN.

```go
func (b *Builder) InnerJoin(table string, onCondition string) BuilderInterface
```

### LeftJoin

Convenience method for LEFT JOIN.

```go
func (b *Builder) LeftJoin(table string, onCondition string) BuilderInterface
```

### RightJoin

Convenience method for RIGHT JOIN.

```go
func (b *Builder) RightJoin(table string, onCondition string) BuilderInterface
```

## Subquery Methods

### Exists

Adds an EXISTS subquery condition.

```go
func (b *Builder) Exists(subquery BuilderInterface) (BuilderInterface, error)
```

**Parameters:**
- `subquery`: Subquery builder

**Returns:**
- `BuilderInterface`: Builder instance for chaining
- `error`: Error if subquery is nil

**Example:**
```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(subquery)
```

### NotExists

Adds a NOT EXISTS subquery condition.

```go
func (b *Builder) NotExists(subquery BuilderInterface) (BuilderInterface, error)
```

### InSubquery

Adds an IN subquery condition.

```go
func (b *Builder) InSubquery(subquery BuilderInterface) (BuilderInterface, error)
```

### NotInSubquery

Adds a NOT IN subquery condition.

```go
func (b *Builder) NotInSubquery(subquery BuilderInterface) (BuilderInterface, error)
```

### Subquery

Creates a subquery builder.

```go
func (b *Builder) Subquery() BuilderInterface
```

## ORDER BY Methods

### OrderBy

Adds an ORDER BY clause.

```go
func (b *Builder) OrderBy(column string, direction string) BuilderInterface
```

**Parameters:**
- `column`: Column name
- `direction`: Sort direction (ASC or DESC)

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    OrderBy("created_at", "DESC")
```

## LIMIT/OFFSET Methods

### Limit

Adds a LIMIT clause.

```go
func (b *Builder) Limit(limit int) BuilderInterface
```

**Parameters:**
- `limit`: Maximum number of rows

**Returns:**
- `BuilderInterface`: Builder instance for chaining

### Offset

Adds an OFFSET clause.

```go
func (b *Builder) Offset(offset int) BuilderInterface
```

**Parameters:**
- `offset`: Number of rows to skip

**Returns:**
- `BuilderInterface`: Builder instance for chaining

## SQL Generation Methods

### Select

Generates a SELECT query.

```go
func (b *Builder) Select(columns []string) (string, []interface{}, error)
```

**Parameters:**
- `columns`: Columns to select

**Returns:**
- `string`: Generated SQL
- `[]interface{`: Parameters for parameterized query
- `error`: Error if validation fails

**Example:**
```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"id", "name", "email"})
```

### Insert

Generates an INSERT query.

```go
func (b *Builder) Insert(data map[string]interface{}) (string, []interface{}, error)
```

**Parameters:**
- `data`: Map of column names to values

**Returns:**
- `string`: Generated SQL
- `[]interface{`: Parameters
- `error`: Error if validation fails

### Update

Generates an UPDATE query.

```go
func (b *Builder) Update(data map[string]interface{}) (string, []interface{}, error)
```

**Parameters:**
- `data`: Map of column names to values

**Returns:**
- `string`: Generated SQL
- `[]interface{`: Parameters
- `error`: Error if validation fails

### Delete

Generates a DELETE query.

```go
func (b *Builder) Delete() (string, []interface{}, error)
```

**Returns:**
- `string`: Generated SQL
- `[]interface{`: Parameters
- `error`: Error if validation fails

## Schema Methods

### Create

Generates a CREATE TABLE query.

```go
func (b *Builder) Create() (string, error)
```

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### CreateIfNotExists

Generates a CREATE TABLE IF NOT EXISTS query.

```go
func (b *Builder) CreateIfNotExists() (string, error)
```

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### Drop

Generates a DROP TABLE query.

```go
func (b *Builder) Drop() (string, error)
```

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### DropIfExists

Generates a DROP TABLE IF EXISTS query.

```go
func (b *Builder) DropIfExists() (string, error)
```

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### Truncate

Generates a TRUNCATE TABLE query.

```go
func (b *Builder) Truncate() (string, error)
```

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### TruncateWithOptions

Generates a TRUNCATE TABLE query with options.

```go
func (b *Builder) TruncateWithOptions(opts TruncateOptions) (string, error)
```

**Parameters:**
- `opts`: Truncate options

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

## Index Methods

### CreateIndex

Generates a CREATE INDEX query.

```go
func (b *Builder) CreateIndex(name string, columns ...string) (string, error)
```

**Parameters:**
- `name`: Index name
- `columns`: Column names

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### CreateIndexWithOptions

Generates a CREATE INDEX query with advanced options.

```go
func (b *Builder) CreateIndexWithOptions(name string, opts IndexOptions) (string, error)
```

**Parameters:**
- `name`: Index name
- `opts`: Index options

**Returns:**
- `string`: Generated SQL
- `error`: Error if validation fails

### CreateUniqueIndex

Generates a CREATE UNIQUE INDEX query.

```go
func (b *Builder) CreateUniqueIndex(name string, columns ...string) (string, error)
```

### CreateCompositeIndex

Generates a CREATE INDEX query with explicit column ordering.

```go
func (b *Builder) CreateCompositeIndex(name string, columns []IndexColumn) (string, error)
```

### CreatePartialIndex

Generates a CREATE INDEX query with WHERE predicate.

```go
func (b *Builder) CreatePartialIndex(name string, where string, columns ...string) (string, error)
```

### CreateCoveringIndex

Generates a CREATE INDEX query with INCLUDE clause.

```go
func (b *Builder) CreateCoveringIndex(name string, include []string, columns ...string) (string, error)
```

### DropIndex

Generates a DROP INDEX query.

```go
func (b *Builder) DropIndex(indexName string) (string, error)
```

### DropIndexIfExists

Generates a DROP INDEX IF EXISTS query.

```go
func (b *Builder) DropIndexIfExists(indexName string) (string, error)
```

### DropIndexWithOptions

Generates a DROP INDEX query with options.

```go
func (b *Builder) DropIndexWithOptions(name string, opts DropIndexOptions) (string, error)
```

## Utility Methods

### WithInterpolatedValues

Disables parameterized queries and uses interpolated values instead (legacy mode).

```go
func (b *Builder) WithInterpolatedValues() BuilderInterface
```

**Returns:**
- `BuilderInterface`: Builder instance for chaining

**Example:**
```go
sql, _, err := builder.
    WithInterpolatedValues().
    Select([]string{"*"})
// Legacy mode: values interpolated directly into SQL
```

## Constants

### Dialect Constants

```go
const (
    DIALECT_MYSQL    = "mysql"
    DIALECT_POSTGRES = "postgres"
    DIALECT_SQLITE   = "sqlite"
    DIALECT_MSSQL    = "mssql"
)
```

### Column Type Constants

```go
const (
    COLUMN_TYPE_STRING     = "string"
    COLUMN_TYPE_INT        = "int"
    COLUMN_TYPE_BIGINT     = "bigint"
    COLUMN_TYPE_FLOAT      = "float"
    COLUMN_TYPE_DECIMAL    = "decimal"
    COLUMN_TYPE_BOOLEAN    = "boolean"
    COLUMN_TYPE_DATE       = "date"
    COLUMN_TYPE_DATETIME   = "datetime"
    COLUMN_TYPE_TIMESTAMP  = "timestamp"
    COLUMN_TYPE_TEXT       = "text"
    COLUMN_TYPE_JSON       = "json"
)
```

### Index Type Constants

```go
const (
    INDEX_TYPE_BTREE    = "BTREE"
    INDEX_TYPE_HASH     = "HASH"
    INDEX_TYPE_GIN      = "GIN"
    INDEX_TYPE_GIST     = "GIST"
    INDEX_TYPE_BRIN     = "BRIN"
    INDEX_TYPE_FULLTEXT = "FULLTEXT"
    INDEX_TYPE_SPATIAL  = "SPATIAL"
)
```

## Error Types

### BuilderError

Structured error type for builder operations.

```go
type BuilderError struct {
    Type    string
    Message string
}

func (e *BuilderError) Error() string
```

### Standard Errors

```go
var (
    ErrEmptyTableName     = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
    ErrEmptyColumnName    = &BuilderError{Type: "ValidationError", Message: "column name cannot be empty"}
    ErrNilSubquery        = &BuilderError{Type: "ArgumentError", Message: "subquery cannot be nil"}
    ErrEmptyOnCondition   = &BuilderError{Type: "ValidationError", Message: "ON condition cannot be empty"}
    ErrEmptyIndexName     = &BuilderError{Type: "ValidationError", Message: "index name cannot be empty"}
    ErrMissingTable       = &BuilderError{Type: "ValidationError", Message: "table name is required"}
    ErrEmptyColumns       = &BuilderError{Type: "ValidationError", Message: "columns cannot be empty"}
)
```

## See Also

- [Overview](overview.md)
- [Features](features/)
- [Architecture](architecture.md)
