# SB <a href="https://gitpod.io/#https://github.com/dracory/sb" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/dracory/sb/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/sb)](https://goreportcard.com/report/github.com/dracory/sb)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/sb)](https://pkg.go.dev/github.com/dracory/sb)

A simplified SQL builder with comprehensive database operations support and robust error handling.

For a full SQL builder functionality check: https://doug-martin.github.io/goqu

Includes a wrapper for the mainstream DB package to allow transparent working with transactions.

**Features:**
- ✅ **Parameterized Queries** - SQL injection protection by default with dialect-specific placeholders
- ✅ **Subquery Support** - IN, NOT IN, EXISTS, NOT EXISTS, and comparison subqueries
- ✅ **JOIN Operations** - INNER, LEFT, RIGHT, FULL OUTER, and CROSS joins with table aliases
- ✅ **Enhanced Index Management** - Advanced index types (GIN, FULLTEXT, partial, covering) with database-specific options
- ✅ **Zero-Panic Error Handling** - All errors returned gracefully, no panics anywhere
- ✅ **Fluent API with Error Collection** - Chain operations with comprehensive error validation

## Security

🔒 **Parameterized Queries by Default:** SB now generates parameterized queries with SQL injection protection. Values are separated from SQL and safely handled by database drivers.

```go
// ✅ Secure - parameterized queries (default)
sql, params, err := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail}).
    Select([]string{"*"})
// SQL: SELECT * FROM users WHERE email = ?
// Params: ["user@example.com"]

// ✅ Legacy mode - interpolated values (backward compatible)
sql, _, err := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail}).
    WithInterpolatedValues().
    Select([]string{"*"})
// SQL: SELECT * FROM users WHERE email = "user@example.com"
```

See [Security Guide](docs/security.md) for detailed safety information and best practices.

```go
// ✅ Safe - values are properly quoted
sql := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail}).
    Select([]string{"*"})

// ❌ Dangerous - never concatenate user input
sql := "SELECT * FROM users WHERE email = '" + userEmail + "'"
```

```go
import "github.com/dracory/sb"

// Create a builder
builder := sb.NewBuilder(sb.DIALECT_MYSQL)

// Build parameterized queries with error handling
sql, params, err := builder.
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})

if err != nil {
    // Handle error gracefully
    log.Fatal("SQL generation error:", err)
}

// Execute with database (parameterized)
myDb := sb.NewDatabaseFromDriver("sqlite3", "test.db")
ctx := context.Background()
result, err := myDb.Exec(ctx, sql, params...)  // Pass params separately
```

## Documentation

### Core Operations

- **[Database Connection](docs/database-connection.md)** - Setup and configuration
- **[Table Operations](docs/tables.md)** - CREATE, DROP, TRUNCATE tables
- **[SELECT Queries](docs/selects.md)** - SELECT statements and subqueries
- **[INSERT Operations](docs/inserts.md)** - INSERT statements
- **[UPDATE & DELETE](docs/updates.md)** - UPDATE and DELETE operations
- **[Error Handling](docs/error-handling.md)** - Comprehensive error management strategy

### Advanced Features

- **[JOIN Operations](docs/joins.md)** - INNER, LEFT, RIGHT, FULL OUTER, CROSS joins
- **[Index Management](docs/indexes.md)** - CREATE and DROP indexes
- **[View Management](docs/views.md)** - CREATE and DROP views
- **[Transactions](docs/transactions.md)** - Transaction management

### Examples

#### Basic Query (Parameterized)
```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})

if err != nil {
    log.Fatal("SQL generation error:", err)
}

// Execute: db.Exec(sql, params...)
// SQL: SELECT `name`, `email` FROM `users` WHERE `status` = ?
// Params: ["active"]
```

#### JOIN Example (Parameterized)
```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
    Select([]string{"orders.*", "users.name"})

if err != nil {
    log.Fatal("SQL generation error:", err)
}

// Execute: db.Query(sql, params...)
// SQL: SELECT `orders`.*, `users`.`name` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id WHERE `orders`.`status` = ?
// Params: ["active"]
```

#### Subquery Example (Parameterized)
```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery).
    Select([]string{"name"})

if err != nil {
    log.Fatal("SQL generation error:", err)
}

// Execute: db.Query(sql, params...)
// SQL: SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > ?)
// Params: ["1000"]
```

#### Enhanced Index Support
```go
// Simple unique index
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    CreateUniqueIndex("idx_users_email", "email")

// PostgreSQL GIN index with covering columns
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("documents").
    CreateIndexWithOptions("idx_search", sb.IndexOptions{
        Unique:      true,
        IfNotExists: true,
        Using:       sb.INDEX_TYPE_GIN,
        Columns:     []sb.IndexColumn{{Name: "search_vector"}},
        Include:     []string{"title", "content"},
        Where:       "published = true",
        Storage:     "fillfactor=90",
    })

// MySQL FULLTEXT index with prefix and comment
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("articles").
    CreateIndexWithOptions("idx_content", sb.IndexOptions{
        Using:   sb.INDEX_TYPE_FULLTEXT,
        Columns: []sb.IndexColumn{
            {Name: "title", Length: 100},
            {Name: "content", Length: 255},
        },
        Comment: "Full-text search index",
    })

// Enhanced DROP with schema support
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
        IfExists: true,
        Schema:   "public",
    })
```

