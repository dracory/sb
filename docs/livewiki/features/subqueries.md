---
path: features/subqueries.md
page-type: tutorial
summary: Complete guide to subquery support in SB SQL Builder including EXISTS, IN, and correlated subqueries.
tags: [subqueries, exists, in-subquery, correlated-subqueries]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# Subqueries

SB provides comprehensive subquery support including EXISTS, IN, NOT IN, and correlated subqueries across all supported database dialects.

## Overview

Subqueries allow you to nest queries within queries, enabling complex data retrieval operations that would otherwise require multiple queries or complex joins.

### Supported Subquery Types

- **EXISTS/NOT EXISTS**: Check for existence of rows
- **IN/NOT IN**: Check if values exist in a set
- **Comparison Operators**: Use subqueries with =, >, <, >=, <=, <>
- **Correlated Subqueries**: Reference outer query columns

## EXISTS Subqueries

### Basic EXISTS

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(subquery)

sql, params, err := builder.Select([]string{"name", "email"})
```

**MySQL Result:**
```sql
SELECT `name`, `email` FROM `users` WHERE EXISTS (SELECT * FROM `orders` WHERE `status` = ?);
-- Params: ["active"]
```

**PostgreSQL Result:**
```sql
SELECT "name", "email" FROM "users" WHERE EXISTS (SELECT * FROM "orders" WHERE "status" = $1);
-- Params: ["active"]
```

### NOT EXISTS

```go
subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "cancelled",
    })

builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    NotExists(subquery)

sql, params, err := builder.Select([]string{"name", "email"})
```

**PostgreSQL Result:**
```sql
SELECT "name", "email" FROM "users" WHERE NOT EXISTS (SELECT * FROM "orders" WHERE "status" = $1);
-- Params: ["cancelled"]
```

## IN Subqueries

### Basic IN

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "total",
        Operator: ">",
        Value:    1000,
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery)

sql, params, err := builder.Select([]string{"name", "email"})
```

**MySQL Result:**
```sql
SELECT `name`, `email` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > ?);
-- Params: [1000]
```

### NOT IN

```go
subquery := sb.NewBuilder(sb.DIALECT_SQLITE).
    Table("banned_users").
    Where(&sb.Where{
        Column:   "reason",
        Operator: "=",
        Value:    "spam",
    })

builder, err := sb.NewBuilder(sb.DIALECT_SQLITE).
    Table("users").
    NotInSubquery(subquery)

sql, params, err := builder.Select([]string{"name", "email"})
```

**SQLite Result:**
```sql
SELECT "name", "email" FROM "users" WHERE "id" NOT IN (SELECT * FROM "banned_users" WHERE "reason" = ?);
-- Params: ["spam"]
```

## Correlated Subqueries

### Basic Correlated Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "user_id",
        Operator: "=",
        Value:    "users.id", // Reference outer query column
    }).
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(subquery)

sql, params, err := builder.Select([]string{"name", "email"})
```

**MySQL Result:**
```sql
SELECT `name`, `email` FROM `users` WHERE EXISTS (SELECT * FROM `orders` WHERE `user_id` = "users.id" AND `status` = ?);
-- Params: ["active"]
```

### Comparison with Correlated Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    Where(&sb.Where{
        Column:   "user_id",
        Operator: "=",
        Value:    "users.id",
    })

builder := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    Where(&sb.Where{
        Column:   "id",
        Operator: ">",
        Subquery: subquery.(*sb.Builder),
    })

sql, params, err := builder.Select([]string{"name", "email"})
```

**PostgreSQL Result:**
```sql
SELECT "name", "email" FROM "users" WHERE "id" > (SELECT * FROM "orders" WHERE "user_id" = "users.id");
-- Params: []
```

## Complex Subquery Examples

### Multiple WHERE Conditions with Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("order_items").
    Where(&sb.Where{
        Column:   "quantity",
        Operator: ">",
        Value:    5,
    }).
    Where(&sb.Where{
        Column:   "price",
        Operator: ">",
        Value:    100,
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "pending",
    }).
    InSubquery(subquery).
    OrderBy("created_at", "DESC").
    Limit(20)

