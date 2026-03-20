# Error Handling Guide

SB uses a **zero-panic error handling strategy**. All errors are collected during fluent chaining and returned when SQL is generated. This ensures the library never panics and provides predictable error handling for production use.

## Overview

- **No Panics**: SB never panics under any circumstances
- **Error Collection**: Errors are collected during method chaining
- **Graceful Returns**: All SQL generation methods return `(string, error)`
- **Structured Errors**: Consistent error types and messages

## Error Types

### ValidationError
Used for invalid parameters or validation failures:
- Missing table name
- Empty conditions
- Invalid column names
- Unsupported dialects

```go
"ValidationError: no table specified"
"ValidationError: ON condition cannot be empty"
"ValidationError: unsupported dialect: unknown"
```

### ArgumentError
Used for invalid arguments:
- Nil subqueries
- Invalid data types

```go
"ArgumentError: subquery cannot be nil"
```

## Error Collection Pattern

### How It Works

1. **Collection Phase**: Errors are collected during fluent chaining
2. **Validation Phase**: Errors are checked when SQL is generated
3. **Return Phase**: First collected error is returned

```go
// Step 1: Errors collected during chaining
builder := sb.NewBuilder("unknown") // Error collected: unsupported dialect
builder.InnerJoin("users", "")    // Error collected: empty ON condition
builder.Table("")                 // Error collected: empty table name

// Step 2: All errors returned when SQL is generated
sql, err := builder.Select([]string{"*"})
// Returns: "ValidationError: unsupported dialect: unknown"
```

### Error Collection Methods

These methods collect errors but don't return them immediately:

```go
// Builder methods that collect errors
NewBuilder(dialect string) *Builder
Table(tableName string) BuilderInterface
InnerJoin(table, condition string) BuilderInterface
LeftJoin(table, condition string) BuilderInterface
RightJoin(table, condition string) BuilderInterface
JoinWithAlias(joinType, table, alias, condition string) BuilderInterface
Column(column Column) BuilderInterface
```

### Error Validation Methods

These methods validate and return collected errors:

```go
// SQL generation methods that validate and return errors
Select(columns []string) (string, error)
Insert(values map[string]string) (string, error)
Update(values map[string]string) (string, error)
Delete() (string, error)
Create() (string, error)
Drop() (string, error)
Truncate() (string, error)
DropIndex(name string) (string, error)
CreateIndex(name string, columns ...string) (string, error)
```

## Best Practices

### Always Check for Errors

```go
// ✅ GOOD: Always check for errors
sql, err := builder.Select([]string{"*"})
if err != nil {
    return fmt.Errorf("failed to generate SQL: %w", err)
}

// ❌ BAD: Ignoring errors
sql, _ := builder.Select([]string{"*"}) // Could be empty string!
```

### Handle Specific Error Cases

```go
sql, err := builder.Select([]string{"*"})
if err != nil {
    switch {
    case strings.Contains(err.Error(), "unsupported dialect"):
        return fmt.Errorf("database dialect not supported: %w", err)
    case strings.Contains(err.Error(), "no table specified"):
        return fmt.Errorf("table name required: %w", err)
    case strings.Contains(err.Error(), "ON condition cannot be empty"):
        return fmt.Errorf("JOIN condition required: %w", err)
    default:
        return fmt.Errorf("SQL generation failed: %w", err)
    }
}
```

### Validate Early

```go
// ✅ GOOD: Validate builder configuration early
builder := sb.NewBuilder(dialect)
if err := builder.validateAndReturnError(); err != nil {
    return fmt.Errorf("invalid builder configuration: %w", err)
}

// Continue with confidence that builder is valid
sql, err := builder.Table("users").Select([]string{"*"})
```

### Error Wrapping

```go
// ✅ GOOD: Wrap errors with context
sql, err := builder.Select([]string{"*"})
if err != nil {
    return fmt.Errorf("failed to generate user query: %w", err)
}

// ✅ GOOD: Use structured error types
var validationErr *sb.ValidationError
if errors.As(err, &validationErr) {
    // Handle validation errors specifically
}
```

## Common Error Scenarios

### Unsupported Dialect

```go
builder := sb.NewBuilder("unknown")
_, err := builder.Select([]string{"*"})
// Error: "ValidationError: unsupported dialect: unknown"
```

### Missing Table

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL)
_, err := builder.Select([]string{"*"})
// Error: "ValidationError: no table specified"
```

### Empty JOIN Condition

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "") // Empty condition

_, err := builder.Select([]string{"*"})
// Error: "ValidationError: ON condition cannot be empty"
```

### Nil Subquery

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(nil) // Nil subquery

_, err := builder.Select([]string{"*"})
// Error: "ArgumentError: subquery cannot be nil"
```

### Empty Index Name

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    DropIndex("") // Empty index name

_, err := builder.Select([]string{"*"})
// Error: "ValidationError: index name cannot be empty"
```

## Migration from Panic-Based Error Handling

