# SB <a href="https://gitpod.io/#https://github.com/dracory/sb" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/dracory/sb/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/sb)](https://goreportcard.com/report/github.com/dracory/sb)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/sb)](https://pkg.go.dev/github.com/dracory/sb)

A simplified SQL builder with comprehensive database operations support.

For a full SQL builder functionality check: https://doug-martin.github.io/goqu

Includes a wrapper for the mainstream DB package to allow transparent working with transactions.

**Features:**
- ✅ **Subquery Support** - IN, NOT IN, EXISTS, NOT EXISTS, and comparison subqueries
- ✅ **JOIN Operations** - INNER, LEFT, RIGHT, FULL OUTER, and CROSS joins with table aliases
- ✅ **Index Management** - CREATE INDEX and DROP INDEX with database-specific options
- ✅ **Table Operations** - CREATE, DROP, TRUNCATE with full dialect support
- ✅ **View Management** - CREATE and DROP views with IF EXISTS support
- ✅ **Multi-dialect Support** - MySQL, PostgreSQL, SQLite, and MSSQL
- ✅ **Transaction Support** - Built-in transaction management
- ✅ **Type Safety** - Strong typing for database operations


## Installation

```ssh
go get -u github.com/dracory/sb
```


## Example Create Table SQL

```go
import "github.com/dracory/sb"

sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("users").
	Column(Column{
		Name:       "id",
		Type:       COLUMN_TYPE_STRING,
		Length:     40,
		PrimaryKey: true,
	}).
	Column(Column{
		Name: "image",
		Type: COLUMN_TYPE_BLOB,
	}).
	Column(Column{
		Name: "price_default",
		Type: COLUMN_TYPE_DECIMAL,
	}).
	Column(Column{
		Name:     "price_custom",
		Type:     COLUMN_TYPE_DECIMAL,
		Length:   12,
		Decimals: 10,
	}).
	Column(Column{
		Name: "created_at",
		Type: COLUMN_TYPE_DATETIME,
	}).
	Column(Column{
		Name:     "deleted_at",
		Type:     COLUMN_TYPE_DATETIME,
		Nullable: true,
	}).
	Create()
```

## Example Table Drop SQL

```go
sql := NewBuilder(DIALECT_MYSQL).
	Table("users").
	Drop()
```


## Example Insert SQL

```go	
sql := sb.NewBuilder(DIALECT_MYSQL).
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

## Example Delete SQL

```go
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("user").
	Where(&sb.Where{
		Column: "id",
		Operator: "==",
		Value: "1",
	}).
	Limit(1).
	Delete()
```

## Example JOIN SQL

```go
// Basic INNER JOIN
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	Select([]string{"orders.*", "users.name"})

// LEFT JOIN with table alias
sql := sb.NewBuilder(DIALECT_POSTGRES).
	Table("orders").
	LeftJoin("profiles", "orders.user_id = profiles.user_id").
	Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
	Select([]string{"orders.*", "profiles.avatar"})

// Multiple JOINs
sql := sb.NewBuilder(DIALECT_SQLITE).
	Table("orders").
	InnerJoin("users", "orders.user_id = users.id").
	LeftJoin("profiles", "users.id = profiles.user_id").
	Select([]string{"orders.total", "users.name", "profiles.avatar"})

// RIGHT JOIN (MSSQL)
sql := sb.NewBuilder(DIALECT_MSSQL).
	Table("orders").
	RightJoin("users", "orders.user_id = users.id").
	Select([]string{"orders.*", "users.name"})

// FULL OUTER JOIN (PostgreSQL)
sql := sb.NewBuilder(DIALECT_POSTGRES).
	Table("orders").
	Join(sb.JOIN_FULL, "users", "orders.user_id = users.id").
	Select([]string{"orders.*", "users.name"})

// CROSS JOIN
sql := sb.NewBuilder(DIALECT_SQLITE).
	Table("orders").
	Join(sb.JOIN_CROSS, "users", "1=1").
	Select([]string{"orders.*", "users.name"})

// JOIN with table alias
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("orders").
	JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
	Select([]string{"orders.*", "p.avatar"})
```

## Example Subquery SQL

```go
// Basic IN Subquery
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("orders").
	Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	InSubquery(subquery).
	Select([]string{"name"})
// Result: SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > "1000");

// EXISTS Subquery
activeOrders := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("users").
	Exists(activeOrders).
	Select([]string{"name", "email"})
// Result: SELECT "name", "email" FROM "users" WHERE EXISTS (SELECT * FROM "orders" WHERE "status" = "active");

// NOT EXISTS Subquery
inactiveOrders := sb.NewBuilder(sb.DIALECT_SQLITE).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"})

sql := sb.NewBuilder(sb.DIALECT_SQLITE).
	Table("users").
	NotExists(inactiveOrders).
	Select([]string{"name"})
// Result: SELECT "name" FROM "users" WHERE NOT EXISTS (SELECT * FROM "orders" WHERE "status" = 'inactive');

// NOT IN Subquery
subquery := sb.NewBuilder(sb.DIALECT_MSSQL).
	Table("orders").
	Where(&sb.Where{Column: "status", Operator: "=", Value: "cancelled"})

sql := sb.NewBuilder(sb.DIALECT_MSSQL).
	Table("users").
	NotInSubquery(subquery).
	Select([]string{"name"})
// Result: SELECT [name] FROM [users] WHERE [id] NOT IN (SELECT * FROM [orders] WHERE [status] = 'cancelled');