sql, params, err := builder.Select([]string{"*"})
```

**MySQL Result:**
```sql
SELECT * FROM `orders` WHERE `status` = ? AND `id` IN (SELECT * FROM `order_items` WHERE `quantity` > ? AND `price` > ?) ORDER BY `created_at` DESC LIMIT ?;
-- Params: ["pending", 5, 100, 20]
```

### Subquery with JOIN

```go
subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("order_items").
    InnerJoin("products", "order_items.product_id = products.id").
    Where(&sb.Where{
        Column:   "products.category",
        Operator: "=",
        Value:    "electronics",
    })

builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    InSubquery(subquery).
    Select([]string{"id", "total"})
```

**PostgreSQL Result:**
```sql
SELECT "id", "total" FROM "orders" WHERE "id" IN (SELECT * FROM "order_items" INNER JOIN "products" ON order_items.product_id = products.id WHERE "products"."category" = $1);
-- Params: ["electronics"]
```

### Nested Subqueries

```go
innerSubquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("order_items").
    Where(&sb.Where{
        Column:   "quantity",
        Operator: ">",
        Value:    10,
    })

middleSubquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "shipped",
    }).
    InSubquery(innerSubquery)

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(middleSubquery).
    Select([]string{"name", "email"})
```

**MySQL Result:**
```sql
SELECT `name`, `email` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `status` = ? AND `id` IN (SELECT * FROM `order_items` WHERE `quantity` > ?));
-- Params: ["shipped", 10]
```

## Subquery Error Handling

### Handle Nil Subquery Error

```go
builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(nil)

if err != nil {
    if err == sb.ErrNilSubquery {
        fmt.Println("Subquery cannot be nil")
        return "", nil, err
    }
    return "", nil, err
}
```

### Handle Subquery Build Error

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(subquery)

if err != nil {
    return "", nil, fmt.Errorf("failed to add subquery: %w", err)
}

sql, params, err := builder.Select([]string{"name", "email"})
if err != nil {
    return "", nil, fmt.Errorf("failed to build query: %w", err)
}
```

## Database-Specific Considerations

### MySQL

- Uses backticks for identifiers
- Double quotes for string values
- Full EXISTS/IN/NOT IN support
- Correlated subqueries with proper quoting

### PostgreSQL

- Uses double quotes for identifiers
- Double quotes for string values (consistent)
- Full EXISTS/IN/NOT IN support
- Excellent correlated subquery performance

### SQLite

- Uses double quotes for identifiers
- Single quotes for string values
- Full EXISTS/IN/NOT IN support
- Good for simple correlated subqueries

### MSSQL

- Uses square brackets for identifiers
- Unquoted numeric values
- Full EXISTS/IN/NOT IN support
- Optimized for complex subqueries

## Performance Considerations

### EXISTS vs IN

**EXISTS** is generally faster for:
- Large datasets
- When you only need to check existence
- When the subquery returns many rows

**IN** is generally faster for:
- Small datasets
- When you need the actual values
- When the subquery returns few rows

### Correlated Subqueries

Correlated subqueries can be slow because:
- They execute once per row in the outer query
- They prevent certain query optimizations

**Optimization Tips:**
- Use EXISTS instead of IN when possible
- Ensure correlated columns are indexed
- Consider rewriting as JOINs for complex cases

### Subquery Caching

Some databases cache subquery results:
- MySQL: Limited subquery caching
- PostgreSQL: Good subquery optimization
- SQLite: Basic subquery caching
- MSSQL: Advanced subquery optimization

## Common Use Cases

### Find Users with Active Orders

```go
activeOrders := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "status",
        Operator: "=",
        Value:    "active",
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Exists(activeOrders).
    Select([]string{"name", "email"})
```

### Find Users NOT in Banned List

```go
bannedUsers := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("banned_users").
    Select([]string{"user_id"})

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    NotInSubquery(bannedUsers).
    Select([]string{"name", "email"})
```

### Find Users with Total Orders > Threshold

```go
highValueOrders := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{
        Column:   "total",
        Operator: ">",
        Value:    10000,
    })

builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(highValueOrders).
    Select([]string{"name", "email"})
```

## Best Practices

1. **Use EXISTS for existence checks** - More efficient than IN for large datasets
2. **Index correlated columns** - Improve subquery performance
3. **Limit subquery results** - Use LIMIT when appropriate
4. **Test with real data** - Subquery performance varies by data distribution
5. **Consider JOINs** - Sometimes JOINs are more efficient than subqueries

## See Also

- [API Reference](../api_reference.md)
- [JOIN Operations](joins.md)
- [Parameterized Queries](parameterized-queries.md)
