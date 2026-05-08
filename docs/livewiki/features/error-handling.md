---
path: features/error-handling.md
page-type: tutorial
summary: Comprehensive error handling system in SB SQL Builder with structured error types and validation patterns.
tags: [error-handling, validation, structured-errors]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# Error Handling

SB provides a comprehensive error handling system with structured error types, validation patterns, and graceful error propagation.

## Overview

SB uses a sophisticated error handling approach:

- **Structured Error Types**: Clear, typed errors with context
- **Error Collection**: Fluent methods collect errors during chaining
- **Validation at Build Time**: Errors are validated when SQL is generated
- **No Panics in Production**: All errors are returned, not panicked

## Error Types

### BuilderError

The base error type for all builder operations:

```go
type BuilderError struct {
    Type    string
    Message string
}

func (e *BuilderError) Error() string {
    return fmt.Sprintf("%s: %s", e.Type, e.Message)
}
```

### Standard Error Types

```go
var (
    // Validation Errors
    ErrEmptyTableName = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
    ErrEmptyColumnName = &BuilderError{Type: "ValidationError", Message: "column name cannot be empty"}
    ErrEmptyColumnType = &BuilderError{Type: "ValidationError", Message: "column type is required"}
    ErrEmptyOnCondition = &BuilderError{Type: "ValidationError", Message: "ON condition cannot be empty"}
    ErrEmptyIndexName = &BuilderError{Type: "ValidationError", Message: "index name cannot be empty"}
    ErrMissingTable = &BuilderError{Type: "ValidationError", Message: "table name is required"}
    ErrEmptyColumns = &BuilderError{Type: "ValidationError", Message: "columns cannot be empty"}
    
    // Argument Errors
    ErrNilSubquery = &BuilderError{Type: "ArgumentError", Message: "subquery cannot be nil"}
)
```

## Error Collection Pattern

### Fluent Methods Collect Errors

Fluent methods collect errors silently during chaining:

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Column(sb.Column{Name: "", Type: sb.COLUMN_TYPE_STRING}) // Error collected
```

### Validation at Build Time

Errors are validated when SQL is generated:

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    // Handle ErrEmptyColumnName
    return "", nil, err
}
```

### Multiple Error Collection

Multiple errors can be collected in a single chain:

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Column(sb.Column{Name: "", Type: sb.COLUMN_TYPE_STRING}).    // Error 1
    Column(sb.Column{Name: "email", Type: ""})                  // Error 2

sql, params, err := builder.Select([]string{"*"})
// err will be the first collected error
```

## SQL Generation Error Handling

### Select Method

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    // Handle error
    return "", nil, err
}
// Use sql and params
```

### Insert Method

```go
sql, params, err := builder.Insert(data)
if err != nil {
    // Handle error
    return "", nil, err
}
// Use sql and params
```

### Update Method

```go
sql, params, err := builder.Update(data)
if err != nil {
    // Handle error
    return "", nil, err
}
// Use sql and params
```

### Delete Method

```go
sql, params, err := builder.Delete()
if err != nil {
    // Handle error
    return "", nil, err
}
// Use sql and params
```

## Schema Operation Error Handling

### Create Method

```go
sql, err := builder.Create()
if err != nil {
    // Handle error
    return "", err
}
// Use sql
```

### Drop Method

```go
sql, err := builder.Drop()
if err != nil {
    // Handle error
    return "", err
}
// Use sql
```

### Truncate Method

```go
sql, err := builder.Truncate()
if err != nil {
    // Handle error
    return "", err
}
// Use sql
```

## Subquery Error Handling

### Exists Method

```go
builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(subquery)
if err != nil {
    // Handle ErrNilSubquery
    return nil, err
}
```

### InSubquery Method

```go
builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery)
if err != nil {
    // Handle ErrNilSubquery
    return nil, err
}
```

## Error Type Checking

### Check for Specific Error Type

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    if builderErr, ok := err.(*sb.BuilderError); ok {
        switch builderErr.Type {
        case "ValidationError":
            fmt.Println("Validation error:", builderErr.Message)
        case "ArgumentError":
            fmt.Println("Argument error:", builderErr.Message)
        default:
            fmt.Println("Unknown error:", builderErr.Message)
        }
    }
    return "", nil, err
}
```

### Check for Specific Error

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    if err == sb.ErrEmptyTableName {
        fmt.Println("Table name is required")
        return "", nil, err
    }
    return "", nil, err
}
```

