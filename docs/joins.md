# JOIN Operations

## Overview

SB SQL builder provides comprehensive JOIN support with all major JOIN types, table aliases, and database-specific dialect handling while maintaining simplicity and consistency.

## Enhanced JOIN Support (Complete Implementation)

The JOIN implementation supports advanced query building with proper database dialect handling and comprehensive testing.

### Quick Examples

#### Basic INNER JOIN
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Select([]string{"orders.*", "users.name"})
// Result: SELECT `orders`.*, `users`.`name` FROM `orders` INNER JOIN `users` ON orders.user_id = users.id;
```

#### JOIN with Table Alias
```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
    Select([]string{"orders.*", "p.avatar"})
// Result: SELECT "orders".*, "p"."avatar" FROM "orders" LEFT JOIN "profiles" AS "p" ON orders.user_id = p.user_id;
```

#### Multiple JOINs with Conditions
```go
sql := sb.NewBuilder(sb.DIALECT_SQLITE).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    LeftJoin("profiles", "users.id = profiles.user_id").
    Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
    OrderBy("orders.created_at", "DESC").
    Select([]string{"orders.total", "users.name", "profiles.avatar"})
```

---

## API Reference

### Core JOIN Methods

#### Join()
```go
Join(joinType JoinType, table string, onCondition string) BuilderInterface
```
Generic JOIN method supporting all join types.

**Parameters:**
- `joinType` - JOIN_INNER, JOIN_LEFT, JOIN_RIGHT, JOIN_FULL, JOIN_CROSS
- `table` - Table name to join
- `onCondition` - JOIN condition

#### JoinWithAlias()
```go
JoinWithAlias(joinType JoinType, table string, alias string, onCondition string) BuilderInterface
```
JOIN with table alias support.

**Parameters:**
- `joinType` - JOIN type
- `table` - Table name to join
- `alias` - Table alias
- `onCondition` - JOIN condition

#### Convenience Methods
```go
InnerJoin(table string, onCondition string) BuilderInterface
LeftJoin(table string, onCondition string) BuilderInterface
RightJoin(table string, onCondition string) BuilderInterface
```

### JOIN Types
```go
const (
    JOIN_INNER JoinType = "INNER"
    JOIN_LEFT  JoinType = "LEFT"
    JOIN_RIGHT JoinType = "RIGHT"
    JOIN_FULL  JoinType = "FULL"
    JOIN_CROSS JoinType = "CROSS"
)
```

---

## Database-Specific Behavior

### MySQL
- **Backtick quoting**: `table_name`, `column_name`
- **All JOIN types supported**: INNER, LEFT, RIGHT, FULL, CROSS
- **Table aliases**: AS clause supported
- **Multiple JOINs**: Fully supported

```sql
SELECT `orders`.*, `users`.`name` 
FROM `orders` 
INNER JOIN `users` ON orders.user_id = users.id 
LEFT JOIN `profiles` AS `p` ON orders.user_id = p.user_id;
```

### PostgreSQL
- **Double quote quoting**: "table_name", "column_name"
- **All JOIN types supported**: INNER, LEFT, RIGHT, FULL, CROSS
- **Table aliases**: AS clause supported
- **Advanced features**: Fully supported

```sql
SELECT "orders".*, "users"."name" 
FROM "orders" 
INNER JOIN "users" ON orders.user_id = users.id 
LEFT JOIN "profiles" AS "p" ON orders.user_id = p.user_id;
```

### SQLite
- **Double quote quoting**: "table_name", "column_name"
- **JOIN types supported**: INNER, LEFT, RIGHT, CROSS
- **FULL OUTER JOIN**: Limited support
- **Table aliases**: AS clause supported

```sql
SELECT "orders".*, "users"."name" 
FROM "orders" 
INNER JOIN "users" ON orders.user_id = users.id 
LEFT JOIN "profiles" AS "p" ON orders.user_id = p.user_id;
```

### MSSQL
- **Square bracket quoting**: [table_name], [column_name]
- **All JOIN types supported**: INNER, LEFT, RIGHT, FULL, CROSS
- **Table aliases**: AS clause supported
- **Multiple JOINs**: Fully supported

```sql
SELECT [orders].*, [users].[name] 
FROM [orders] 
INNER JOIN [users] ON orders.user_id = users.id 
LEFT JOIN [profiles] AS [p] ON orders.user_id = p.user_id;
```

---

## Usage Patterns

### Multiple JOINs
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    LeftJoin("profiles", "users.id = profiles.user_id").
    LeftJoin("addresses", "users.address_id = addresses.id").
    Select([]string{
        "orders.total",
        "users.name", 
        "users.email",
        "profiles.avatar",
        "addresses.city",
    })
```

