# Index Management

## CREATE INDEX

Create indexes to improve query performance on frequently searched columns.

### Basic Index Creation

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_email", "email")
// Result: CREATE INDEX `idx_users_email` ON `users` (`email`);
```

### Composite Index

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_user_status", "user_id", "status")
// Result: CREATE INDEX `idx_orders_user_status` ON `orders` (`user_id`, `status`);
```

### Unique Index

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	CreateUniqueIndex("uq_users_email", "email")
// Result: CREATE UNIQUE INDEX "uq_users_email" ON "users" ("email");
```

## DROP INDEX

Remove indexes from tables.

### Basic Index Drop

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	DropIndex("idx_users_email")
// Result: DROP INDEX `idx_users_email` ON `users`;
```

### Safe Index Drop (IF EXISTS)

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	DropIndexIfExists("idx_users_email")
// Result: DROP INDEX IF EXISTS "idx_users_email";
```

### Schema-Aware Index Drop (PostgreSQL)

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	DropIndexWithSchema("idx_users_email", "public")
// Result: DROP INDEX IF EXISTS "public"."idx_users_email";
```

## Database-Specific Behavior

### MySQL

```sql
-- Create Index
CREATE INDEX `idx_name` ON `table_name` (`column1`, `column2`);

-- Drop Index
DROP INDEX `idx_name` ON `table_name`;

-- Note: MySQL doesn't support IF EXISTS for indexes
```

### PostgreSQL

```sql
-- Create Index
CREATE INDEX "idx_name" ON "table_name" ("column1", "column2");

-- Create Unique Index
CREATE UNIQUE INDEX "uq_name" ON "table_name" ("column");

-- Drop Index
DROP INDEX "idx_name";

-- Safe Drop
DROP INDEX IF EXISTS "idx_name";

-- Schema-Aware Drop
DROP INDEX IF EXISTS "schema"."idx_name";
```

### SQLite

```sql
-- Create Index
CREATE INDEX "idx_name" ON "table_name" ("column1", "column2");

-- Create Unique Index
CREATE UNIQUE INDEX "uq_name" ON "table_name" ("column");

-- Drop Index
DROP INDEX "idx_name";

-- Safe Drop
DROP INDEX IF EXISTS "idx_name";
```

### MSSQL

```sql
-- Create Index
CREATE INDEX [idx_name] ON [table_name] ([column1], [column2]);

-- Create Unique Index
CREATE UNIQUE INDEX [uq_name] ON [table_name] ([column]);

-- Drop Index
DROP INDEX [idx_name] ON [table_name];

-- Safe Drop (SQL Server 2016+)
DROP INDEX IF EXISTS [idx_name] ON [table_name];
```

## Index Methods Reference

### Creation Methods

- `CreateIndex(indexName, columns...)` - Create basic index
- `CreateUniqueIndex(indexName, columns...)` - Create unique index

### Drop Methods

- `DropIndex(indexName)` - Drop index
- `DropIndexIfExists(indexName)` - Drop index with IF EXISTS
- `DropIndexWithSchema(indexName, schema)` - Drop index with schema (PostgreSQL)

## Best Practices

### When to Create Indexes

```go
// Good: Index frequently searched columns
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_email", "email")  // Login queries

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_user_date", "user_id", "created_at")  // User order history

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("products").
	CreateIndex("idx_products_category_price", "category", "price")  // Category browsing
```

### Composite Index Strategy

```go
// Most selective column first
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_status_user_date", "status", "user_id", "created_at")

// Covers queries like:
// - WHERE status = 'active'
// - WHERE status = 'active' AND user_id = 123
// - WHERE status = 'active' AND user_id = 123 ORDER BY created_at
```

### Unique Constraints

```go
// Enforce data uniqueness
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	CreateUniqueIndex("uq_users_email", "email")

sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	CreateUniqueIndex("uq_users_username", "username")

sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("products").
	CreateUniqueIndex("uq_products_sku", "sku")
```

## Common Index Patterns

### User Table Indexes

```go
// Primary email lookup (login)
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_email", "email")

// Username search
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_username", "username")

// Status filtering
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_status", "status")

// Registration date queries
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_created_at", "created_at")
```

### Order Table Indexes

```go
// User order history
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_user_id", "user_id")

// Status filtering
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_status", "status")

// Date range queries
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_created_at", "created_at")

// Complex query support
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	CreateIndex("idx_orders_status_user_date", "status", "user_id", "created_at")
```

### Product Table Indexes

```go
// Category browsing
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("products").
	CreateIndex("idx_products_category", "category")

// Price range queries
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("products").
	CreateIndex("idx_products_price", "price")

// Search functionality
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("products").
	CreateIndex("idx_products_name", "name")

// Category + price sorting
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("products").
	CreateIndex("idx_products_category_price", "category", "price")
```

## Index Management Examples

### Complete Index Lifecycle

```go
func SetupUserIndexes(db *sb.Database) error {
    ctx := context.Background()
    
    // Create indexes
    indexes := []string{
        sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("users").
            CreateIndex("idx_users_email", "email"),
            
        sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("users").
            CreateIndex("idx_users_status", "status"),
            
        sb.NewBuilder(sb.DIALECT_MYSQL).
            Table("users").
            CreateUniqueIndex("uq_users_username", "username"),
    }
    
    for _, sql := range indexes {
        if _, err := db.Exec(ctx, sql); err != nil {
            return err
        }
    }
    
    return nil
}
```

### Safe Index Recreation

```go
func RecreateIndex(db *sb.Database, tableName, indexName string, columns ...string) error {
    ctx := context.Background()
    
    // Drop if exists
    dropSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
        Table(tableName).
        DropIndexIfExists(indexName)
    
    db.Exec(ctx, dropSQL) // Ignore error if index doesn't exist
    
    // Create new index
    createSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
        Table(tableName).
        CreateIndex(indexName, columns...)
    
    _, err := db.Exec(ctx, createSQL)
    return err
}
```

## Performance Considerations

1. **Don't over-index** - Each index slows down writes
2. **Choose columns wisely** - High cardinality columns make better indexes
3. **Monitor query performance** - Use EXPLAIN to verify index usage
4. **Consider index size** - Large indexes consume memory and disk space
5. **Regular maintenance** - Rebuild fragmented indexes periodically

## Index Naming Conventions

- `idx_table_column` - Basic index
- `idx_table_columns` - Composite index (column names separated by underscore)
- `uq_table_column` - Unique index
- `fk_table_column` - Foreign key index
- `ix_table_purpose` - Purpose-specific index
