# Technical Context: Simplified SQL Builder (SB)

Last updated at: 2025-07-22

## Technologies Used
- Go programming language (1.23.3+)
- `database/sql` package (Go's standard database library)
- Database drivers:
  - MySQL: `github.com/go-sql-driver/mysql`
  - PostgreSQL: `github.com/lib/pq`
  - SQLite: `"modernc.org/sqlite"`
- Additional utility packages:
  - `github.com/samber/lo` for functional programming helpers
  - `github.com/gouniverse/*` for various utility functions
  - `github.com/georgysavva/scany` for scanning SQL rows into structs

## Development Setup
- Go development environment
- A suitable IDE or text editor
- Access to MySQL, PostgreSQL, or SQLite databases for testing

## Technical Constraints
- The library should have minimal dependencies to reduce project size and complexity.
- The generated SQL should be compatible with the target database dialects.
- The database wrapper should not introduce significant performance overhead.
- Transaction management should be reliable and consistent across different database drivers.

## Dependencies
### Core Dependencies
- `database/sql` (Go's standard database library)
- Database drivers:
  - `github.com/go-sql-driver/mysql` for MySQL
  - `github.com/lib/pq` for PostgreSQL
  - `_ "modernc.org/sqlite"` for SQLite

### Utility Dependencies
- `github.com/samber/lo` - Functional programming helpers
- `github.com/georgysavva/scany` - SQL row scanning
- `github.com/gouniverse/*` - Various utility packages
- `github.com/spf13/cast` - Type conversion utilities

### Development Dependencies
- Testing frameworks (implicit with Go's standard library)
- `github.com/stretchr/testify` for assertions

## Date Handling

### NULL Date Implementation
The library uses `0002-01-01` as the NULL date value (rather than `0001-01-01`) for better compatibility across different database systems.

**Rationale**:
- Avoids issues with MySQL and other systems that may reject `0001-01-01`
- Less likely to be confused with system minimum values
- Lower chance of collision with default or minimum date values
- Still obviously an invalid date for most real-world applications

**Implementation**:
```go
// NULL_DATE represents an invalid/empty date value
const NULL_DATE = "0002-01-01"
const NULL_DATETIME = "0002-01-01 00:00:00"
```

Full documentation is available in `docs/null_date.md`
