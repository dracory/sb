---
path: llm-context.md
page-type: reference
summary: Structured documentation for AI-assisted development, providing LLMs with comprehensive context about the SB SQL Builder library.
tags: [llm, ai-context, structured-data, metadata]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# LLM Context

This document provides structured context for AI-assisted development of the SB SQL Builder library.

## Project Metadata

**Project Name**: SB (Simplified SQL Builder)
**Language**: Go
**Version**: v0.18.0
**Go Version**: 1.23.3+
**License**: See project repository
**Repository**: https://github.com/dracory/sb

## Project Purpose

SB is a lightweight, type-safe SQL builder library for Go that provides a fluent interface for constructing SQL queries across multiple database dialects (MySQL, PostgreSQL, SQLite, MSSQL).

## Key Design Principles

1. **Fluent Interface**: Method chaining for readable, expressive code
2. **Type Safety**: Strong typing to prevent runtime errors
3. **Multi-Dialect Support**: Consistent API across different databases
4. **Security First**: Parameterized queries by default
5. **Minimal Dependencies**: Primarily relies on Go standard library
6. **Error Handling**: Structured errors, no panics in production

## Architecture Overview

### Core Components

- **Builder**: Main entry point for constructing queries
- **Database**: Wrapper around sql.DB with enhanced functionality
- **Column**: Represents table columns with properties
- **Where**: Handles query conditions and subqueries
- **Join**: Represents JOIN operations with aliases
- **IndexColumn**: Represents index columns with direction and prefix length

### Design Patterns

- **Builder Pattern**: Core pattern for constructing complex SQL queries
- **Strategy Pattern**: Different SQL generation strategies per dialect
- **Facade Pattern**: Simplified interface over database/sql
- **Template Method**: Common SQL patterns with dialect variations

## Important Files

### Core Implementation

- `builder.go`: Main Builder implementation with all query methods
- `builder_where.go`: WHERE clause and subquery handling
- `builder_quote.go`: Database-specific identifier quoting
- `interfaces.go`: BuilderInterface and related interfaces
- `errors.go`: Structured error types and error handling

### Testing

- `builder_test.go`: Core Builder tests
- `parameterized_queries_test.go`: Parameterized query tests
- `integration_test.go`: Integration tests

### Documentation

- `README.md`: Main project documentation
- `docs/livewiki/`: LiveWiki documentation (this directory)
- `docs/ai-memory-bank/`: AI memory bank with project context

## Database Dialects

### Supported Dialects

- `DIALECT_MYSQL`: MySQL database support
- `DIALECT_POSTGRES`: PostgreSQL database support
- `DIALECT_SQLITE`: SQLite database support
- `DIALECT_MSSQL`: Microsoft SQL Server support

### Dialect-Specific Features

| Feature | MySQL | PostgreSQL | SQLite | MSSQL |
|---------|-------|------------|--------|-------|
| Parameter Placeholders | `?` | `$1, $2, $3` | `?` | `@p1, @p2, @p3` |
| Identifier Quoting | Backticks | Double quotes | Double quotes | Square brackets |
| IF NOT EXISTS (CREATE) | ✗ | ✓ | ✓ | ✗ |
| IF EXISTS (DROP) | ✗ | ✓ | ✓ | ✓ |
| Partial Indexes (WHERE) | ✗ | ✓ | ✓ | ✓ |
| Covering Indexes (INCLUDE) | ✗ | ✓ | ✗ | ✓ |
| FULLTEXT Index | ✓ | ✗ | ✗ | ✗ |
| GIN/GIST Index | ✗ | ✓ | ✗ | ✗ |
| Prefix Length | ✓ | ✗ | ✗ | ✗ |

## Core API Patterns

### Query Building Pattern

```go
// Standard pattern
builder := sb.NewBuilder(dialect).
    Table("table_name").
    Where(&sb.Where{Column: "col", Operator: "=", Value: "val"}).
    Select([]string{"col1", "col2"})

sql, params, err := builder.Select(columns)
```

### Error Handling Pattern

```go
// SQL generation methods return (string, []interface{}, error)
sql, params, err := builder.Select(columns)
if err != nil {
    // Handle error
}

// Schema methods return (string, error)
sql, err := builder.Create()
if err != nil {
    // Handle error
}

// Subquery methods return (BuilderInterface, error)
builder, err := builder.Exists(subquery)
if err != nil {
    // Handle error
}
```

### Parameterized Query Pattern

```go
// Default: parameterized queries
sql, params, err := builder.Select(columns)
// SQL: SELECT * FROM table WHERE col = ?
// Params: ["value"]

// Legacy mode: interpolated values
sql, _, err := builder.WithInterpolatedValues().Select(columns)
// SQL: SELECT * FROM table WHERE col = "value"
// Params: []
```