#### Error Handling Example
```go
// All SQL generation methods return (string, []interface{}, error)
sql, params, err := builder.
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: "test@example.com"}).
    Select([]string{"*"})

if err != nil {
    switch err.Error() {
    case "ValidationError: no table specified":
        // Handle missing table error
    case "ValidationError: unsupported dialect: unknown":
        // Handle unsupported dialect error
    default:
        // Handle other errors
        log.Fatal("SQL error:", err)
    }
}

// Execute safely with parameters
rows, err := db.Query(sql, params...)
```

## Migration Guide

### Breaking Changes in v0.18.0

SB now uses parameterized queries by default for improved security. This is a **breaking change** that affects method signatures.

#### Before (v0.17.x)
```go
sql, err := builder.
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"*"})

db.Exec(sql)  // Direct execution
```

#### After (v0.18.0)
```go
sql, params, err := builder.
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"*"})

db.Exec(sql, params...)  // Pass parameters separately
```

### Migration Strategies

#### 1. Quick Migration (Recommended)
Update your code to handle the new 3-value return:

```go
// Before
sql, err := builder.Select([]string{"*"})
db.Exec(sql)

// After  
sql, params, err := builder.Select([]string{"*"})
db.Exec(sql, params...)
```

#### 2. Gradual Migration
Use `WithInterpolatedValues()` for backward compatibility:

```go
// Legacy mode - no changes needed
sql, _, err := builder.
    WithInterpolatedValues().
    Select([]string{"*"})
db.Exec(sql)  // Same as before
```

#### 3. Database-Specific Placeholders

| Database | Placeholder Format | Example |
|----------|------------------|---------|
| MySQL | `?` | `WHERE id = ?` |
| PostgreSQL | `$1, $2, $3` | `WHERE id = $1` |
| SQLite | `?` | `WHERE id = ?` |
| MSSQL | `@p1, @p2, @p3` | `WHERE id = @p1` |

### Benefits of Migration

✅ **Security**: SQL injection protection by default  
✅ **Performance**: Better query plan caching  
✅ **Compatibility**: Works with all major Go database drivers  
✅ **Flexibility**: Optional legacy mode available  

### Need Help?

- Check the [parameterized queries test file](parameterized_queries_test.go) for examples
- Review the [implementation documentation](docs/ai-memory-bank/parameterized-queries-implementation.md)
- Open an issue for migration questions

#### Transaction Example
```go
err := myDb.ExecInTransaction(ctx, func(tx *sb.Database) error {
    if _, err := tx.Exec(ctx, sql1); err != nil {
        return err
    }
    if _, err := tx.Exec(ctx, sql2); err != nil {
        return err
    }
    return nil
})
```

## Error Handling

SB uses a **zero-panic error handling strategy**. All errors are collected during fluent chaining and returned when SQL is generated.

### Error Types

- **ValidationError** - Invalid parameters (missing table, empty conditions, etc.)
- **ArgumentError** - Invalid arguments (nil subqueries, etc.)

### Error Collection Pattern

```go
// Errors are collected during chaining
builder := sb.NewBuilder("unknown") // Unsupported dialect error collected
builder.InnerJoin("users", "")    // Empty JOIN condition error collected

// All errors returned when SQL is generated
sql, err := builder.Select([]string{"*"})
// Returns: "ValidationError: unsupported dialect: unknown"
```

### Best Practices

```go
// ✅ Always check for errors
sql, err := builder.Select([]string{"*"})
if err != nil {
    return fmt.Errorf("failed to generate SQL: %w", err)
}

// ✅ Handle specific error cases
if strings.Contains(err.Error(), "unsupported dialect") {
    // Handle dialect error
}

// ✅ Use structured error handling
switch {
case strings.Contains(err.Error(), "no table specified"):
    // Handle missing table
case strings.Contains(err.Error(), "ON condition cannot be empty"):
    // Handle empty JOIN condition
default:
    return err
}
```

## Database Support

| Database | Dialect Constant | Status |
|----------|------------------|---------|
| MySQL | `sb.DIALECT_MYSQL` | ✅ Full Support |
| PostgreSQL | `sb.DIALECT_POSTGRES` | ✅ Full Support |
| SQLite | `sb.DIALECT_SQLITE` | ✅ Full Support |
| MSSQL | `sb.DIALECT_MSSQL` | ✅ Full Support |

## Recently Implemented Features

- ✅ **Parameterized Queries** - SQL injection protection by default with dialect-specific placeholders
- ✅ **Zero-Panic Error Handling** - Complete error collection system, no panics anywhere
- ✅ **Subquery Support** - IN, NOT IN, EXISTS, NOT EXISTS, and comparison subqueries with correlation
- ✅ **JOIN Support** - INNER, LEFT, RIGHT, FULL OUTER, and CROSS joins with table aliases
- ✅ **Index Management** - Complete CREATE INDEX and DROP INDEX functionality
- ✅ **Table Truncation** - TRUNCATE TABLE support with database-specific options
- ✅ **Multi-dialect Support** - MySQL, PostgreSQL, SQLite, and MSSQL compatibility

## Developers

```sh
podman run -it --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=test -e MYSQL_DATABASE=test -e MYSQL_USER=test -e MYSQL_PASSWORD=test mysql:latest
```

```sh
podman run -it --rm -p 5432:5432 -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -e POSTGRES_USER=test postgres:latest
```

## Similar Projects

- https://doug-martin.github.io/goqu - Best SQL Builder for Golang
- https://github.com/elgris/golang-sql-builder-benchmark
- https://github.com/es-code/gql
- https://github.com/cosiner/go-sqldb
- https://github.com/simukti/sqldb-logger
- https://github.com/elgs/gosqlcrud
- https://github.com/nandrechetan/gomb

## TODO

- github.com/stapelberg/postgrestest
