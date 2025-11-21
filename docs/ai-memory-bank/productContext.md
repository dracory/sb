# Product Context: Simplified SQL Builder (SB)

Last updated at: 2025-07-22

## Problem
Many Go developers need a straightforward way to generate SQL queries programmatically. While there are several SQL builder libraries available, they often come with:
- Steep learning curves
- Unnecessary complexity for common use cases
- Heavy dependencies
- Inconsistent behavior across different database systems

## Solution
SB provides a clean, type-safe SQL builder with these key features:
- **Fluent Interface**: Chainable methods for building queries
- **Database Support**: MySQL, PostgreSQL, and SQLite with consistent API
- **Minimal Dependencies**: Primarily relies on Go's standard library
- **Type Safety**: Strongly-typed query building
- **Transaction Support**: Simplified transaction management
- **Schema Operations**: Table and column management
- **View Support**: Create and manage database views

## How it Works

### Query Building
```go
// Example: Building a SELECT query
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Select([]string{"id", "name", "email"}).
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    }).
    OrderBy("created_at", "DESC").
    Limit(10)
```

### Database Operations
```go
// Example: Database operations with transactions
db := sb.NewDatabaseFromDriver("mysql", "user:pass@/dbname")
err := db.ExecInTransaction(func(tx *sb.Database) error {
    // Perform multiple operations in transaction
    _, err := tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com")
    return err
})
```

### Schema Management
```go
// Example: Table creation
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("products").
    Column(sb.Column{
        Name:       "id",
        Type:       sb.COLUMN_TYPE_STRING,
        Length:     40,
        PrimaryKey: true,
    }).
    Column(sb.Column{
        Name:   "price",
        Type:   sb.COLUMN_TYPE_DECIMAL,
        Length: 10,
    }).
    Create()
```

## User Experience Goals

### For Developers
- **Intuitive API**: Method chaining that reads like a sentence
- **Consistent Behavior**: Same API works across supported databases
- **Helpful Errors**: Clear error messages for common mistakes
- **Comprehensive Examples**: Well-documented examples for common use cases

### For Maintainers
- **Test Coverage**: High test coverage for all features
- **Documentation**: Clear, up-to-date documentation
- **Backward Compatibility**: Careful versioning and deprecation policies
- **Performance**: Minimal overhead over raw SQL

## Integration
SB is designed to work seamlessly with:
- Existing `database/sql` connections
- Go modules
- Common web frameworks (Gin, Echo, etc.)
- Testing frameworks (Testify, etc.)