// Correlated Subquery
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

// Complex Subquery with Multiple Conditions
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

// Using Subquery() Method
subqueryBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).Subquery()
subquery := subqueryBuilder.
	Table("orders").
	Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	InSubquery(subquery).
	Select([]string{"name"})
// Result: SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > "1000");
```

## Example Index Operations

```go
// Create Index
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("users").
	CreateIndex("idx_users_email", "email")

// Drop Index
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("users").
	DropIndex("idx_users_email")

// Drop Index with IF EXISTS (PostgreSQL, SQLite, MSSQL)
sql := sb.NewBuilder(DIALECT_POSTGRES).
	Table("users").
	DropIndexIfExists("idx_users_email")

// Drop Index with Schema (PostgreSQL)
sql := sb.NewBuilder(DIALECT_POSTGRES).
	Table("users").
	DropIndexWithSchema("idx_users_email", "public")
```

## Example Truncate Table

```go
// Basic Truncate
sql := sb.NewBuilder(DIALECT_MYSQL).
	Table("users").
	Truncate()

// Truncate with Options (PostgreSQL CASCADE)
sql := sb.NewBuilder(DIALECT_POSTGRES).
	Table("orders").
	TruncateWithOptions(sb.TruncateOptions{Cascade: true})

// Truncate with Options (MSSQL Reset Identity)
sql := sb.NewBuilder(DIALECT_MSSQL).
	Table("users").
	TruncateWithOptions(sb.TruncateOptions{ResetIdentity: true})
```

## Initiating Database Instance

1) From existing Go DB instance
```
myDb := sb.NewDatabase(sqlDb, sb.DIALECT_MYSQL)
```

3) From driver
```
myDb, err := sb.NewDatabaseFromDriver("sqlite3", "test.db")
```

## Example SQL Execute

```
myDb := sb.NewDatabase(sqlDb, sb.DIALECT_MYSQL)

ctx := context.Background()

_, err := myDb.Exec(ctx, sql)
```

## Example Transaction (using ExecInTransaction)

```go
import _ "modernc.org/sqlite"

myDb = sb.NewDatabaseFromDriver("sqlite3", "test.db")

ctx := context.Background()

err := myDb.ExecInTransaction(ctx, func(tx *sb.Database) error {
	if _, err := tx.Exec(ctx, sql1); err != nil {
		return err
	}

	if _, err := tx.Exec(ctx, sql2); err != nil {
		return err
	}

	return nil
})

if err != nil {
	// handle error
}
```

## Example Transaction (using BeginTransaction, CommitTransaction and RollbackTransaction)

```go
import _ "modernc.org/sqlite"

myDb = sb.NewDatabaseFromDriver("sqlite3", "test.db")

myDb.BeginTransaction()

ctx := context.Background()

_, err := myDb.Exec(ctx, sql1)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

_, err := myDb.Exec(ctx, sql2)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

myDb.CommitTransaction()

```

## Example Create View SQL

```go
selectSQL := sb.NewBuilder(DIALECT_POSTGRES).
	Table("users").
	Select([]string{"FirstName", "LastName"})

createViewSql := NewBuilder(DIALECT_POSTGRES).
	View("v_users").
	ViewColumns([]string{"first_name", "last_name"}).
	ViewSQL(selectSQL).
	Create()
```

## Example Create View If Not Exists SQL

```go
selectSQL := sb.NewBuilder(DIALECT_POSTGRES).
	Table("users").
	Select([]string{"FirstName", "LastName"})

createViewSql := NewBuilder(DIALECT_POSTGRES).
	View("v_users").
	ViewColumns([]string{"first_name", "last_name"}).
	ViewSQL(selectSQL).
	CreateIfNotExists()
```


## Example Drop View SQL

```go
dropiewSql := ab.NewBuilder(DIALECT_POSTGRES).
	View("v_users").
	Drop()
```


## Example Select as Map

Executes a select query and returns map[string]any

```go
ctx := context.Background()

mapAny := myDb.SelectToMapAny(ctx, sql)

```

Executes a select query and returns map[string]string

```go

ctx := context.Background()

mapString := myDb.SelectToMapString(ctx, sql)

```


## Developers

```sh
podman run -it --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=test -e MYSQL_DATABASE=test -e MYSQL_USER=test -e MYSQL_PASSWORD=test mysql:latest
```

```sh
podman run -it --rm -p 5432:5432 -e POSTGRES_PASSWORD=test -e POSTGRES_DB=test -e POSTGRES_USER=test postgres:latest
```



## Similar

- https://doug-martin.github.io/goqu - Best SQL Builder for Golang
- https://github.com/elgris/golang-sql-builder-benchmark
- https://github.com/es-code/gql
- https://github.com/cosiner/go-sqldb
- https://github.com/simukti/sqldb-logger
- https://github.com/elgs/gosqlcrud
- https://github.com/nandrechetan/gomb

## TODO
- github.com/stapelberg/postgrestest

## Recently Implemented Features
- ✅ **Subquery Support** - IN, NOT IN, EXISTS, NOT EXISTS, and comparison subqueries with correlation
- ✅ **JOIN Support** - INNER, LEFT, RIGHT, FULL OUTER, and CROSS joins with table aliases
- ✅ **Index Management** - Complete CREATE INDEX and DROP INDEX functionality
- ✅ **Table Truncation** - TRUNCATE TABLE support with database-specific options
- ✅ **Multi-dialect Support** - MySQL, PostgreSQL, SQLite, and MSSQL compatibility