## Key Features

### 1. Parameterized Queries (v0.18.0)

**Status**: ✅ Complete
**Breaking Change**: Yes - method signatures changed from `(string, error)` to `(string, []interface{}, error)`

**Implementation Details**:
- All SQL generation methods now return parameters
- Dialect-specific placeholder generation
- Parameter ordering: WHERE first, then INSERT/UPDATE values
- NULL handling: NULL values skip placeholder generation
- Backward compatibility via `WithInterpolatedValues()`

### 2. Error Handling Refactoring

**Status**: ✅ Complete
**Breaking Change**: Yes - subquery methods now return errors

**Implementation Details**:
- Structured error types with Type and Message fields
- Error collection pattern for fluent methods
- Validation at build time (SQL generation)
- No panics in production code
- Consistent error propagation

### 3. Subquery Support

**Status**: ✅ Complete
**Breaking Change**: Yes - subquery methods return `(BuilderInterface, error)`

**Supported Types**:
- EXISTS / NOT EXISTS
- IN / NOT IN
- Comparison operators (=, >, <, >=, <=, <>)
- Correlated subqueries

### 4. JOIN Operations

**Status**: ✅ Complete
**Breaking Change**: No

**Supported Types**:
- INNER JOIN
- LEFT JOIN
- RIGHT JOIN
- FULL JOIN
- CROSS JOIN
- Table aliases
- Multiple JOINs

### 5. Enhanced Index Support

**Status**: ✅ Complete
**Breaking Change**: No

**Features**:
- Basic indexes
- Unique indexes
- Composite indexes with direction
- Partial indexes (WHERE clause)
- Covering indexes (INCLUDE clause)
- Database-specific types (FULLTEXT, GIN, GIST)
- MySQL prefix length support
- PostgreSQL schema support

## Error Types

### Standard Errors

```go
var (
    // Validation Errors
    ErrEmptyTableName     = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
    ErrEmptyColumnName    = &BuilderError{Type: "ValidationError", Message: "column name cannot be empty"}
    ErrEmptyColumnType    = &BuilderError{Type: "ValidationError", Message: "column type is required"}
    ErrEmptyOnCondition   = &BuilderError{Type: "ValidationError", Message: "ON condition cannot be empty"}
    ErrEmptyIndexName     = &BuilderError{Type: "ValidationError", Message: "index name cannot be empty"}
    ErrMissingTable       = &BuilderError{Type: "ValidationError", Message: "table name is required"}
    ErrEmptyColumns       = &BuilderError{Type: "ValidationError", Message: "columns cannot be empty"}
    
    // Argument Errors
    ErrNilSubquery        = &BuilderError{Type: "ArgumentError", Message: "subquery cannot be nil"}
)
```

### Error Handling Strategy

- **Configuration errors** → Panic (NewBuilder with invalid dialect)
- **Validation errors** → Error collection (Column, Join methods)
- **Runtime errors** → Error returns (SQL generation, subquery methods)

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

### Join Type Constants

```go
const (
    JOIN_INNER JoinType = "INNER"
    JOIN_LEFT  JoinType = "LEFT"
    JOIN_RIGHT JoinType = "RIGHT"
    JOIN_FULL  JoinType = "FULL"
    JOIN_CROSS JoinType = "CROSS"
)
```

## Recent Major Changes

### v0.18.0 (2026-03-20)

**Parameterized Queries Implementation**
- All SQL generation methods now return `(string, []interface{}, error)`
- Dialect-specific placeholder generation
- Backward compatibility via `WithInterpolatedValues()`
- 97/97 tests passing

### Error Handling Refactoring

- All panic-based error handling replaced with structured errors
- Error collection pattern for fluent methods
- Validation at build time
- Subquery methods return `(BuilderInterface, error)`

### Subquery Implementation

- EXISTS / NOT EXISTS support
- IN / NOT IN support
- Correlated subqueries
- Comprehensive test coverage (19 tests)

### JOIN Implementation

- All join types supported
- Table aliases
- Multiple JOINs
- 13 comprehensive tests

### Enhanced Index Support

- Advanced index types (FULLTEXT, GIN, GIST)
- Partial indexes (WHERE clause)
- Covering indexes (INCLUDE clause)
- 32 comprehensive tests

## Testing Status

- **Total Tests**: 97
- **Passing**: 97 (100%)
- **Test Coverage**: Comprehensive for all major features

## Dependencies

### Core Dependencies

