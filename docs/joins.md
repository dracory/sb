# JOIN Operations

## Basic JOIN Types

SB supports all major JOIN types with proper database dialect handling.

### INNER JOIN

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	Select([]string{"orders.*", "users.name"})
// Result: SELECT `orders`.*, `users`.`name` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id;
```

### LEFT JOIN

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	LeftJoin("profiles", "orders.user_id = profiles.user_id").
	Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
	Select([]string{"orders.*", "profiles.avatar"})
// Result: SELECT "orders".*, "profiles"."avatar" FROM "orders" LEFT JOIN "profiles" ON orders.user_id = profiles.user_id WHERE "orders"."status" = 'active';
```

### RIGHT JOIN

```go
sql := sb.NewBuilder(sb.DIALECT_MSSQL).
	Table("orders").
	RightJoin("users", "orders.user_id = users.id").
	Select([]string{"orders.*", "users.name"})
// Result: SELECT [orders].*, [users].[name] FROM [orders] RIGHT JOIN [users] ON orders.user_id = users.id;
```

### FULL OUTER JOIN

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	Join(sb.JOIN_FULL, "users", "orders.user_id = users.id").
	Select([]string{"orders.*", "users.name"})
// Result: SELECT "orders".*, "users"."name" FROM "orders" FULL OUTER JOIN "users" ON orders.user_id = users.id;
```

### CROSS JOIN

```go
sql := sb.NewBuilder(sb.DIALECT_SQLITE).
	Table("orders").
	Join(sb.JOIN_CROSS, "users", "1=1").
	Select([]string{"orders.*", "users.name"})
// Result: SELECT "orders".*, "users"."name" FROM "orders" CROSS JOIN "users" ON 1=1;
```

## JOIN with Table Aliases

Use table aliases for cleaner queries and self-joins.

### Basic Table Alias

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
	Select([]string{"orders.*", "p.avatar"})
// Result: SELECT `orders`.*, `p`.`avatar` FROM `orders` LEFT JOIN `profiles` AS `p` ON orders.user_id = p.user_id;
```

### Self-Join with Alias

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("employees").
	JoinWithAlias(sb.JOIN_LEFT, "employees", "manager", "employees.manager_id = manager.id").
	Select([]string{"employees.name", "manager.name as manager_name"})
// Result: SELECT `employees`.`name`, `manager`.`name` as `manager_name` FROM `employees` LEFT JOIN `employees` AS `manager` ON employees.manager_id = manager.id;
```

## Multiple JOINs

Combine multiple JOINs in a single query.

```go
sql := sb.NewBuilder(sb.DIALECT_SQLITE).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	LeftJoin("profiles", "users.id = profiles.user_id").
	Select([]string{"orders.total", "users.name", "profiles.avatar"})
// Result: SELECT "orders"."total", "users"."name", "profiles"."avatar" FROM "orders" INNER JOIN "users" ON orders.user_id = users.id LEFT JOIN "profiles" ON users.id = profiles.user_id;
```

## Complex JOIN Queries

### JOIN with WHERE and ORDER BY

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	LeftJoin("users", "orders.user_id = users.id").
	Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
	Where(&sb.Where{Column: "users.status", Operator: "=", Value: "verified"}).
	OrderBy("orders.created_at", "DESC").
	Select([]string{"orders.*", "users.name", "users.email"})
```

### JOIN with Subquery

```go
// Subquery for high-value orders
highValueOrders := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"}).
	Select([]string{"user_id"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	InnerJoin("orders", "users.id = orders.user_id").
	InSubquery(highValueOrders).
	Select([]string{"users.name", "COUNT(orders.id) as order_count"})
```

## JOIN Types Reference

### JoinType Constants

```go
const (
    JOIN_INNER JoinType = "INNER"
    JOIN_LEFT  JoinType = "LEFT"
    JOIN_RIGHT JoinType = "RIGHT"
    JOIN_FULL  JoinType = "FULL"
    JOIN_CROSS JoinType = "CROSS"
)
```

### Method Reference

- `Join(joinType, table, onCondition)` - Generic JOIN method
- `JoinWithAlias(joinType, table, alias, onCondition)` - JOIN with table alias
- `InnerJoin(table, onCondition)` - Convenience method for INNER JOIN
- `LeftJoin(table, onCondition)` - Convenience method for LEFT JOIN
- `RightJoin(table, onCondition)` - Convenience method for RIGHT JOIN

## Database-Specific Behavior

### MySQL
- Uses backticks for identifiers: `table`, `column`
- Supports all JOIN types
- Table alias syntax: `table AS alias`

```sql
SELECT `orders`.*, `users`.`name` 
FROM `orders` 
INNER JOIN `users` ON orders.user_id = users.id;
```

### PostgreSQL
- Uses double quotes for identifiers: "table", "column"
- Supports all JOIN types including FULL OUTER
- Table alias syntax: "table" AS "alias"

```sql
SELECT "orders".*, "users"."name" 
FROM "orders" 
INNER JOIN "users" ON orders.user_id = users.id;
```

### SQLite
- Uses double quotes for identifiers: "table", "column"
- Supports all JOIN types
- Table alias syntax: "table" AS "alias"

```sql
SELECT "orders".*, "users"."name" 
FROM "orders" 
INNER JOIN "users" ON orders.user_id = users.id;
```

### MSSQL
- Uses square brackets for identifiers: [table], [column]
- Supports all JOIN types
- Table alias syntax: [table] AS [alias]

```sql
SELECT [orders].*, [users].[name] 
FROM [orders] 
INNER JOIN [users] ON orders.user_id = users.id;
```

## Best Practices

1. **Use explicit column names** instead of `*` in production
2. **Choose appropriate JOIN type** based on data requirements
3. **Use table aliases** for self-joins and complex queries
4. **Optimize JOIN conditions** with proper indexing
5. **Consider performance** with large datasets

## Performance Tips

```go
// Good: Use indexed columns in JOIN conditions
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").  // user_id should be indexed
	Select([]string{"orders.total", "users.name"})

// Avoid: JOIN on non-indexed columns in large tables
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.notes = users.description").  // Poor performance
	Select([]string{"orders.total", "users.name"})
```

## Common Patterns

### User Profile Query

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	LeftJoin("profiles", "users.id = profiles.user_id").
	LeftJoin("user_settings", "users.id = user_settings.user_id").
	Where(&sb.Where{Column: "users.id", Operator: "=", Value: userID}).
	Select([]string{
		"users.id", "users.name", "users.email",
		"profiles.avatar", "profiles.bio",
		"user_settings.theme", "user_settings.notifications",
	})
```

### Order with Customer Details

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	LeftJoin("addresses", "orders.billing_address_id = addresses.id").
	Where(&sb.Where{Column: "orders.id", Operator: "=", Value: orderID}).
	Select([]string{
		"orders.*",
		"users.name as customer_name", "users.email as customer_email",
		"addresses.address_line1", "addresses.city", "addresses.country",
	})
```