### JOIN with Subqueries
```go
// Subquery for JOIN condition
highValueOrders := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "total", Operator: ">", Value: 1000}).
    Select([]string{"user_id"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InnerJoin("orders", "users.id = orders.user_id").
    InSubquery(highValueOrders).
    Select([]string{"users.name", "COUNT(orders.id) as order_count"})
```

### Complex Query with Multiple Conditions
```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    LeftJoin("users", "orders.user_id = users.id").
    LeftJoin("products", "orders.product_id = products.id").
    Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
    Where(&sb.Where{Column: "users.verified", Operator: "=", Value: true}).
    Where(&sb.Where{Column: "products.available", Operator: "=", Value: true}).
    OrderBy("orders.created_at", "DESC").
    Limit(20).
    Select([]string{"orders.*", "users.name", "products.name"})
```

### Self-Join with Alias
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("employees").
    JoinWithAlias(sb.JOIN_LEFT, "employees", "manager", "employees.manager_id = manager.id").
    Select([]string{"employees.name", "manager.name as manager_name"})
// Result: SELECT `employees`.`name`, `manager`.`name` as `manager_name` 
// FROM `employees` LEFT JOIN `employees` AS `manager` ON employees.manager_id = manager.id;
```

---

## Advanced Features

### JOIN in UPDATE Queries
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Where(&sb.Where{Column: "users.id", Operator: "=", Value: 1}).
    Update(map[string]string{"orders.status": "processed"})
// Result: UPDATE `orders` INNER JOIN `users` ON orders.user_id = users.id 
// SET `orders`.`status` = 'processed' WHERE `users`.`id` = 1;
```

### JOIN with Aggregations
```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    InnerJoin("orders", "users.id = orders.user_id").
    GroupBy("users.id", "users.name").
    OrderBy("order_count", "DESC").
    Select([]string{"users.name", "COUNT(orders.id) as order_count"})
```

### Conditional JOIN Patterns
```go
// Different JOIN based on conditions
baseBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).Table("orders")

if includeProfiles {
    baseBuilder.LeftJoin("profiles", "orders.user_id = profiles.user_id")
}

if includeDetails {
    baseBuilder.InnerJoin("order_details", "orders.id = order_details.order_id")
}

sql, _ := baseBuilder.Select([]string{"orders.*"})
```

---

## Best Practices

### Performance Considerations
```go
// ✅ Good: JOIN on indexed columns
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").  // user_id should be indexed
    Select([]string{"orders.total", "users.name"})

// ❌ Avoid: JOIN on non-indexed columns in large tables
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.notes = users.description").  // Poor performance
    Select([]string{"orders.total", "users.name"})
```

### JOIN Order
```go
// Best practice: Start with main table, then join related tables
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").                                    // Main table
    LeftJoin("profiles", "users.id = profiles.user_id").    // Related data
    LeftJoin("user_settings", "users.id = user_settings.user_id"). // More related data
    Where(&sb.Where{Column: "users.id", Operator: "=", Value: userID}).
    Select([]string{
        "users.id", "users.name", "users.email",
        "profiles.avatar", "user_settings.theme",
    })
```

### Complex Multi-Table Queries
```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    LeftJoin("addresses", "orders.billing_address_id = addresses.id").
    Where(&sb.Where{Column: "orders.id", Operator: "=", Value: orderID}).
    Select([]string{
        "orders.*",
        "users.name as customer_name",
        "users.email as customer_email",
        "addresses.street",
        "addresses.city",
        "addresses.country",
    })
```

---

## Error Handling

JOIN methods follow the standard SB error handling pattern:

```go
builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id")

if err != nil {
    return "", err
}

sql, err := builder.Select([]string{"orders.*", "users.name"})
if err != nil {
    return "", err
}
```

Common JOIN-related errors:
- Empty table names
- Invalid JOIN conditions
- Missing required parameters

---

## Testing

The JOIN implementation includes **13 comprehensive tests** covering:

- All JOIN types across all database dialects
- Table alias functionality
- Multiple JOIN scenarios
- Error handling cases
- Integration with WHERE, ORDER BY, LIMIT
- UPDATE queries with JOINs

Run tests:
```bash
go test -v -run "TestBuilder.*Join.*"
```

---

## Migration Notes

### From Manual SQL to SB JOINs

**Before (Manual SQL):**
```go
sql := "SELECT orders.*, users.name FROM orders INNER JOIN users ON orders.user_id = users.id WHERE orders.status = ?"
```

**After (SB JOINs):**
```go
sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id").
    Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
    Select([]string{"orders.*", "users.name"})
```

Benefits:
- **Database-agnostic** syntax
- **Parameterized queries** by default
- **Type safety** with structured conditions
- **Composable** query building

---

## Conclusion

The JOIN implementation provides comprehensive, database-agnostic JOIN support while maintaining SB's philosophy of simplicity and consistency. The implementation supports all major JOIN patterns with proper dialect handling and extensive testing.

For detailed implementation information, see the memory bank documentation for the JOIN support implementation.
