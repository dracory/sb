# SELECT Queries and Subqueries

## Basic SELECT

Retrieve data from tables with column selection and filtering.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Select([]string{"id", "name", "email"})
```

## WHERE Conditions

Filter results using various operators and conditions.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{
		Column: "status",
		Operator: "=",
		Value: "active",
	}).
	Select([]string{"name", "email"})
```

### Where Operators

- `=` - Equal to
- `!=` or `<>` - Not equal to
- `>` - Greater than
- `<` - Less than
- `>=` - Greater than or equal to
- `<=` - Less than or equal to
- `LIKE` - Pattern matching
- `IN` - In list of values
- `NOT IN` - Not in list of values

## ORDER BY

Sort query results.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	OrderBy("created_at", "DESC").
	OrderBy("name", "ASC").
	Select([]string{"*"})
```

## LIMIT

Limit the number of returned rows.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Limit(10).
	Select([]string{"*"})
```

## Subqueries

Use nested queries for complex data retrieval.

### IN Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	InSubquery(subquery).
	Select([]string{"name"})
// Result: SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > "1000");
```

### EXISTS Subquery

```go
activeOrders := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	Exists(activeOrders).
	Select([]string{"name", "email"})
// Result: SELECT "name", "email" FROM "users" WHERE EXISTS (SELECT * FROM "orders" WHERE "status" = "active");
```

### NOT EXISTS Subquery

```go
inactiveOrders := sb.NewBuilder(sb.DIALECT_SQLITE).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"})

sql := sb.NewBuilder(sb.DIALECT_SQLITE).
	Table("users").
	NotExists(inactiveOrders).
	Select([]string{"name"})
// Result: SELECT "name" FROM "users" WHERE NOT EXISTS (SELECT * FROM "orders" WHERE "status" = 'inactive');
```

### NOT IN Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_MSSQL).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "cancelled"})

sql := sb.NewBuilder(sb.DIALECT_MSSQL).
	Table("users").
	NotInSubquery(subquery).
	Select([]string{"name"})
// Result: SELECT [name] FROM [users] WHERE [id] NOT IN (SELECT * FROM [orders] WHERE [status] = 'cancelled');
```

### Correlated Subquery

```go
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"}).
	Where(&sb.Where{Column: "total", Operator: ">", Value: "5000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Where(&sb.Where{
		Column:   "id",
		Operator: ">",
		Subquery: subquery.(*sb.Builder),
	}).
	Select([]string{"name"})
// Result: SELECT `name` FROM `users` WHERE `id` > (SELECT * FROM `orders` WHERE `user_id` = "users.id" AND `total` > "5000");
```

### Complex Subquery with Multiple Conditions

```go
subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("order_items").
	Where(&sb.Where{Column: "quantity", Operator: ">", Value: "5"}).
	Where(&sb.Where{Column: "price", Operator: ">", Value: "100"})

sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
	InSubquery(subquery).
	OrderBy("created_at", "DESC").
	Limit(20).
	Select([]string{"*"})
// Result: SELECT * FROM "orders" WHERE "status" = "active" AND "id" IN (SELECT * FROM "order_items" WHERE "quantity" > "5" AND "price" > "100") ORDER BY "created_at" DESC LIMIT 20;
```

### Using Subquery() Method

```go
subqueryBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).Subquery()
subquery := subqueryBuilder.
	Table("orders").
	Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	InSubquery(subquery).
	Select([]string{"name"})
```

## Subquery Methods

- `Subquery()` - Create a new subquery builder
- `Exists(subquery)` - Add EXISTS condition
- `NotExists(subquery)` - Add NOT EXISTS condition
- `InSubquery(subquery)` - Add IN condition with subquery
- `NotInSubquery(subquery)` - Add NOT IN condition with subquery

## Database-Specific Quoting

- **MySQL**: Backticks for identifiers (`table`), double quotes for values (`"value"`)
- **PostgreSQL**: Double quotes for both identifiers and values (`"table"`, `"value"`)
- **SQLite**: Double quotes for identifiers (`"table"`), single quotes for values (`'value'`)
- **MSSQL**: Square brackets for identifiers (`[table]`), unquoted numeric values

## Complex Queries

Combine multiple clauses for sophisticated queries.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
	Where(&sb.Where{Column: "total", Operator: ">", Value: "100"}).
	OrderBy("created_at", "DESC").
	Limit(10).
	Select([]string{"id", "user_id", "total", "created_at"})
```
