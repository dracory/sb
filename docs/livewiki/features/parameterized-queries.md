---
path: features/parameterized-queries.md
page-type: tutorial
summary: Comprehensive guide to parameterized queries in SB SQL Builder, the default secure query method.
tags: [parameterized-queries, security, sql-injection-prevention]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# Parameterized Queries

SB uses **parameterized queries by default** to provide SQL injection protection while maintaining clean, readable code.

## Overview

Parameterized queries separate SQL logic from data, preventing SQL injection attacks and improving database performance through query plan caching.

### Why Parameterized Queries?

- **Security**: Prevents SQL injection attacks
- **Performance**: Better query plan caching
- **Correctness**: Proper handling of special characters
- **Best Practice**: Industry standard for database operations

## Basic Usage

### Simple SELECT with Parameters

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "email",
        Operator: "=",
        Value:    "user@example.com",
    }).
    Select([]string{"id", "name", "email"})
```

**Result:**
```sql
-- SQL
SELECT `id`, `name`, `email` FROM `users` WHERE `email` = ?;

-- Parameters
["user@example.com"]
```

### INSERT with Parameters

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Insert(map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
        "status": "active",
    })
```

**Result:**
```sql
-- SQL
INSERT INTO `users` (`name`, `email`, `status`) VALUES (?, ?, ?);

-- Parameters
["John Doe", "john@example.com", "active"]
```

### UPDATE with Parameters

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "id",
        Operator: "=",
        Value:    1,
    }).
    Update(map[string]interface{}{
        "status": "inactive",
    })
```

**Result:**
```sql
-- SQL
UPDATE `users` SET `status` = ? WHERE `id` = ?;

-- Parameters
["inactive", 1]
```

### DELETE with Parameters

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "id",
        Operator: "=",
        Value:    1,
    }).
    Delete()
```

**Result:**
```sql
-- SQL
DELETE FROM `users` WHERE `id` = ?;

-- Parameters
[1]
```

## Database-Specific Placeholders

SB automatically generates the correct placeholder syntax for each database:

| Database | Placeholder | Example |
|----------|------------|---------|
| MySQL | `?` | `WHERE id = ?` |
| PostgreSQL | `$1, $2, $3` | `WHERE id = $1` |
| SQLite | `?` | `WHERE id = ?` |
| MSSQL | `@p1, @p2, @p3` | `WHERE id = @p1` |

### PostgreSQL Example

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    Where(&sb.Where{
        Column:   "email",
        Operator: "=",
        Value:    "user@example.com",
    }).
    Select([]string{"*"})
```

**Result:**
```sql
-- SQL
SELECT * FROM "users" WHERE "email" = $1;

-- Parameters
["user@example.com"]
```

### MSSQL Example

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MSSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "email",
        Operator: "=",
        Value:    "user@example.com",
    }).
    Select([]string{"*"})
```

**Result:**
```sql
-- SQL
SELECT * FROM [users] WHERE [email] = @p1;

-- Parameters
["user@example.com"]
```

## Parameter Handling

### NULL Values

NULL values are handled correctly - no placeholder is generated:

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "deleted_at",
        Operator: "=",
        Value:    "", // Empty string treated as NULL
    }).
    Select([]string{"*"})
```

**Result:**
```sql
-- SQL
SELECT * FROM `users` WHERE `deleted_at` IS NULL;

-- Parameters
[]
```

### Empty Strings

Empty strings are treated as empty strings, not NULL:

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Insert(map[string]interface{}{
        "name":  "",
        "email": "test@example.com",
    })
```

**Result:**
```sql
-- SQL
INSERT INTO `users` (`name`, `email`) VALUES (?, ?);

-- Parameters
["", "test@example.com"]
```

### Parameter Ordering

Parameters are ordered consistently:
1. WHERE clause parameters (first)
2. INSERT/UPDATE values (second)
3. Subquery parameters (nested)

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    }).
    Where(&sb.Where{
        Column:   "role",
        Operator: "=",
        Value:    "admin",
    }).
    Update(map[string]interface{}{
        "last_login": "2026-05-08",
    })
```

**Result:**
```sql
-- SQL
UPDATE `users` SET `last_login` = ? WHERE `status` = ? AND `role` = ?;

