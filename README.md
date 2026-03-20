# SB <a href="https://gitpod.io/#https://github.com/dracory/sb" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/dracory/sb/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/sb)](https://goreportcard.com/report/github.com/dracory/sb)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/sb)](https://pkg.go.dev/github.com/dracory/sb)

A simplified SQL builder with comprehensive database operations support and robust error handling.

For a full SQL builder functionality check: https://doug-martin.github.io/goqu

Includes a wrapper for the mainstream DB package to allow transparent working with transactions.

**Features:**
- ✅ **Subquery Support** - IN, NOT IN, EXISTS, NOT EXISTS, and comparison subqueries
- ✅ **JOIN Operations** - INNER, LEFT, RIGHT, FULL OUTER, and CROSS joins with table aliases
- ✅ **Index Management** - CREATE INDEX and DROP INDEX with database-specific options
- ✅ **Zero-Panic Error Handling** - All errors returned gracefully, no panics anywhere
- ✅ **Fluent API with Error Collection** - Chain operations with comprehensive error validation

## Security

⚠️ **Important:** SB generates SQL strings with value quoting. While this provides basic protection against SQL injection, consider using parameterized queries when available (planned v0.18.0).

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

// Build queries with error handling
sql, err := builder.
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})

if err != nil {
    // Handle error gracefully
    log.Fatal("SQL generation error:", err)
}

// Execute with database
myDb := sb.NewDatabaseFromDriver("sqlite3", "test.db")
ctx := context.Background()
result, err := myDb.Exec(ctx, sql)
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

#### Basic Query
```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})

if err != nil {
    log.Fatal("SQL generation error:", err)
}
```

#### JOIN Example
```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})

if err != nil {
    log.Fatal("SQL generation error:", err)
}
```

#### Subquery Example
```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery).
    Select([]string{"name"})

if err != nil {
    log.Fatal("SQL generation error:", err)
}
```

#### Error Handling Example
```go
// All SQL generation methods return (string, error)
sql, err := builder.
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
```

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