- `database/sql` (Go standard library)
- Database drivers:
  - `github.com/go-sql-driver/mysql` (MySQL)
  - `github.com/lib/pq` (PostgreSQL)
  - `modernc.org/sqlite` (SQLite)

### Utility Dependencies

- `github.com/samber/lo` - Functional programming helpers
- `github.com/georgysavva/scany` - SQL row scanning
- `github.com/spf13/cast` - Type conversion utilities

## Development Guidelines

### Adding New Features

1. Add state to Builder struct if needed
2. Add fluent method to Builder
3. Add method to BuilderInterface
4. Implement SQL generation logic
5. Add error handling with structured errors
6. Write comprehensive tests
7. Update documentation

### Adding New Dialects

1. Implement ColumnSQLGenerator interface
2. Add dialect constant
3. Update NewBuilder switch statement
4. Add quoting logic
5. Implement dialect-specific features
6. Write dialect-specific tests
7. Update documentation

### Error Handling Guidelines

1. Use structured error types (BuilderError)
2. Define standard errors in errors.go
3. Use error collection for fluent methods
4. Validate at build time for SQL generation
5. Return errors, never panic in production
6. Document error types and messages

## Performance Considerations

- Minimal allocations during query building
- Efficient string concatenation with strings.Builder
- No reflection for type operations
- Query plan caching through parameterization
- Index join columns for JOIN operations

## Security Considerations

- Parameterized queries by default (SQL injection prevention)
- Input validation for all operations
- Type-safe operations
- No string concatenation for user input
- Structured error handling (no information leakage)

## Common Patterns

### SELECT with WHERE

```go
sql, params, err := sb.NewBuilder(dialect).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"id", "name"})
```

### INSERT with Data

```go
sql, params, err := sb.NewBuilder(dialect).
    Table("users").
    Insert(map[string]interface{}{
        "name":  "John",
        "email": "john@example.com",
    })
```

### UPDATE with WHERE

```go
sql, params, err := sb.NewBuilder(dialect).
    Table("users").
    Where(&sb.Where{Column: "id", Operator: "=", Value: 1}).
    Update(map[string]interface{}{
        "status": "inactive",
    })
```

### DELETE with WHERE

```go
sql, params, err := sb.NewBuilder(dialect).
    Table("users").
    Where(&sb.Where{Column: "id", Operator: "=", Value: 1}).
    Delete()
```

### CREATE TABLE

```go
sql, err := sb.NewBuilder(dialect).
    Table("users").
    Column(sb.Column{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true}).
    Column(sb.Column{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 255}).
    Create()
```

## Known Limitations

1. No ORM-like functionality (by design)
2. Limited support for complex window functions
3. No automatic schema migrations
4. Some database-specific features not abstracted
5. Limited support for very complex queries

## Future Enhancement Opportunities

1. Common Table Expressions (CTEs)
2. Window functions support
3. Batch operations with multiple parameter sets
4. Named parameters (:name syntax)
5. Query optimization hints
6. Index information retrieval
7. Schema introspection methods

## Documentation Structure

- `overview.md`: Main entry point and introduction
- `api_reference.md`: Complete API documentation
- `architecture.md`: System architecture and design patterns
- `features/`: Detailed feature documentation
  - `parameterized-queries.md`: Parameterized query guide
  - `error-handling.md`: Error handling system
  - `subqueries.md`: Subquery support
  - `joins.md`: JOIN operations
  - `indexes.md`: Index management
- `llm-context.md`: This document (AI consumption)
- `table_of_contents.md`: Navigation structure
- `_sidebar.md`: Docsify sidebar configuration

## AI Development Tips

When working with this codebase:

1. **Always handle errors** - Never ignore returned errors
2. **Use parameterized queries** - Default is secure
3. **Check dialect support** - Not all features work on all databases
4. **Follow error collection pattern** - For fluent method validation
5. **Test with multiple dialects** - Ensure cross-database compatibility
6. **Update documentation** - Keep docs in sync with code changes
7. **Use structured errors** - Define standard errors in errors.go
8. **Maintain backward compatibility** - Consider breaking changes carefully

## Contact and Support

- **Repository**: https://github.com/dracory/sb
- **Issues**: https://github.com/dracory/sb/issues
- **Documentation**: See this LiveWiki directory

## Version History

- **v0.18.0** (2026-03-20): Parameterized queries, error handling refactoring, subqueries, JOINs, enhanced indexes
- **v0.17.x**: Previous version without parameterized queries

## Summary

SB is a mature, well-tested SQL builder library with comprehensive error handling, parameterized query support, and multi-database compatibility. The library prioritizes security, type safety, and developer experience while maintaining simplicity and performance.
