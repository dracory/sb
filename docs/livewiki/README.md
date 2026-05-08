---
path: README.md
page-type: overview
summary: Default entry point for SB SQL Builder documentation, providing an introduction to the library and its capabilities.
tags: [overview, introduction, getting-started]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# SB SQL Builder - Overview

SB (Simplified SQL Builder) is a lightweight, type-safe SQL builder library for Go that provides a fluent interface for constructing SQL queries across multiple database dialects.

## What is SB?

SB is a Go library designed to make SQL query building simple, safe, and consistent across different database systems. It provides:

- **Fluent Interface**: Chainable methods that read like natural language
- **Multi-Dialect Support**: MySQL, PostgreSQL, SQLite, and MSSQL with consistent API
- **Type Safety**: Strongly-typed query building to prevent runtime errors
- **Security**: Parameterized queries by default to prevent SQL injection
- **Transaction Support**: Simplified transaction management
- **Schema Operations**: Table and column management with dialect-specific optimizations

## Quick Start

### Installation

```bash
go get github.com/dracory/sb
```

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/dracory/sb"
)

func main() {
    // Build a SELECT query
    sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{
            Column:   "status",
            Operator: "=",
            Value:    "active",
        }).
        OrderBy("created_at", "DESC").
        Limit(10).
        Select([]string{"id", "name", "email"})
    
    if err != nil {
        panic(err)
    }
    
    fmt.Println("SQL:", sql)
    fmt.Println("Params:", params)
}
```

## Key Features

### Parameterized Queries (Default)

SB uses parameterized queries by default to prevent SQL injection:

```go
sql, params, err := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userInput}).
    Select([]string{"*"})
// SQL: SELECT * FROM users WHERE email = ?
// Params: ["user_input"]
```

### Error Handling

Comprehensive error handling with structured error types:

```go
sql, err := builder.Create()
if err != nil {
    // Handle structured errors like ErrEmptyTableName, ErrNilSubquery, etc.
}
```

### Subquery Support

Full support for subqueries including EXISTS, IN, and correlated subqueries:

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery).
    Select([]string{"name"})
```

### JOIN Operations

Complete JOIN support with aliases and multiple join types:

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    LeftJoin("profiles", "users.id = profiles.user_id").
    Select([]string{"orders.*", "users.name", "profiles.avatar"})
```

### Advanced Index Support

Enhanced index creation with database-specific features:

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("documents").
    CreateIndexWithOptions("idx_documents_search", sb.IndexOptions{
        Unique:      true,
        IfNotExists: true,
        Using:       sb.INDEX_TYPE_GIN,
        Columns:     []sb.IndexColumn{{Name: "search_vector"}},
        Include:     []string{"title", "content"},
        Where:       "published = true",
    })
```

## Database Support

SB supports the following database dialects:

| Database   | Dialect Constant | Status  |
|------------|------------------|---------|
| MySQL      | `DIALECT_MYSQL`    | ✅ Full |
| PostgreSQL | `DIALECT_POSTGRES` | ✅ Full |
| SQLite     | `DIALECT_SQLITE`   | ✅ Full |
| MSSQL      | `DIALECT_MSSQL`    | ✅ Full |

## Architecture

SB uses the Builder pattern with fluent method chaining:

```
Builder
├── Table
│   ├── Column
│   │   ├── Data Type
│   │   ├── Constraints
│   │   └── Default Values
│   └── Indexes
├── Where (Conditions)
│   ├── AND/OR logic
│   ├── Comparison operators
│   └── Subquery support
├── JOIN Operations
│   ├── Multiple join types
│   ├── Table aliases
│   └── Join conditions
└── Query Operations
    ├── SELECT
    ├── INSERT
    ├── UPDATE
    ├── DELETE
    └── Schema Operations
```

## Security

SB prioritizes security with:

- **Parameterized queries by default** - SQL injection protection
- **Structured error handling** - No panics in production code
- **Type-safe operations** - Compile-time error detection
- **Input validation** - Comprehensive validation for all operations

## Performance

SB is optimized for:

- **Minimal overhead** - Efficient SQL generation
- **Query plan caching** - Better database optimization through parameterization
- **No reflection** - Static typing for better performance
- **Memory efficient** - Minimal allocations during query building

## Documentation Structure

- [API Reference](api_reference.md) - Complete API documentation
- [Features](features/) - Detailed feature documentation
- [Architecture](architecture.md) - System architecture and design patterns
- [LLM Context](llm-context.md) - Structured documentation for AI consumption

## Version Information

- **Current Version**: v0.18.0
- **Go Version**: 1.23.3+
- **License**: See project repository

## Getting Help

- Check the [API Reference](api_reference.md) for detailed method documentation
- Review [Features](features/) for specific functionality
- Consult the [LLM Context](llm-context.md) for AI-assisted development

## See Also

- [Parameterized Queries](features/parameterized-queries.md)
- [Error Handling](features/error-handling.md)
- [Subqueries](features/subqueries.md)
- [JOIN Operations](features/joins.md)
- [Index Management](features/indexes.md)