## Error Handling Patterns

### Pattern 1: Immediate Error Return

```go
func GetUser(db *sql.DB, email string) (*User, error) {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{
            Column:   "email",
            Operator: "=",
            Value:    email,
        })
    
    sql, params, err := builder.Select([]string{"id", "name", "email"})
    if err != nil {
        return nil, fmt.Errorf("failed to build query: %w", err)
    }
    
    rows, err := db.Query(sql, params...)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %w", err)
    }
    defer rows.Close()
    
    // ... process rows
}
```

### Pattern 2: Error Wrapping

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    return "", nil, fmt.Errorf("select query failed: %w", err)
}
```

### Pattern 3: Error Logging

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    log.Printf("Query build error: %v", err)
    return "", nil, err
}
```

### Pattern 4: Error Recovery

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    // Attempt fallback
    fallbackSQL := "SELECT * FROM users"
    fallbackParams := []interface{}{}
    return fallbackSQL, fallbackParams, nil
}
```

## Common Error Scenarios

### Scenario 1: Empty Table Name

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table(""). // Empty table name
    Select([]string{"*"})

sql, params, err := builder.Select([]string{"*"})
// Error: ValidationError: table name cannot be empty
```

### Scenario 2: Empty Column Name

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Column(sb.Column{Name: "", Type: sb.COLUMN_TYPE_STRING})

sql, err := builder.Create()
// Error: ValidationError: column name cannot be empty
```

### Scenario 3: Nil Subquery

```go
builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(nil) // Nil subquery
// Error: ArgumentError: subquery cannot be nil
```

### Scenario 4: Empty ON Condition

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "") // Empty ON condition

sql, params, err := builder.Select([]string{"*"})
// Error: ValidationError: ON condition cannot be empty
```

## Error Handling Best Practices

1. **Always Check Errors** - Never ignore returned errors
2. **Use Error Wrapping** - Provide context with `fmt.Errorf` and `%w`
3. **Handle Specific Errors** - Check for specific error types when appropriate
4. **Log Errors** - Log errors for debugging and monitoring
5. **Provide Context** - Add context to errors for better debugging
6. **Don't Panic** - Handle errors gracefully, never panic in production code

## Error Handling in Transactions

```go
err := db.ExecInTransaction(func(tx *sb.Database) error {
    // First operation
    sql1, params1, err1 := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Insert(map[string]interface{}{
            "name":  "John",
            "email": "john@example.com",
        })
    if err1 != nil {
        return fmt.Errorf("failed to build insert query: %w", err1)
    }
    
    _, err := tx.Exec(sql1, params1...)
    if err != nil {
        return fmt.Errorf("failed to insert user: %w", err)
    }
    
    // Second operation
    sql2, params2, err2 := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("orders").
        Insert(map[string]interface{}{
            "user_id": 1,
            "total":   100,
        })
    if err2 != nil {
        return fmt.Errorf("failed to build order insert query: %w", err2)
    }
    
    _, err = tx.Exec(sql2, params2...)
    if err != nil {
        return fmt.Errorf("failed to insert order: %w", err)
    }
    
    return nil
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

## Testing Error Handling

### Test for Expected Error

```go
func TestEmptyTableName(t *testing.T) {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("").
        Select([]string{"*"})
    
    _, _, err := builder.Select([]string{"*"})
    
    if err == nil {
        t.Fatal("Expected error but got none")
    }
    
    if err != sb.ErrEmptyTableName {
        t.Fatalf("Expected ErrEmptyTableName but got: %v", err)
    }
}
```

### Test Error Type

```go
func TestErrorType(t *testing.T) {
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("").
        Select([]string{"*"})
    
    _, _, err := builder.Select([]string{"*"})
    
    if err == nil {
        t.Fatal("Expected error but got none")
    }
    
    builderErr, ok := err.(*sb.BuilderError)
    if !ok {
        t.Fatal("Expected BuilderError type")
    }
    
    if builderErr.Type != "ValidationError" {
        t.Fatalf("Expected ValidationError but got: %s", builderErr.Type)
    }
}
```

## See Also

- [API Reference](../api_reference.md)
- [Parameterized Queries](parameterized-queries.md)
- [Architecture](../architecture.md)
