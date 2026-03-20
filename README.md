# SB <a href="https://gitpod.io/#https://github.com/dracory/sb" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/dracory/sb/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/sb)](https://goreportcard.com/report/github.com/dracory/sb)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/sb)](https://pkg.go.dev/github.com/dracory/sb)

A simplified SQL builder with comprehensive database operations support.

For a full SQL builder functionality check: https://doug-martin.github.io/goqu

Includes a wrapper for the mainstream DB package to allow transparent working with transactions.

**Features:**
- ✅ **Subquery Support** - IN, NOT IN, EXISTS, NOT EXISTS, and comparison subqueries
- ✅ **JOIN Operations** - INNER, LEFT, RIGHT, FULL OUTER, and CROSS joins with table aliases
- ✅ **Index Management** - CREATE INDEX and DROP INDEX with database-specific options

```go
import "github.com/dracory/sb"

// Create a builder
builder := sb.NewBuilder(sb.DIALECT_MYSQL)

// Build queries
sql := builder.
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})

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

### Advanced Features

- **[JOIN Operations](docs/joins.md)** - INNER, LEFT, RIGHT, FULL OUTER, CROSS joins
- **[Index Management](docs/indexes.md)** - CREATE and DROP indexes
- **[View Management](docs/views.md)** - CREATE and DROP views
- **[Transactions](docs/transactions.md)** - Transaction management

### Examples

#### Basic Query
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})
```

#### JOIN Example
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

#### Subquery Example
```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery).
    Select([]string{"name"})
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

## Database Support

| Database | Dialect Constant | Status |
|----------|------------------|---------|
| MySQL | `sb.DIALECT_MYSQL` | ✅ Full Support |
| PostgreSQL | `sb.DIALECT_POSTGRES` | ✅ Full Support |
| SQLite | `sb.DIALECT_SQLITE` | ✅ Full Support |
| MSSQL | `sb.DIALECT_MSSQL` | ✅ Full Support |

## Recently Implemented Features

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
