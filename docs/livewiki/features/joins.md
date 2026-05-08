---
path: features/joins.md
page-type: tutorial
summary: Complete guide to JOIN operations in SB SQL Builder including all join types and table aliases.
tags: [joins, inner-join, left-join, right-join, table-aliases]
created: 2026-05-08
updated: 2026-05-08
version: 1.0.0
---

# JOIN Operations

SB provides comprehensive JOIN support including INNER, LEFT, RIGHT, FULL, and CROSS joins with table aliases and multiple join support.

## Overview

JOIN operations allow you to combine rows from two or more tables based on a related column between them.

### Supported Join Types

- **INNER JOIN**: Returns rows when there's a match in both tables
- **LEFT JOIN**: Returns all rows from the left table, and matched rows from the right table
- **RIGHT JOIN**: Returns all rows from the right table, and matched rows from the left table
- **FULL JOIN**: Returns rows when there's a match in either table
- **CROSS JOIN**: Returns the Cartesian product of the two tables

## Basic JOIN Syntax

### INNER JOIN

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**MySQL Result:**
```sql
SELECT `orders`.*, `users`.`name` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id;
```

**PostgreSQL Result:**
```sql
SELECT "orders".*, "users"."name" FROM "orders" INNER JOIN "users" ON orders.user_id = users.id;
```

### LEFT JOIN

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_LEFT, "users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**MySQL Result:**
```sql
SELECT `orders`.*, `users`.`name` FROM `orders` LEFT JOIN `users` ON orders.user_id = users.id;
```

### RIGHT JOIN

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_RIGHT, "users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**MySQL Result:**
```sql
SELECT `orders`.*, `users`.`name` FROM `orders` RIGHT JOIN `users` ON orders.user_id = users.id;
```

## Table Aliases

### Join with Alias

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
    Select([]string{"orders.*", "p.avatar"})
```

**MySQL Result:**
```sql
SELECT `orders`.*, `p`.`avatar` FROM `orders` LEFT JOIN `profiles` AS `p` ON orders.user_id = p.user_id;
```

**PostgreSQL Result:**
```sql
SELECT "orders".*, "p"."avatar" FROM "orders" LEFT JOIN "profiles" AS "p" ON orders.user_id = p.user_id;
```

## Convenience Methods

### InnerJoin

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

### LeftJoin

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    LeftJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

### RightJoin

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    RightJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

## Multiple JOINs

### Two Joins

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    LeftJoin("profiles", "users.id = profiles.user_id").
    Select([]string{"orders.total", "users.name", "profiles.avatar"})
```

**MySQL Result:**
```sql
SELECT `orders`.`total`, `users`.`name`, `profiles`.`avatar` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id LEFT JOIN `profiles` ON users.id = profiles.user_id;
```

### Three Joins

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    LeftJoin("profiles", "users.id = profiles.user_id").
    LeftJoin("addresses", "users.id = addresses.user_id").
    Select([]string{"orders.*", "users.name", "profiles.avatar", "addresses.city"})
```

## JOIN with WHERE and ORDER BY

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    LeftJoin("users", "orders.user_id = users.id").
    Where(&sb.Where{
        Column:   "orders.status",
        Operator: "=",
        Value:    "active",
    }).
    OrderBy("orders.created_at", "DESC").
    Select([]string{"orders.*", "users.name"})
```

**PostgreSQL Result:**
```sql
SELECT "orders".*, "users"."name" FROM "orders" LEFT JOIN "users" ON orders.user_id = users.id WHERE "orders"."status" = $1 ORDER BY "orders"."created_at" DESC;
-- Params: ["active"]
```

## JOIN with Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("order_items").
    Where(&sb.Where{
        Column:   "quantity",
        Operator: ">",
        Value:    10,
    })

sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Where(&sb.Where{
        Column:   "orders.id",
        Operator: "IN",
        Subquery: subquery.(*sb.Builder),
    }).
    Select([]string{"orders.*", "users.name"})
```

## JOIN in UPDATE Queries

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Where(&sb.Where{
        Column:   "users.id",
        Operator: "=",
        Value:    1,
    }).
    Update(map[string]interface{}{
        "orders.status": "processed",
    })
```

