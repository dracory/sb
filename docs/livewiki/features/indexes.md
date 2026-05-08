---
path: features/indexes.md
page-type: tutorial
summary: Complete guide to index management in SB SQL Builder including advanced index types and database-specific features.
tags: [indexes, create-index, drop-index, unique-index, partial-index]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# Index Management

SB provides comprehensive index management including basic indexes, unique indexes, composite indexes, partial indexes, and database-specific features like FULLTEXT and GIN indexes.

## Overview

Indexes are essential for database performance, allowing fast data retrieval. SB supports creating and dropping indexes across all supported database dialects with advanced features.

### Supported Index Types

- **Basic Index**: Single or multi-column indexes
- **Unique Index**: Enforces uniqueness on indexed columns
- **Composite Index**: Multi-column indexes with explicit ordering
- **Partial Index**: Indexes with WHERE predicates (PostgreSQL, SQLite, MSSQL)
- **Covering Index**: Indexes with INCLUDE clause (PostgreSQL, MSSQL)
- **Database-Specific**: FULLTEXT (MySQL), GIN/GIST (PostgreSQL), etc.

## Basic Index Operations

### Create Basic Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    CreateIndex("idx_users_email", "email")
```

**MySQL Result:**
```sql
CREATE INDEX `idx_users_email` ON `users` (`email`);
```

**PostgreSQL Result:**
```sql
CREATE INDEX "idx_users_email" ON "users" ("email");
```

### Create Unique Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    CreateUniqueIndex("idx_users_email_unique", "email")
```

**MySQL Result:**
```sql
CREATE UNIQUE INDEX `idx_users_email_unique` ON `users` (`email`);
```

### Create Composite Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateCompositeIndex("idx_orders_user_status", []sb.IndexColumn{
        {Name: "user_id", Direction: "ASC"},
        {Name: "status", Direction: "ASC"},
    })
```

**MySQL Result:**
```sql
CREATE INDEX `idx_orders_user_status` ON `users` (`user_id` ASC, `status` ASC);
```

## Advanced Index Operations

### Create Partial Index (PostgreSQL, SQLite, MSSQL)

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    CreatePartialIndex("idx_users_active_email", "deleted_at IS NULL", "email")
```

**PostgreSQL Result:**
```sql
CREATE INDEX "idx_users_active_email" ON "users" ("email") WHERE deleted_at IS NULL;
```

### Create Covering Index (PostgreSQL, MSSQL)

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    CreateCoveringIndex("idx_orders_user_covering", []string{"order_date", "total"}, "user_id")
```

**PostgreSQL Result:**
```sql
CREATE INDEX "idx_orders_user_covering" ON "orders" ("user_id") INCLUDE ("order_date", "total");
```

### Create Index with Options

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("documents").
    CreateIndexWithOptions("idx_documents_search", sb.IndexOptions{
        Unique:      true,
        IfNotExists: true,
        Using:       sb.INDEX_TYPE_GIN,
        Columns:     []sb.IndexColumn{{Name: "search_vector"}},
        Include:     []string{"title", "content"},
        Where:       "published = true",
        Storage:     "fillfactor=90",
    })
```

**PostgreSQL Result:**
```sql
CREATE UNIQUE INDEX IF NOT EXISTS "idx_documents_search" ON "documents" USING GIN ("search_vector") INCLUDE ("title", "content") WHERE published = true WITH (fillfactor=90);
```

## Database-Specific Index Features

### MySQL FULLTEXT Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("products").
    CreateIndexWithOptions("idx_products_search", sb.IndexOptions{
        Using:   sb.INDEX_TYPE_FULLTEXT,
        Columns: []sb.IndexColumn{
            {Name: "name", PrefixLength: 100},
            {Name: "description", PrefixLength: 255},
        },
        Comment: "Full-text search index for products",
    })