### Before (Panic-based)

```go
// Old approach - could panic
func someFunction() {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("Panic recovered: %v", r)
        }
    }()
    
    sql := sb.NewBuilder("unknown").Select([]string{"*"})
    // Might panic!
}
```

### After (Error-based)

```go
// New approach - predictable error handling
func someFunction() error {
    sql, err := sb.NewBuilder("unknown").Select([]string{"*"})
    if err != nil {
        return fmt.Errorf("SQL generation failed: %w", err)
    }
    
    // Use sql safely
    return nil
}
```

## Error Handling in Different Contexts

### Web Applications

```go
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "id", Operator: "=", Value: r.URL.Query().Get("id")})
    
    sql, err := builder.Select([]string{"name", "email"})
    if err != nil {
        http.Error(w, "Database query error", http.StatusInternalServerError)
        log.Printf("SQL generation error: %v", err)
        return
    }
    
    // Execute query safely
}
```

### CLI Applications

```go
func exportUsers() error {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "active", Operator: "=", Value: true})
    
    sql, err := builder.Select([]string{"*"})
    if err != nil {
        return fmt.Errorf("failed to generate user export query: %w", err)
    }
    
    fmt.Printf("Generated SQL: %s\n", sql)
    return nil
}
```

### Library Code

```go
func FindUser(db *sql.DB, id string) (*User, error) {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "id", Operator: "=", Value: id})
    
    sql, err := builder.Select([]string{"name", "email"})
    if err != nil {
        return nil, fmt.Errorf("failed to generate user query: %w", err)
    }
    
    // Execute query and return user
    return executeUserQuery(db, sql)
}
```

## Testing Error Handling

### Test Error Cases

```go
func TestUnsupportedDialect(t *testing.T) {
    builder := sb.NewBuilder("unknown")
    _, err := builder.Select([]string{"*"})
    
    if err == nil {
        t.Fatal("Expected error for unsupported dialect")
    }
    
    expected := "ValidationError: unsupported dialect: unknown"
    if err.Error() != expected {
        t.Fatalf("Expected %s but got: %v", expected, err)
    }
}
```

### Test Error Recovery

```go
func TestErrorCollection(t *testing.T) {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        InnerJoin("orders", "") // Empty condition
    
    // Multiple errors can be collected
    _, err := builder.Select([]string{"*"})
    if err == nil {
        t.Fatal("Expected error for empty JOIN condition")
    }
    
    // First error is returned
    expected := "ValidationError: ON condition cannot be empty"
    if err.Error() != expected {
        t.Fatalf("Expected %s but got: %v", expected, err)
    }
}
```

## Performance Considerations

### Error Collection Overhead

The error collection system has minimal performance impact:

- **Memory**: Small slice to store errors (usually empty or 1-2 items)
- **CPU**: Simple append operations during chaining
- **Validation**: Single check during SQL generation

### Optimization Tips

```go
// ✅ GOOD: Early validation to avoid unnecessary work
builder := sb.NewBuilder(dialect)
if err := builder.validateAndReturnError(); err != nil {
    return err // Fail fast
}

// ✅ GOOD: Reuse validated builders
builder := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users")
sql1, err1 := builder.Select([]string{"name"})
sql2, err2 := builder.Select([]string{"email"})
```

## Implementation Details

### Error Types

```go
// Standard error types
var (
    ErrEmptyTableName     = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
    ErrEmptyColumnName    = &BuilderError{Type: "ValidationError", Message: "column name cannot be empty"}
    ErrEmptyIndexName     = &BuilderError{Type: "ValidationError", Message: "index name cannot be empty"}
    ErrMissingTable       = &BuilderError{Type: "ValidationError", Message: "no table specified"}
    ErrEmptyColumns       = &BuilderError{Type: "ValidationError", Message: "columns cannot be empty"}
    ErrNilSubquery        = &BuilderError{Type: "ValidationError", Message: "subquery cannot be nil"}
    ErrEmptyOnCondition   = &BuilderError{Type: "ValidationError", Message: "ON condition cannot be empty"}
)

// Custom error type
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
```

### Error Collection Implementation

```go
type Builder struct {
    // ... other fields ...
    sqlErrors []error // Collect errors during fluent chaining
}

// Error collection methods
func (b *Builder) addError(err error) {
    if err != nil {
        b.sqlErrors = append(b.sqlErrors, err)
    }
}

// Error validation methods
func (b *Builder) validateAndReturnError() error {
    if len(b.sqlErrors) == 0 {
        return nil
    }
    return b.sqlErrors[0] // Return first collected error
}
```

## Summary

SB's error handling provides:

- **Predictable Behavior**: No panics, only structured errors
- **Clear Error Messages**: Descriptive error types and messages
- **Easy Integration**: Standard Go error handling patterns
- **Production Ready**: Robust error handling for real applications

By following these patterns, you can build reliable applications that handle database query generation errors gracefully and predictably.
