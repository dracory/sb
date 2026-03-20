# Error Handling Guidelines

## Philosophy

The SB SQL builder should use error returns for SQL generation methods and reserve panics only for programmer errors that cannot be recovered from.

## Panic vs Error Return

### Use Panic For:
- **Programmer errors** that should never occur in production
- **Invalid configuration** that prevents the library from functioning
- **Type assertion failures** when interface contracts are violated
- **Unsupported database dialects** when no fallback is possible

### Use Error Return For:
- **SQL generation methods** that can fail due to invalid input
- **Missing required fields** (table names, column names)
- **Invalid method arguments** that can be corrected by the caller
- **Validation failures** that provide actionable feedback

## Current State

### Methods That Currently Panic (Should Return Errors)
```go
// These should return (string, error):
CreateIndex() string
DropIndex() string  
DropIndexIfExists() string
DropIndexWithSchema() string
Truncate() string
TruncateWithOptions() string
Select() string
Create() string
CreateIfNotExists() string
Drop() string
DropIfExists() string
Delete() string
Insert() string
Update() string
```

### Methods Already Returning Errors (Correct Pattern)
```go
// These already follow the correct pattern:
TableColumnAdd() (string, error)
TableColumnChange() (string, error)
TableColumnDrop() (string, error)
TableColumnExists() (string, []any, error)
TableColumnRename() (string, error)
TableRename() (string, error)
```

### Builder Methods (Keep Panic - Programmer Errors)
```go
// These can keep panics as they are programmer errors:
Column() BuilderInterface
Join() BuilderInterface
Where() BuilderInterface
OrderBy() BuilderInterface
Limit() BuilderInterface
Offset() BuilderInterface
GroupBy() BuilderInterface
```

## Error Types

### Standard Error Types
```go
// Common error types
var (
    ErrEmptyTableName     = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
    ErrEmptyColumnName    = &BuilderError{Type: "ValidationError", Message: "column name cannot be empty"}
    ErrEmptyIndexName     = &BuilderError{Type: "ValidationError", Message: "index name cannot be empty"}
    ErrInvalidDialect     = &BuilderError{Type: "ConfigurationError", Message: "invalid database dialect"}
    ErrMissingTable       = &BuilderError{Type: "ValidationError", Message: "no table specified"}
    ErrEmptyColumns       = &BuilderError{Type: "ValidationError", Message: "columns cannot be empty"}
    ErrNilSubquery        = &BuilderError{Type: "ValidationError", Message: "subquery cannot be nil"}
    ErrEmptyOnCondition   = &BuilderError{Type: "ValidationError", Message: "ON condition cannot be empty"}
)
```

### Custom Error Type
```go
type BuilderError struct {
    Type    string
    Message string
    Cause   error
}

func (e *BuilderError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
    }
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *BuilderError) Unwrap() error {
    return e.Cause
}
```

## Implementation Patterns

### SQL Generation Methods
```go
// Before (panic):
func (b *Builder) CreateIndex(indexName string, columnName ...string) string {
    if indexName == "" {
        panic("index name cannot be empty!")
    }
    // ... generate SQL
}

// After (error return):
func (b *Builder) CreateIndex(indexName string, columnName ...string) (string, error) {
    if indexName == "" {
        return "", ErrEmptyIndexName
    }
    if b.sqlTableName == "" {
        return "", ErrMissingTable
    }
    // ... generate SQL
    return sql, nil
}
```

### Builder Methods (Keep Panic)
```go
// These can keep panics as they represent programmer errors:
func (b *Builder) Column(column Column) BuilderInterface {
    if column.Name == "" {
        panic("column name is required") // Programmer error
    }
    // ... continue
    return b
}
```

## Migration Strategy

### Phase 1: Add Error Types
- Create `errors.go` with custom error types
- Define standard error variables

### Phase 2: Update SQL Generation Methods
- Change method signatures to return `(string, error)`
- Replace panics with error returns
- Update interface accordingly

### Phase 3: Update Tests
- Update test expectations to handle errors
- Add error case testing

### Phase 4: Update Documentation
- Document error behavior in method comments
- Update examples to show error handling

## Backward Compatibility

This is a **breaking change** that requires:
- Major version bump (v1.0.0)
- Migration guide for users
- Deprecation warnings in transition period

## Examples

### Error Handling Pattern
```go
sql, err := builder.Table("users").CreateIndex("idx_email", "email")
if err != nil {
    // Handle error appropriately
    return fmt.Errorf("failed to create index: %w", err)
}

// Execute SQL
_, err = db.Exec(sql)
if err != nil {
    return fmt.Errorf("failed to execute index creation: %w", err)
}
```

### Builder Method (No Error Handling Needed)
```go
// Builder methods chain without error handling
sql := builder.
    Table("users").
    Where(&Where{Column: "id", Operator: ">", Value: "100"}).
    OrderBy("name", "ASC").
    Limit(10).
    Select([]string{"name", "email"}) // Only Select returns error
```

## Benefits

1. **System Stability**: No unexpected panics from invalid input
2. **Better Error Messages**: Structured errors with context
3. **Graceful Degradation**: Callers can handle errors appropriately
4. **Testing**: Easier to test error conditions
5. **Debugging**: Clear error types and messages

## Migration Timeline

- **v0.18.x**: Add error types, prepare for migration
- **v0.19.x**: Implement new error-returning methods alongside old ones
- **v1.0.0**: Remove old panic-based methods, breaking change
