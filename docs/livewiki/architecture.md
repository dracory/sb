---
path: architecture.md
page-type: reference
summary: System architecture and design patterns used in the SB SQL Builder library.
tags: [architecture, design-patterns, system-design]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# Architecture

SB uses a combination of design patterns to provide a clean, extensible SQL building API.

## Core Patterns

### Builder Pattern

The Builder pattern is the foundation of SB, allowing complex SQL queries to be constructed through method chaining:

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    OrderBy("created_at", "DESC").
    Limit(10)
```

**Benefits:**
- Readable, expressive code
- Step-by-step query construction
- Immutability where appropriate
- Type-safe operations

### Strategy Pattern

Different SQL generation strategies for each database dialect:

```go
type ColumnSQLGenerator interface {
    GenerateColumnSQL(column Column) string
}

type MySQLColumnSQLGenerator struct{}
type PostgreSQLColumnSQLGenerator struct{}
type SQLiteColumnSQLGenerator struct{}
type MSSQLColumnSQLGenerator struct{}
```

**Benefits:**
- Database-specific optimizations
- Easy to add new dialects
- Consistent API across dialects
- Isolated dialect logic

### Facade Pattern

Simplified interface over Go's `database/sql` package:

```go
type Database struct {
    db *sql.DB
    dialect string
}
```

**Benefits:**
- Simpler transaction management
- Consistent error handling
- Connection pooling abstraction
- Easier testing

### Template Method Pattern

Common SQL patterns with database-specific variations:

```go
func (b *Builder) Select(columns []string) (string, []interface{}, error) {
    // Common logic
    // Dialect-specific quoting
    // Parameter handling
}
```

**Benefits:**
- Code reuse across dialects
- Consistent query structure
- Easy to maintain
- Dialect-specific optimizations

## Component Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         Builder                             │
├─────────────────────────────────────────────────────────────┤
│  State:                                                     │
│  - sqlTable string                                          │
│  - sqlColumns []Column                                      │
│  - sqlWhere []Where                                         │
│  - sqlJoins []Join                                          │
│  - sqlParameters []interface{}                              │
│  - sqlErrors []error                                        │
│  - dialect string                                           │
├─────────────────────────────────────────────────────────────┤
│  Core Methods:                                              │
│  - Table(name string) BuilderInterface                      │
│  - Column(column Column) BuilderInterface                   │
│  - Where(where *Where) BuilderInterface                     │
│  - Select(columns []string) (string, []interface{}, error)  │
├─────────────────────────────────────────────────────────────┤
│  Advanced Methods:                                          │
│  - Join(joinType, table, onCondition) BuilderInterface     │
│  - Exists(subquery) (BuilderInterface, error)               │
│  - InSubquery(subquery) (BuilderInterface, error)            │
├─────────────────────────────────────────────────────────────┤
│  Schema Methods:                                            │
│  - Create() (string, error)                                 │
│  - Drop() (string, error)                                   │
│  - Truncate() (string, error)                               │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

### Query Building Flow

```
1. NewBuilder(dialect)
   ↓
2. Table(name)
   ↓
3. Column(column) × N
   ↓
4. Where(where) × N
   ↓
5. Join(joinType, table, condition) × N
   ↓
6. Select(columns) / Insert(data) / etc.
   ↓
7. SQL Generation + Parameter Collection
   ↓
8. Return (sql, params, error)
```

### Error Handling Flow

```
1. Fluent Method Call
   ↓
2. Validation Check
   ↓
3. Error Collection (if invalid)
   ↓
4. Return BuilderInterface (for chaining)
   ↓
5. SQL Generation Method Call
   ↓
6. validateAndReturnError()
   ↓
7. Return error if collected
   ↓
8. Generate SQL (if valid)
```

### Parameter Handling Flow

```
1. Where Clause with Value
   ↓
2. Value Processing
   ↓
3. NULL Check → Skip if NULL
   ↓
4. Parameter Addition to sqlParameters
   ↓
5. Placeholder Generation (dialect-specific)
   ↓
