# View Management

## CREATE VIEW

Create database views for simplified data access and security.

### Basic View Creation

```go
selectSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	Select([]string{"FirstName", "LastName"})

createViewSql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	View("v_users").
	ViewColumns([]string{"first_name", "last_name"}).
	ViewSQL(selectSQL).
	Create()
```

### View with Complex Query

```go
selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "completed"}).
	Select([]string{
		"orders.id",
		"orders.total",
		"users.name as customer_name",
		"orders.created_at",
	})

createViewSql := sb.NewBuilder(sb.DIALECT_MYSQL).
	View("v_completed_orders").
	ViewSQL(selectSQL).
	Create()
```

### View with Aggregation

```go
selectSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	Select([]string{
		"user_id",
		"COUNT(*) as order_count",
		"SUM(total) as total_spent",
	}).
	GroupBy("user_id")

createViewSql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	View("v_user_order_summary").
	ViewSQL(selectSQL).
	Create()
```

## CREATE VIEW IF NOT EXISTS

Create views safely without errors if they already exist.

```go
selectSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	Select([]string{"FirstName", "LastName"})

createViewSql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	View("v_users").
	ViewColumns([]string{"first_name", "last_name"}).
	ViewSQL(selectSQL).
	CreateIfNotExists()
```

## DROP VIEW

Remove views from the database.

### Basic View Drop

```go
dropViewSql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	View("v_users").
	Drop()
```

### Safe View Drop

```go
dropViewSql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	View("v_users").
	DropIfExists()
```

## Database-Specific Behavior

### MySQL

```sql
-- Create View
CREATE VIEW `view_name` AS SELECT ...;

-- Create View with Columns
CREATE VIEW `view_name` (`col1`, `col2`) AS SELECT ...;

-- Create View IF NOT EXISTS
CREATE VIEW IF NOT EXISTS `view_name` AS SELECT ...;

-- Drop View
DROP VIEW `view_name`;

-- Drop View IF EXISTS
DROP VIEW IF EXISTS `view_name`;
```

### PostgreSQL

```sql
-- Create View
CREATE VIEW "view_name" AS SELECT ...;

-- Create View with Columns
CREATE VIEW "view_name" ("col1", "col2") AS SELECT ...;

-- Create View IF NOT EXISTS
CREATE VIEW IF NOT EXISTS "view_name" AS SELECT ...;

-- Drop View
DROP VIEW "view_name";

-- Drop View IF EXISTS
DROP VIEW IF EXISTS "view_name";
```

### SQLite

```sql
-- Create View
CREATE VIEW "view_name" AS SELECT ...;

-- Create View with Columns
CREATE VIEW "view_name" ("col1", "col2") AS SELECT ...;

-- Create View IF NOT EXISTS
CREATE VIEW IF NOT EXISTS "view_name" AS SELECT ...;

-- Drop View
DROP VIEW "view_name";

-- Drop View IF EXISTS
DROP VIEW IF EXISTS "view_name";
```

### MSSQL

```sql
-- Create View
CREATE VIEW [view_name] AS SELECT ...;

-- Create View with Columns
CREATE VIEW [view_name] ([col1], [col2]) AS SELECT ...;

-- Drop View
DROP VIEW [view_name];

-- Drop View IF EXISTS (SQL Server 2016+)
DROP VIEW IF EXISTS [view_name];
```

## View Methods Reference

### Creation Methods

- `View(viewName)` - Set view name
- `ViewColumns(columns...)` - Set view column names
- `ViewSQL(selectBuilder)` - Set the SELECT query for the view
- `Create()` - Create the view
- `CreateIfNotExists()` - Create view if it doesn't exist

### Drop Methods

- `Drop()` - Drop the view
- `DropIfExists()` - Drop view if it exists

## Common View Patterns

### User-Friendly Data Views

```go
// Simplified user profile view
selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	LeftJoin("profiles", "users.id = profiles.user_id").
	Select([]string{
		"users.id",
		"users.name",
		"users.email",
		"profiles.avatar",
		"profiles.bio",
	})

createViewSql := sb.NewBuilder(sb.DIALECT_MYSQL).
	View("v_user_profiles").
	ViewSQL(selectSQL).
	Create()
```

### Security Views

```go
// Limited order information for customer service
selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	Select([]string{
		"id",
		"user_id",
		"status",
		"total",
		"created_at",
	}).
	Where(&sb.Where{Column: "status", Operator: "!=", Value: "deleted"})

createViewSql := sb.NewBuilder(sb.DIALECT_MYSQL).
	View("v_order_summary").
	ViewSQL(selectSQL).
	Create()
```

