# UPDATE and DELETE Operations

## UPDATE Statements

Modify existing records in tables.

### Basic UPDATE

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{
		Column: "id",
		Operator: "=",
		Value: "1",
	}).
	Update(map[string]string{
		"name":       "John Doe",
		"email":      "john@example.com",
		"updated_at": time.Now().Format("2006-01-02T15:04:05"),
	})
```

### UPDATE with Multiple Conditions

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
	Where(&sb.Where{Column: "last_login", Operator: "<", Value: "2023-01-01"}).
	Update(map[string]string{
		"status":     "inactive",
		"updated_at": time.Now().Format("2006-01-02T15:04:05"),
	})
```

### UPDATE with JOIN

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	Where(&sb.Where{Column: "users.id", Operator: "=", Value: "1"}).
	Update(map[string]string{"orders.status": "processed"})
```

## DELETE Statements

Remove records from tables.

### Basic DELETE

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{
		Column: "id",
		Operator: "=",
		Value: "1",
	}).
	Delete()
```

### DELETE with Limit

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{
		Column: "status",
		Operator: "=",
		Value: "inactive",
	}).
	Limit(10).
	Delete()
```

### DELETE with Multiple Conditions

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"}).
	Where(&sb.Where{Column: "last_login", Operator: "<", Value: "2023-01-01"}).
	Delete()
```

## Database-Specific Behavior

### MySQL
```sql
-- UPDATE
UPDATE `table_name` SET `column1` = 'value1' WHERE `id` = '1';

-- DELETE
DELETE FROM `table_name` WHERE `id` = '1';
```

### PostgreSQL
```sql
-- UPDATE
UPDATE "table_name" SET "column1" = 'value1' WHERE "id" = '1';

-- DELETE
DELETE FROM "table_name" WHERE "id" = '1';
```

### SQLite
```sql
-- UPDATE
UPDATE "table_name" SET "column1" = 'value1' WHERE "id" = '1';

-- DELETE
DELETE FROM "table_name" WHERE "id" = '1';
```

### MSSQL
```sql
-- UPDATE
UPDATE [table_name] SET [column1] = 'value1' WHERE [id] = '1';

-- DELETE
DELETE FROM [table_name] WHERE [id] = '1';
```

## Best Practices

### UPDATE Best Practices

1. **Always use WHERE clause** to avoid updating all records
2. **Include updated_at timestamp** for audit trails
3. **Use transactions** for critical updates
4. **Validate data** before updating

```go
func UpdateUserStatus(db *sb.Database, userID int, status string) error {
    ctx := context.Background()
    
    sql := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(userID)}).
        Update(map[string]string{
            "status":     status,
            "updated_at": time.Now().Format("2006-01-02T15:04:05"),
        })
    
    _, err := db.Exec(ctx, sql)
    return err
}
```

### DELETE Best Practices

1. **Use soft deletes** when possible (update status instead of deleting)
2. **Always use WHERE clause** to avoid deleting all records
3. **Use transactions** for cascading deletes
4. **Consider foreign key constraints**

```go
func SoftDeleteUser(db *sb.Database, userID int) error {
    ctx := context.Background()
    
    sql := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "id", Operator: "=", Value: strconv.Itoa(userID)}).
        Update(map[string]string{
            "status":     "deleted",
            "deleted_at": time.Now().Format("2006-01-02T15:04:05"),
        })
    
    _, err := db.Exec(ctx, sql)
    return err
}
```

## Transaction Examples

### UPDATE in Transaction

```go
err := db.ExecInTransaction(ctx, func(tx *sb.Database) error {
    // Update user
    userSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "id", Operator: "=", Value: userID}).
        Update(map[string]string{
            "status":     "active",
            "updated_at": time.Now().Format("2006-01-02T15:04:05"),
        })
    
    if _, err := tx.Exec(ctx, userSQL); err != nil {
        return err
    }
    
    // Update related records
    profileSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("profiles").
        Where(&sb.Where{Column: "user_id", Operator: "=", Value: userID}).
        Update(map[string]string{
            "verified":   "1",
            "updated_at": time.Now().Format("2006-01-02T15:04:05"),
        })
    
    if _, err := tx.Exec(ctx, profileSQL); err != nil {
        return err
    }
    
    return nil
})
```

### DELETE in Transaction

```go
err := db.ExecInTransaction(ctx, func(tx *sb.Database) error {
    // Delete user's orders first (foreign key constraint)
    ordersSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("orders").
        Where(&sb.Where{Column: "user_id", Operator: "=", Value: userID}).
        Delete()
    
    if _, err := tx.Exec(ctx, ordersSQL); err != nil {
        return err
    }
    
    // Delete user profile
    profileSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("profiles").
        Where(&sb.Where{Column: "user_id", Operator: "=", Value: userID}).
        Delete()
    
    if _, err := tx.Exec(ctx, profileSQL); err != nil {
        return err
    }
    
    // Delete user
    userSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Where(&sb.Where{Column: "id", Operator: "=", Value: userID}).
        Delete()
    
    if _, err := tx.Exec(ctx, userSQL); err != nil {
        return err
    }
    
    return nil
})
```