```

**MySQL Result:**
```sql
CREATE FULLTEXT INDEX `idx_products_search` ON `products` (`name`(100), `description`(255)) COMMENT 'Full-text search index for products';
```

### MySQL SPATIAL Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("locations").
    CreateIndexWithOptions("idx_locations_geo", sb.IndexOptions{
        Using:   sb.INDEX_TYPE_SPATIAL,
        Columns: []sb.IndexColumn{{Name: "coordinates"}},
    })
```

**MySQL Result:**
```sql
CREATE SPATIAL INDEX `idx_locations_geo` ON `locations` (`coordinates`);
```

### PostgreSQL GIN Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("documents").
    CreateIndexWithOptions("idx_documents_content", sb.IndexOptions{
        Unique:  true,
        Using:   sb.INDEX_TYPE_GIN,
        Columns: []sb.IndexColumn{{Name: "content_vector"}},
        Storage: "fastupdate=on",
    })
```

**PostgreSQL Result:**
```sql
CREATE UNIQUE INDEX "idx_documents_content" ON "documents" USING GIN ("content_vector") WITH (fastupdate=on);
```

### PostgreSQL GIST Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("locations").
    CreateIndexWithOptions("idx_locations_geo", sb.IndexOptions{
        Using:   sb.INDEX_TYPE_GIST,
        Columns: []sb.IndexColumn{{Name: "coordinates"}},
    })
```

**PostgreSQL Result:**
```sql
CREATE INDEX "idx_locations_geo" ON "locations" USING GIST ("coordinates");
```

## Drop Index Operations

### Drop Basic Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    DropIndex("idx_users_email")
```

**MySQL Result:**
```sql
DROP INDEX `idx_users_email` ON `users`;
```

**PostgreSQL Result:**
```sql
DROP INDEX "idx_users_email";
```

### Drop Index IF EXISTS

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    DropIndexIfExists("idx_users_email")
```

**PostgreSQL Result:**
```sql
DROP INDEX IF EXISTS "idx_users_email";
```

### Drop Index with Schema (PostgreSQL)

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
        IfExists: true,
        Schema:   "public",
    })
```

**PostgreSQL Result:**
```sql
DROP INDEX IF EXISTS "public"."idx_users_email";
```

## Index Direction

### ASC Direction (Default)

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateCompositeIndex("idx_orders_date", []sb.IndexColumn{
        {Name: "created_at", Direction: "ASC"},
    })
```

**Result:**
```sql
CREATE INDEX `idx_orders_date` ON `orders` (`created_at` ASC);
```

### DESC Direction

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateCompositeIndex("idx_orders_date_desc", []sb.IndexColumn{
        {Name: "created_at", Direction: "DESC"},
    })
```

**Result:**
```sql
CREATE INDEX `idx_orders_date_desc` ON `orders` (`created_at` DESC);
```

### Mixed Directions

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateCompositeIndex("idx_orders_user_date", []sb.IndexColumn{
        {Name: "user_id", Direction: "ASC"},
        {Name: "created_at", Direction: "DESC"},
    })
```

**Result:**
```sql
CREATE INDEX `idx_orders_user_date` ON `orders` (`user_id` ASC, `created_at` DESC);
```

## MySQL Prefix Length

### Text Column with Prefix

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("articles").
    CreateIndexWithOptions("idx_articles_title", sb.IndexOptions{
        Columns: []sb.IndexColumn{
            {Name: "title", PrefixLength: 50},
        },
    })
```

**MySQL Result:**
```sql
CREATE INDEX `idx_articles_title` ON `articles` (`title`(50));
```

### Multiple Columns with Prefixes

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("products").
    CreateIndexWithOptions("idx_products_search", sb.IndexOptions{
        Columns: []sb.IndexColumn{
            {Name: "name", PrefixLength: 100},
            {Name: "description", PrefixLength: 255},
        },
    })
```

**MySQL Result:**
```sql
CREATE INDEX `idx_products_search` ON `products` (`name`(100), `description`(255));
```

## Error Handling

### Empty Index Name

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users")

