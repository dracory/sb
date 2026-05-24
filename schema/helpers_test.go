package schema_test

import (
	"context"
	"database/sql"
	"os"
	"strings"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
	"github.com/samber/lo"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

var TestsWithMySQL = true

func initMySQLWithTable(tableName string, columns []sb.Column) (db *sql.DB, err error) {
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbUser := os.Getenv("MYSQL_USER")
	dbPass := os.Getenv("MYSQL_PASS")
	dbName := os.Getenv("MYSQL_DATABASE")

	host = lo.Ternary(host == "", "localhost", host)
	port = lo.Ternary(port == "", "33306", port)
	dbUser = lo.Ternary(dbUser == "", "test", dbUser)
	dbPass = lo.Ternary(dbPass == "", "test", dbPass)
	dbName = lo.Ternary(dbName == "", "test", dbName)

	db, err = database.Open(database.Options().
		SetDatabaseType(database.DATABASE_TYPE_MYSQL).
		SetDatabaseHost(host).
		SetDatabasePort(port).
		SetDatabaseName(dbName).
		SetUserName(dbUser).
		SetPassword(dbPass))

	if err != nil {
		if strings.Contains(err.Error(), "could not be pinge") {
			TestsWithMySQL = false
		}
		return nil, err
	}

	err = schema.TableDropIfExists(database.Context(context.Background(), db), tableName)
	if err != nil {
		return nil, err
	}

	err = schema.TableCreate(context.Background(), db, tableName, columns)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initSQLiteWithTable(tableName string, columns []sb.Column) (db *sql.DB, err error) {
	db, err = database.Open(database.Options().
		SetDatabaseType(database.DATABASE_TYPE_SQLITE).
		SetDatabaseName(":memory:"))

	if err != nil {
		return nil, err
	}

	err = schema.TableDropIfExists(database.Context(context.Background(), db), tableName)
	if err != nil {
		return nil, err
	}

	err = schema.TableCreate(context.Background(), db, tableName, columns)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func testTableColumns_columns() []sb.Column {
	return []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "title", Type: sb.COLUMN_TYPE_STRING, Length: 100, Unique: true},
		{Name: "image", Type: sb.COLUMN_TYPE_BLOB},
		{Name: "price", Type: sb.COLUMN_TYPE_DECIMAL},
		{Name: "price_custom", Type: sb.COLUMN_TYPE_DECIMAL, Length: 12, Decimals: 10},
		{Name: "short_description", Type: sb.COLUMN_TYPE_TEXT},
		{Name: "long_description", Type: sb.COLUMN_TYPE_TEXT},
		{Name: "created_at", Type: sb.COLUMN_TYPE_DATETIME},
		{Name: "deleted_at", Type: sb.COLUMN_TYPE_DATETIME, Nullable: true},
	}
}
