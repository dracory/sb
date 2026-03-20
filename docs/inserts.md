# INSERT Operations

## Basic INSERT

Insert new records into tables with column-value mappings.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("cache").
	Insert(map[string]string{
		"ID":         uid.NanoUid(),
		"CacheKey":   token,
		"CacheValue": string(emailJSON),
		"ExpiresAt":  expiresAt.Format("2006-01-02T15:04:05"),
		"CreatedAt":  time.Now().Format("2006-01-02T15:04:05"),
		"UpdatedAt":  time.Now().Format("2006-01-02T15:04:05"),
	})
```

## INSERT with Specific Columns

Insert data with explicit column specification.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Insert(map[string]string{
		"name":       "John Doe",
		"email":      "john@example.com",
		"created_at": time.Now().Format("2006-01-02T15:04:05"),
	})
```

## Database-Specific Behavior

### MySQL
```sql
INSERT INTO `table_name` (`column1`, `column2`) VALUES ('value1', 'value2');
```

### PostgreSQL
```sql
INSERT INTO "table_name" ("column1", "column2") VALUES ('value1', 'value2');
```

### SQLite
```sql
INSERT INTO "table_name" ("column1", "column2") VALUES ('value1', 'value2');
```

### MSSQL
```sql
INSERT INTO [table_name] ([column1], [column2]) VALUES ('value1', 'value2');
```

## Data Types

When inserting data, ensure proper formatting for different column types:

### String/Varchar
```go
"name": "John Doe"
```

### Integer
```go
"age": "25"
```

### Decimal/Numeric
```go
"price": "99.99"
```

### DateTime
```go
"created_at": time.Now().Format("2006-01-02T15:04:05")
```

### Boolean
```go
"active": "1"  // or "true" depending on database
```

## Best Practices

1. **Always specify table name** before calling Insert()
2. **Use proper data formatting** for timestamps and numbers
3. **Include created_at/updated_at** timestamps for audit trails
4. **Handle nil values** by omitting them from the map or using empty strings
5. **Use prepared statements** for actual execution to prevent SQL injection

## Example: User Registration

```go
func RegisterUser(db *sb.Database, name, email string) error {
    ctx := context.Background()
    
    sql := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Insert(map[string]string{
            "name":       name,
            "email":      email,
            "status":     "active",
            "created_at": time.Now().Format("2006-01-02T15:04:05"),
            "updated_at": time.Now().Format("2006-01-02T15:04:05"),
        })
    
    _, err := db.Exec(ctx, sql)
    return err
}
```

## Example: Bulk Insert Pattern

While SB doesn't have a built-in bulk insert method, you can create a pattern:

```go
func BulkInsertUsers(db *sb.Database, users []User) error {
    ctx := context.Background()
    
    for _, user := range users {
        sql := sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("users").
            Insert(map[string]string{
                "name":       user.Name,
                "email":      user.Email,
                "created_at": user.CreatedAt.Format("2006-01-02T15:04:05"),
            })
        
        if _, err := db.Exec(ctx, sql); err != nil {
            return err
        }
    }
    
    return nil
}
```