-- Parameters
["2026-05-08", "active", "admin"]
```

## Subquery Parameters

Subqueries also use parameterized queries:

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "total",
        Operator: ">",
        Value:    1000,
    })

sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery).
    Select([]string{"name"})
```

**Result:**
```sql
-- SQL
SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > ?);

-- Parameters
[1000]
```

## Executing Parameterized Queries

### With database/sql

```go
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

db, err := sql.Open("mysql", "user:pass@/dbname")
if err != nil {
    panic(err)
}

// Build query
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "email",
        Operator: "=",
        Value:    "user@example.com",
    }).
    Select([]string{"id", "name", "email"})
if err != nil {
    panic(err)
}

// Execute with parameters
rows, err := db.Query(sql, params...)
if err != nil {
    panic(err)
}
defer rows.Close()

for rows.Next() {
    var id int
    var name, email string
    err := rows.Scan(&id, &name, &email)
    if err != nil {
        panic(err)
    }
    fmt.Println(id, name, email)
}
```

### With SB Database Wrapper

```go
db := sb.NewDatabaseFromDriver("mysql", "user:pass@/dbname")

// Build query
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "email",
        Operator: "=",
        Value:    "user@example.com",
    }).
    Select([]string{"id", "name", "email"})
if err != nil {
    panic(err)
}

// Execute with parameters
results, err := db.Query(sql, params...)
if err != nil {
    panic(err)
}
```

## Legacy Mode (Interpolated Values)

For backward compatibility or specific use cases, you can disable parameterization:

```go
sql, _, err := builder.
    WithInterpolatedValues().
    Select([]string{"*"})
```

**Warning:** This mode is vulnerable to SQL injection and should be used with caution.

## Migration Guide

### From v0.17.x to v0.18.0

**Before (v0.17.x):**
```go
sql, err := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userInput}).
    Select([]string{"*"})
```

**After (v0.18.0):**
```go
sql, params, err := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userInput}).
    Select([]string{"*"})
```

Simply handle the additional return value:

```go
sql, params, err := builder.Select([]string{"*"})
if err != nil {
    return err
}
// Use sql and params...
```

## Security Benefits

### SQL Injection Prevention

**Vulnerable (without parameterization):**
```go
userInput := "' OR '1'='1"
sql := "SELECT * FROM users WHERE email = '" + userInput + "'"
// Results in: SELECT * FROM users WHERE email = '' OR '1'='1'
// This returns all users!
```

**Secure (with parameterization):**
```go
userInput := "' OR '1'='1"
sql, params := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: userInput}).
    Select([]string{"*"})
// Results in: SELECT * FROM users WHERE email = ?
// Params: ["' OR '1'='1"]
// This searches for the literal string
```

## Performance Benefits

### Query Plan Caching

Parameterized queries allow databases to cache query plans:

```go
// First execution
db.Query("SELECT * FROM users WHERE id = ?", 1)
// Database parses and caches plan

// Second execution
db.Query("SELECT * FROM users WHERE id = ?", 2)
// Database reuses cached plan (faster)
```

### Prepared Statement Optimization

```go
stmt, err := db.Prepare("SELECT * FROM users WHERE id = ?")
if err != nil {
    panic(err)
}

// Execute multiple times with different parameters
stmt.Query(1)
stmt.Query(2)
stmt.Query(3)
```

## Best Practices

1. **Always use parameterized queries** - They're the default in SB for security
2. **Handle parameters correctly** - Pass them as a variadic to database methods
3. **Validate input separately** - Parameterization doesn't replace input validation
4. **Use appropriate types** - Let the database driver handle type conversion
5. **Test with real data** - Ensure parameters work with your actual data

## Common Issues

### Issue: Too Many Parameters

**Problem:** Some databases have limits on parameter counts.

**Solution:** Use batch operations or split queries.

### Issue: Parameter Type Mismatch

**Problem:** Database rejects parameter due to type mismatch.

**Solution:** Ensure Go types match database column types.

### Issue: NULL Handling

**Problem:** Confusion between empty strings and NULL.

**Solution:** Use explicit NULL values or database-specific NULL handling.

## See Also

- [API Reference](../api_reference.md)
- [Error Handling](error-handling.md)
- [Security](../security.md)
- [Migration Guide](../migration.md)