sql, err := builder.CreateIndex("", "email")
// Error: ValidationError: index name cannot be empty
```

### Missing Table

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL)
// No table set

sql, err := builder.CreateIndex("idx_users_email", "email")
// Error: ValidationError: table name is required
```

### Empty Columns

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    CreateIndexWithOptions("idx_users", sb.IndexOptions{
        Columns: []sb.IndexColumn{}, // Empty columns
    })
// Error: ValidationError: columns cannot be empty
```

## Index Feature Matrix

| Feature              | MySQL | PostgreSQL | SQLite | MSSQL |
|----------------------|-------|------------|--------|-------|
| UNIQUE index         | ✓     | ✓          | ✓      | ✓     |
| Composite index      | ✓     | ✓          | ✓      | ✓     |
| Column direction     | ✓     | ✓          | ✓      | ✓     |
| IF NOT EXISTS        | ✗     | ✓          | ✓      | ✗     |
| Partial index (WHERE)| ✗     | ✓          | ✓      | ✓     |
| INCLUDE columns      | ✗     | ✓          | ✗      | ✓     |
| USING clause         | ✓*    | ✓          | ✗      | ✗     |
| FULLTEXT index       | ✓     | ✗          | ✗      | ✗     |
| SPATIAL index        | ✓     | ✗          | ✗      | ✗     |
| GIN / GIST / BRIN    | ✗     | ✓          | ✗      | ✗     |
| Prefix length        | ✓     | ✗          | ✗      | ✗     |
| Storage params (WITH)| ✗     | ✓          | ✗      | ✗     |
| Index COMMENT        | ✓     | ✗          | ✗      | ✗     |

*MySQL USING goes inside the column list parentheses, not before it.

## Common Use Cases

### Email Uniqueness

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    CreateUniqueIndex("idx_users_email_unique", "email")
```

### Foreign Key Performance

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateIndex("idx_orders_user_id", "user_id")
```

### Composite Index for Common Queries

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateCompositeIndex("idx_orders_user_status", []sb.IndexColumn{
        {Name: "user_id", Direction: "ASC"},
        {Name: "status", Direction: "ASC"},
    })
```

### Partial Index for Active Users

```go
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    CreatePartialIndex("idx_users_active_email", "status = 'active'", "email")
```

### Full-Text Search Index

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("products").
    CreateIndexWithOptions("idx_products_search", sb.IndexOptions{
        Using: sb.INDEX_TYPE_FULLTEXT,
        Columns: []sb.IndexColumn{
            {Name: "name", PrefixLength: 100},
            {Name: "description", PrefixLength: 255},
        },
    })
```

## Performance Considerations

### When to Create Indexes

- **Frequent WHERE clauses**: Columns used in WHERE conditions
- **JOIN columns**: Foreign keys and join conditions
- **ORDER BY columns**: Columns used for sorting
- **Frequently queried columns**: High-selectivity columns

### When NOT to Create Indexes

- **Small tables**: Full table scan may be faster
- **Frequently updated tables**: Indexes slow down INSERT/UPDATE/DELETE
- **Low-selectivity columns**: Columns with few unique values
- **Write-heavy workloads**: Too many indexes hurt write performance

### Index Maintenance

- **Monitor index usage**: Remove unused indexes
- **Rebuild indexes**: Periodically rebuild fragmented indexes
- **Analyze query plans**: Ensure indexes are being used
- **Consider index size**: Large indexes consume memory and disk space

## Best Practices

1. **Index foreign keys**: Always index foreign key columns
2. **Use unique indexes for unique constraints**: Enforce data integrity
3. **Consider composite indexes**: For multi-column queries
4. **Use partial indexes**: For frequently filtered subsets
5. **Monitor index usage**: Remove unused indexes
6. **Test with real data**: Index effectiveness varies by data distribution
7. **Consider write overhead**: Indexes slow down writes

## See Also

- [API Reference](../api_reference.md)
- [Schema Operations](../schema.md)
- [Performance](../performance.md)