**MySQL Result:**
```sql
UPDATE `orders` INNER JOIN `users` ON orders.user_id = users.id SET `orders`.`status` = ? WHERE `users`.`id` = ?;
-- Params: ["processed", 1]
```

## Database-Specific Syntax

### MySQL

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**Result:**
```sql
SELECT `orders`.*, `users`.`name` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id;
```

### PostgreSQL

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**Result:**
```sql
SELECT "orders".*, "users"."name" FROM "orders" INNER JOIN "users" ON orders.user_id = users.id;
```

### SQLite

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_SQLITE).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**Result:**
```sql
SELECT "orders".*, "users"."name" FROM "orders" INNER JOIN "users" ON orders.user_id = users.id;
```

### MSSQL

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MSSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
```

**Result:**
```sql
SELECT [orders].*, [users].[name] FROM [orders] INNER JOIN [users] ON orders.user_id = users.id;
```

## Complex JOIN Examples

### Self-Join

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("employees").
    JoinWithAlias(sb.JOIN_INNER, "employees", "manager", "employees.manager_id = manager.id").
    Select([]string{"employees.name", "manager.name AS manager_name"})
```

**MySQL Result:**
```sql
SELECT `employees`.`name`, `manager`.`name` AS manager_name FROM `employees` INNER JOIN `employees` AS `manager` ON employees.manager_id = manager.id;
```

### JOIN with Multiple Conditions

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "orders.user_id = users.id AND users.status = 'active'").
    Select([]string{"orders.*", "users.name"})
```

**MySQL Result:**
```sql
SELECT `orders`.*, `users`.`name` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id AND users.status = 'active';
```

### JOIN with Aggregation

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    LeftJoin("orders", "users.id = orders.user_id").
    Select([]string{"users.name", "COUNT(orders.id) AS order_count"})
```

**MySQL Result:**
```sql
SELECT `users`.`name`, COUNT(orders.id) AS order_count FROM `users` LEFT JOIN `orders` ON users.id = orders.user_id;
```

## Error Handling

### Empty ON Condition

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "") // Empty ON condition

sql, params, err := builder.Select([]string{"*"})
// Error: ValidationError: ON condition cannot be empty
```

### Handle Join Error

```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "orders.user_id = users.id")

sql, params, err := builder.Select([]string{"*"})
if err != nil {
    if err == sb.ErrEmptyOnCondition {
        fmt.Println("ON condition is required")
        return "", nil, err
    }
    return "", nil, err
}
```

## Performance Considerations

### Indexing Join Columns

Ensure join columns are indexed for optimal performance:

```sql
CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_users_id ON users(id);
```

### Join Order Matters

- **Smaller table first**: Join smaller tables first when possible
- **Selective conditions**: Apply WHERE conditions before JOINs when possible
- **Index usage**: Ensure indexes are used by checking query plans

### JOIN vs Subquery

- **JOIN**: Generally faster for simple relationships
- **Subquery**: Sometimes clearer for complex conditions
- **Test both**: Performance varies by data distribution

## Common Use Cases

### Get Orders with User Details

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name", "users.email"})
```

### Get Users with Order Count

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    LeftJoin("orders", "users.id = orders.user_id").
    Select([]string{"users.name", "COUNT(orders.id) AS order_count"})
```

### Get Products with Categories

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("products").
    InnerJoin("categories", "products.category_id = categories.id").
    Select([]string{"products.*", "categories.name AS category_name"})
```

### Get Users with Profile and Address

```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    LeftJoin("profiles", "users.id = profiles.user_id").
    LeftJoin("addresses", "users.id = addresses.user_id").
    Select([]string{"users.*", "profiles.avatar", "addresses.city"})
```

## Best Practices

1. **Use INNER JOIN when you need matches in both tables** - Most common use case
2. **Use LEFT JOIN when you need all rows from the left table** - Preserves left table data
3. **Always specify ON condition** - Required for most JOIN types
4. **Use table aliases for clarity** - Makes queries more readable
5. **Index join columns** - Improves performance significantly
6. **Limit JOINed columns** - Only select needed columns to reduce data transfer
7. **Test with real data** - JOIN performance varies by data distribution

## See Also

- [API Reference](../api_reference.md)
- [Subqueries](subqueries.md)
- [Parameterized Queries](parameterized-queries.md)