6. SQL Generation with Placeholders
   ↓
7. Return (sql, params, error)
```

## Key Components

### Builder

Main entry point for constructing queries. Maintains query state and provides fluent API.

**Responsibilities:**
- Query state management
- Method chaining
- SQL generation
- Parameter tracking
- Error collection

### Database

Wrapper around `sql.DB` with enhanced functionality.

**Responsibilities:**
- Connection management
- Transaction handling
- Query execution
- Result processing

### Column

Represents table columns with their properties.

**Responsibilities:**
- Column metadata
- Type information
- Constraint definitions
- Default values

### Where

Handles query conditions and subqueries.

**Responsibilities:**
- Condition representation
- Operator handling
- Subquery support
- AND/OR logic

### Join

Represents JOIN operations with table aliases.

**Responsibilities:**
- Join type specification
- Table aliasing
- Join conditions
- Column selection

## Error Handling Architecture

SB uses a sophisticated error handling system:

### Error Collection Pattern

Fluent methods collect errors during chaining:

```go
func (b *Builder) Column(column Column) BuilderInterface {
    if column.Name == "" {
        b.sqlErrors = append(b.sqlErrors, ErrEmptyColumnName)
        return b
    }
    // ... validation logic
    return b
}
```

### Error Validation at Build Time

Errors are validated when SQL is generated:

```go
func (b *Builder) Create() (string, error) {
    if err := b.validateAndReturnError(); err != nil {
        return "", err
    }
    // ... SQL generation
}
```

### Structured Error Types

```go
type BuilderError struct {
    Type    string
    Message string
}

var (
    ErrEmptyTableName = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
    ErrNilSubquery   = &BuilderError{Type: "ArgumentError", Message: "subquery cannot be nil"}
)
```

## Parameterization Architecture

### Parameter Tracking

The Builder struct tracks parameters:

```go
type Builder struct {
    // ... other fields
    sqlParameters []interface{}
    parameterIndex int
}
```

### Dialect-Specific Placeholders

Each dialect uses different placeholder syntax:

| Database | Placeholder | Example |
|----------|------------|---------|
| MySQL | `?` | `WHERE id = ?` |
| PostgreSQL | `$1, $2, $3` | `WHERE id = $1` |
| SQLite | `?` | `WHERE id = ?` |
| MSSQL | `@p1, @p2, @p3` | `WHERE id = @p1` |

### Parameter Ordering

Parameters are ordered consistently:
1. WHERE clause parameters (first)
2. INSERT/UPDATE values (second)
3. Subquery parameters (nested)

## Extensibility

### Adding New Dialects

1. Implement `ColumnSQLGenerator` interface
2. Add dialect constant
3. Update `NewBuilder` switch statement
4. Add quoting logic
5. Implement dialect-specific features

### Adding New SQL Operations

1. Add state to Builder struct
2. Add fluent method to Builder
3. Add method to BuilderInterface
4. Implement SQL generation logic
5. Add error handling
6. Write comprehensive tests

### Adding New Error Types

1. Define error variable in `errors.go`
2. Use in validation logic
3. Update error collection
4. Document error type
5. Add error handling tests

## Performance Considerations

### Memory Efficiency

- Minimal allocations during query building
- Efficient string concatenation using `strings.Builder`
- Parameter reuse where possible

### SQL Generation

- Single-pass generation where possible
- Efficient quoting logic
- Minimal conditional branching

### Error Handling

- No performance penalty for valid operations
- Fast-fail for invalid operations
- Minimal error collection overhead

## Testing Architecture

### Unit Tests

- Individual method testing
- Error condition testing
- Edge case coverage

### Integration Tests

- End-to-end query building
- Multi-dialect testing
- Database interaction testing

### Parameterized Query Tests

- Placeholder generation
- Parameter ordering
- SQL injection prevention

## See Also

- [API Reference](api_reference.md)
- [Error Handling](features/error-handling.md)
- [Parameterized Queries](features/parameterized-queries.md)
- [Design Patterns](system-patterns.md)
