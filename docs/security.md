# Security Guide

## ⚠️ Important Notice

SB is a SQL builder library, not an ORM. It generates SQL strings by quoting values. While this provides basic protection against SQL injection, it's not as secure as parameterized queries.

## SQL Injection Prevention

### ✅ Safe: Using the Builder API

The builder API properly quotes values to prevent SQL injection:

```go
// Safe - values are properly quoted
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail}).
    Select([]string{"*"})

// Result: SELECT * FROM `users` WHERE `email` = 'user@example.com'
```

### ⚠️ Caution: Raw SQL

Using the Raw field allows arbitrary SQL and can be dangerous:

```go
// Potentially unsafe - user input in Raw field
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Raw: "email = '" + userEmail + "'"}).
    Select([]string{"*"})
```

### ❌ Dangerous: String Concatenation

Never concatenate user input directly into SQL strings:

```go
// NEVER do this - vulnerable to injection
sql := "SELECT * FROM users WHERE email = '" + userEmail + "'"
```

## Security Best Practices

### 1. Always Use the Builder API

```go
// ✅ Recommended
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "name", Operator: "LIKE", Value: "%" + search + "%"}).
    Select([]string{"*"})
```

### 2. Validate Input Before Using

```go
// Validate email format before using in query
if !isValidEmail(userEmail) {
    return errors.New("invalid email format")
}

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail}).
    Select([]string{"*"})
```

### 3. Use Parameterized Queries When Available

SB will support parameterized queries in v0.18.0:

```go
// Future (v0.18.0)
sql, params := builder.
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: userEmail}).
    SelectWithParams([]string{"*"})

// Execute safely
db.Exec(sql, params...)
```

### 4. Avoid Raw SQL with User Input

```go
// ❌ Dangerous
rawCondition := fmt.Sprintf("name = '%s'", userName)
sql := builder.Where(&sb.Where{Raw: rawCondition}).Select([]string{"*"})

// ✅ Safe
sql := builder.
    Where(&sb.Where{Column: "name", Operator: "=", Value: userName}).
    Select([]string{"*"})
```

## Database-Specific Considerations

### MySQL

MySQL uses backtick quoting for identifiers:

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: "user@example.com"}).
    Select([]string{"*"})

// Result: SELECT * FROM `users` WHERE `email` = 'user@example.com'
```

### PostgreSQL

PostgreSQL uses double quote quoting:

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: "user@example.com"}).
    Select([]string{"*"})

// Result: SELECT * FROM "users" WHERE "email" = 'user@example.com'
```

### SQLite

SQLite uses double quote quoting for identifiers:

```go
sql := sb.NewBuilder(sb.DIALECT_SQLITE).
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: "user@example.com"}).
    Select([]string{"*"})

// Result: SELECT * FROM "users" WHERE "email" = 'user@example.com'
```

### MSSQL

MSSQL uses square bracket quoting:

```go
sql := sb.NewBuilder(sb.DIALECT_MSSQL).
    Table("users").
    Where(&sb.Where{Column: "email", Operator: "=", Value: "user@example.com"}).
    Select([]string{"*"})

// Result: SELECT * FROM [users] WHERE [email] = 'user@example.com'
```

## Common Attack Vectors

### 1. Injection Through Column Names

```go
// ❌ Dangerous - column name from user input
columnName := req.FormValue("column")
sql := builder.Where(&sb.Where{Column: columnName, Operator: "=", Value: "value"})

// ✅ Safe - validate column name against whitelist
validColumns := map[string]bool{"name": true, "email": true, "id": true}
columnName := req.FormValue("column")
if !validColumns[columnName] {
    return errors.New("invalid column name")
}
sql := builder.Where(&sb.Where{Column: columnName, Operator: "=", Value: "value"})
```

### 2. Injection Through Table Names

```go
// ❌ Dangerous - table name from user input
tableName := req.FormValue("table")
sql := builder.Table(tableName).Select([]string{"*"})

// ✅ Safe - validate table name
validTables := map[string]bool{"users": true, "orders": true, "products": true}
tableName := req.FormValue("table")
if !validTables[tableName] {
    return errors.New("invalid table name")
}
sql := builder.Table(tableName).Select([]string{"*"})
```

### 3. Injection Through Raw SQL

```go
// ❌ Dangerous - raw SQL with user input
rawSQL := fmt.Sprintf("name LIKE '%%%s%%'", searchTerm)
sql := builder.Where(&sb.Where{Raw: rawSQL}).Select([]string{"*"})

// ✅ Safe - use builder API
sql := builder.
    Where(&sb.Where{Column: "name", Operator: "LIKE", Value: "%" + searchTerm + "%"}).
    Select([]string{"*"})
```

## Migration to Parameterized Queries

When parameterized queries are available (v0.18.0), migrate existing code:

### Current (v0.16.x)
```go
sql := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: email}).
    Select([]string{"*"})

// Execute
rows, err := db.Query(sql)
```

### Future (v0.18.0)
```go
sql, params := builder.
    Where(&sb.Where{Column: "email", Operator: "=", Value: email}).
    SelectWithParams([]string{"*"})

// Execute safely
rows, err := db.Query(sql, params...)
```

## Security Checklist

- [ ] Always validate user input before using in queries
- [ ] Use the Builder API instead of string concatenation
- [ ] Avoid using the Raw field with user input
- [ ] Validate column and table names against whitelists
- [ ] Use parameterized queries when available (v0.18.0+)
- [ ] Review generated SQL for unexpected patterns
- [ ] Keep SB updated to the latest version

## Reporting Security Issues

If you discover a security vulnerability in SB:

1. **Do not** create a public issue
2. Email security details to: security@dracory.com
3. Include details about the vulnerability
4. Provide a proof-of-concept if possible
5. Allow time for the issue to be fixed before disclosure

## Additional Resources

- [OWASP SQL Injection Prevention](https://owasp.org/www-community/attacks/SQL_Injection)
- [Go Database Security Best Practices](https://go.dev/doc/database/sql-injection)
- [Parameterized Queries Guide](https://use-the-index-luke.com/sql/injection/)

---

**Last Updated:** March 20, 2026  
**Version:** 0.16.x  
**Next Update:** v0.18.0 (parameterized queries)
