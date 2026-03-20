# Golang Declarative SQL Builders

Here's a list of popular SQL builders in the Go ecosystem that follow a declarative approach:

## 1. [go-sqlbuilder](https://github.com/huandu/go-sqlbuilder)
- **Features**:
  - Supports MySQL, PostgreSQL, SQLite, and more
  - Type-safe query building
  - Supports both standard and named parameter binding
  - Extensible with custom SQL dialects
  - Supports subqueries and complex joins

## 2. [goqu](https://github.com/doug-martin/goqu)
- **Features**:
  - Full-featured SQL builder and query library
  - Supports CRUD operations with a fluent API
  - Supports transactions and prepared statements
  - Extensible with custom SQL dialects
  - Comprehensive test coverage

## 3. [squirrel](https://github.com/Masterminds/squirrel)
- **Features**:
  - Fluent API for building SQL queries
  - Supports PostgreSQL, MySQL, and SQLite
  - Composable query building
  - Supports prepared statements
  - Well-documented with examples

## 4. [dbr](https://github.com/gocraft/dbr)
- **Features**:
  - Fast and simple SQL builder
  - Supports PostgreSQL, MySQL, and SQLite
  - Built-in connection pooling
  - Supports transactions and prepared statements
  - Lightweight and focused API

## 5. [sqlx](https://github.com/jmoiron/sqlx)
- **Features**:
  - Extends the standard `database/sql` package
  - Named parameter support
  - Struct scanning
  - Supports PostgreSQL, MySQL, and SQLite
  - Lightweight and minimal abstraction

## 6. [bun](https://bun.uptrace.dev/)
- **Features**:
  - SQL-first query builder
  - Supports PostgreSQL, MySQL, and SQLite
  - Type-safe query building
  - Supports migrations
  - Built on top of `database/sql`

## 7. [gorm](https://gorm.io/)
- **Features**:
  - Full-featured ORM with SQL builder
  - Supports multiple databases
  - Migrations
  - Hooks and callbacks
  - Associations and relationships

## 8. [ent](https://entgo.io/)
- **Features**:
  - Entity framework for Go
  - Schema as code
  - Type-safe queries
  - GraphQL and gRPC support
  - Built-in migration system

## 9. [sqlc](https://github.com/sqlc-dev/sqlc)
- **Features**:
  - Generates type-safe Go from SQL
  - Supports PostgreSQL and MySQL
  - Compile-time query checking
  - No ORM, just plain SQL with type safety

## 10. [squid](https://github.com/abiosoft/squid)
- **Features**:
  - Minimal SQL builder
  - Supports PostgreSQL and MySQL
  - Lightweight and fast
  - Simple and intuitive API

## Comparison Table

| Builder | Type | Key Features |
|---------|------|-------------|
| go-sqlbuilder | SQL Builder | Dialect support, Type-safe, Named parameters |
| goqu | SQL Builder | Full-featured, Fluent API, Transactions |
| squirrel | SQL Builder | Fluent API, Composable queries |
| dbr | SQL Builder | Connection pooling, Simple API |
| sqlx | SQL Extensions | Lightweight, Named parameters |
| bun | SQL Builder | Type-safe, Migrations |
| GORM | ORM | Full ORM, Migrations, Hooks |
| ent | Entity Framework | Schema as code, Type-safe, GraphQL |
| sqlc | Code Generator | Type-safe from SQL, No runtime |
| squid | SQL Builder | Minimal, Fast, Simple |

## When to Use What

- **For simple applications**: Consider `sqlx` or `squid` for minimal overhead
- **For complex queries**: `goqu` or `go-sqlbuilder` offer powerful query building
- **For ORM features**: `GORM` or `ent` provide full ORM capabilities
- **For type safety**: `sqlc` generates type-safe code from your SQL
- **For schema migrations**: `bun` or `GORM` have built-in migration support

## References

- [Awesome Go SQL Builders](https://awesome-go.com/sql-query-builders/)
- [Go Database/SQL Tutorial](http://go-database-sql.org/)
- [SQLBoiler](https://github.com/volatiletech/sqlboiler) - Another popular ORM alternative