### Reporting Views

```go
// Monthly sales report view
selectSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	Select([]string{
		"DATE_TRUNC('month', created_at) as month",
		"COUNT(*) as order_count",
		"SUM(total) as total_revenue",
		"AVG(total) as avg_order_value",
	}).
	Where(&sb.Where{Column: "status", Operator: "=", Value: "completed"}).
	GroupBy("DATE_TRUNC('month', created_at)").
	OrderBy("month", "DESC")

createViewSql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	View("v_monthly_sales").
	ViewSQL(selectSQL).
	Create()
```

### Dashboard Views

```go
// User dashboard statistics
selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	LeftJoin("(SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) as order_stats", "users.id = order_stats.user_id").
	Select([]string{
		"users.id",
		"users.name",
		"COALESCE(order_stats.order_count, 0) as order_count",
	})

createViewSql := sb.NewBuilder(sb.DIALECT_MYSQL).
	View("v_user_dashboard").
	ViewSQL(selectSQL).
	Create()
```

## View Management Examples

### Complete View Lifecycle

```go
func SetupReportingViews(db *sb.Database) error {
    ctx := context.Background()
    
    // Monthly sales view
    monthlySalesSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
        Table("orders").
        Select([]string{
            "DATE_TRUNC('month', created_at) as month",
            "COUNT(*) as order_count",
            "SUM(total) as total_revenue",
        }).
        Where(&sb.Where{Column: "status", Operator: "=", Value: "completed"}).
        GroupBy("DATE_TRUNC('month', created_at)")

    createMonthlySales := sb.NewBuilder(sb.DIALECT_POSTGRES).
        View("v_monthly_sales").
        ViewSQL(monthlySalesSQL).
        CreateIfNotExists()

    // User order summary view
    userOrderSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).
        Table("orders").
        Select([]string{
            "user_id",
            "COUNT(*) as order_count",
            "SUM(total) as total_spent",
        }).
        GroupBy("user_id")

    createUserOrders := sb.NewBuilder(sb.DIALECT_POSTGRES).
        View("v_user_orders").
        ViewSQL(userOrderSQL).
        CreateIfNotExists()

    // Create views
    views := []string{createMonthlySales, createUserOrders}
    
    for _, sql := range views {
        if _, err := db.Exec(ctx, sql); err != nil {
            return err
        }
    }
    
    return nil
}
```

### View Refresh Pattern

```go
func RefreshMaterializedViews(db *sb.Database) error {
    ctx := context.Background()
    
    // Drop existing views
    dropViews := []string{
        sb.NewBuilder(sb.DIALECT_POSTGRES).View("v_monthly_sales").DropIfExists(),
        sb.NewBuilder(sb.DIALECT_POSTGRES).View("v_user_orders").DropIfExists(),
    }
    
    for _, sql := range dropViews {
        db.Exec(ctx, sql) // Ignore errors
    }
    
    // Recreate views
    return SetupReportingViews(db)
}
```

## Best Practices

### View Design Principles

1. **Keep views simple** - Complex views are hard to maintain
2. **Use meaningful names** - Prefix with `v_` to distinguish views
3. **Document view purpose** - Comments explaining business logic
4. **Avoid nested views** - Views referencing other views can cause performance issues
5. **Consider performance** - Views don't automatically improve query performance

### Security Considerations

```go
// Restrictive view for customer service (no sensitive data)
selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Select([]string{
		"id",
		"name",
		"email",
		"created_at",
	}).
	Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

createViewSql := sb.NewBuilder(sb.DIALECT_MYSQL).
	View("v_customer_service_users").
	ViewSQL(selectSQL).
	Create()
```

### Performance Optimization

```go
// Indexed view-friendly query
selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "completed"}).
	Select([]string{
		"orders.id",
		"users.name as customer_name",
		"orders.total",
		"orders.created_at",
	})

createViewSql := sb.NewBuilder(sb.DIALECT_MYSQL).
	View("v_completed_orders").
	ViewSQL(selectSQL).
	Create()
```

## Common Use Cases

### Data Abstraction
- Hide complex table structures
- Provide simplified interfaces
- Encapsulate business logic

### Security
- Restrict access to sensitive columns
- Implement row-level security
- Provide read-only access

### Reporting
- Pre-computed aggregations
- Simplified reporting queries
- Consistent data formatting

### Legacy Support
- Maintain API compatibility during schema changes
- Provide backward-compatible interfaces
- Abstract schema migrations
