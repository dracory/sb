# SB <a href="https://gitpod.io/#https://github.com/dracory/sb" style="float:right:"><img src="https://gitpod.io/button/open-in-gitpod.svg" alt="Open in Gitpod" loading="lazy"></a>

![tests](https://github.com/dracory/sb/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/dracory/sb)](https://goreportcard.com/report/github.com/dracory/sb)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/dracory/sb)](https://pkg.go.dev/github.com/dracory/sb)

A simplified SQL builder (with limited functionality).

For a full SQL builder functionality check: https://doug-martin.github.io/goqu

Includes a wrapper for the mainstream DB package to allow transparent working with transactions.


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
	Where(sb.Where{
		Column: "id",
		Operator: "==",
		Value: "1",
	}).
	Limit(1).
	Delete()
```

## Initiating Database Instance

1) From existing Go DB instance
```
myDb := sb.NewDatabaseFromDb(sqlDb, DIALECT_MYSQL)
```

3) From driver
```
myDb = sql.NewDatabaseFromDriver("sqlite3", "test.db")
```

## Example SQL Execute

```
myDb := sb.NewDatabaseFromDb(sqlDb, DIALECT_MYSQL)
err := myDb.Exec(sql)
```

## Example Transaction

```go
import _ "github.com/glebarez/sqlite"

myDb = sb.NewDatabaseFromDriver("sqlite3", "test.db")

myDb.BeginTransaction()

err := Database.Exec(sql1)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

err := Database.Exec(sql2)

if err != nil {
	myDb.RollbackTransaction()
	return err
}

myDB.CommitTransaction()

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

mapAny := myDb.SelectToMapAny(sql)

```

Executes a select query and returns map[string]string

```go

mapString := myDb.SelectToMapAny(sql)

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

