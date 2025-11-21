package sb_test

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/samber/lo"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dracory/database"
	"github.com/dracory/sb"

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

	err = sb.TableDropIfExists(database.Context(context.Background(), db), tableName)

	if err != nil {
		return nil, err
	}

	err = sb.TableCreate(context.Background(), db, tableName, columns)

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

	err = sb.TableDropIfExists(database.Context(context.Background(), db), tableName)

	if err != nil {
		return nil, err
	}

	err = sb.TableCreate(context.Background(), db, tableName, columns)

	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestTableColumnsMySQL(t *testing.T) {
	columns := _TestTableColumns_columns()

	db, err := initMySQLWithTable("test_table_columns", columns)

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnsMySQL test")
		return
	}

	defer db.Close()

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	columns, err = sb.TableColumns(database.Context(context.Background(), db), "test_table_columns", true)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if len(columns) != 9 {
		t.Fatal("Error must be 9 but got: ", len(columns))
	}

	expecteds := []struct {
		columnName   string
		columnType   string
		isPrimaryKey bool
		length       int
		decimals     int
		isNullable   bool
		Default      string
		isUnique     bool
	}{
		{"id", sb.COLUMN_TYPE_STRING, true, 40, 0, false, "", false},
		{"title", sb.COLUMN_TYPE_STRING, false, 100, 0, false, "", true},
		{"image", sb.COLUMN_TYPE_BLOB, false, 0, 0, false, "", false},
		{"price", sb.COLUMN_TYPE_DECIMAL, false, 0, 0, false, "", false},
		{"price_custom", sb.COLUMN_TYPE_DECIMAL, false, 12, 10, false, "", false},
		{"short_description", sb.COLUMN_TYPE_TEXT, false, 0, 0, false, "", false},
		{"long_description", sb.COLUMN_TYPE_TEXT, false, 0, 0, false, "", false},
		{"created_at", sb.COLUMN_TYPE_DATETIME, false, 0, 0, false, "", false},
		{"deleted_at", sb.COLUMN_TYPE_DATETIME, false, 0, 0, true, "", false},
	}

	for _, expected := range expecteds {
		column, found := lo.Find(columns, func(column sb.Column) bool {
			return column.Name == expected.columnName
		})

		if !found {
			t.Fatal("Error column '"+expected.columnName+"' must be found but got: ", found)
		}

		if column.Type != expected.columnType {
			t.Fatal("Error column '"+expected.columnName+"' type must be '"+expected.columnType+"' but got: ", column.Type)
		}
	}
}

func TestTableColumnsSQLite(t *testing.T) {
	columns := _TestTableColumns_columns()

	db, err := initSQLiteWithTable("test_table_columns", columns)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	defer db.Close()

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	columns, err = sb.TableColumns(database.Context(context.Background(), db), "test_table_columns", true)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if len(columns) != 9 {
		t.Fatal("Error must be 9 but got: ", len(columns))
	}

	expecteds := []struct {
		columnName   string
		columnType   string
		isPrimaryKey bool
		length       int
		decimals     int
		isNullable   bool
		Default      string
		isUnique     bool
	}{
		{"id", sb.COLUMN_TYPE_STRING, true, 40, 0, false, "", false},
		{"title", sb.COLUMN_TYPE_STRING, false, 100, 0, false, "", true},
		{"image", sb.COLUMN_TYPE_BLOB, false, 0, 0, false, "", false},
		{"price", sb.COLUMN_TYPE_DECIMAL, false, 0, 0, false, "", false},
		{"price_custom", sb.COLUMN_TYPE_DECIMAL, false, 12, 10, false, "", false},
		{"short_description", sb.COLUMN_TYPE_TEXT, false, 0, 0, false, "", false},
		{"long_description", sb.COLUMN_TYPE_TEXT, false, 0, 0, false, "", false},
		{"created_at", sb.COLUMN_TYPE_DATETIME, false, 0, 0, false, "", false},
		{"deleted_at", sb.COLUMN_TYPE_DATETIME, false, 0, 0, true, "", false},
	}

	for _, expected := range expecteds {
		column, found := lo.Find(columns, func(column sb.Column) bool {
			return column.Name == expected.columnName
		})

		if !found {
			t.Fatal("Error column '"+expected.columnName+"' must be found but got: ", found)
		}

		if column.Type != expected.columnType {
			t.Fatal("Error column '"+expected.columnName+"' type must be '"+expected.columnType+"' but got: ", column.Type)
		}
	}
}

func _TestTableColumns_columns() []sb.Column {
	columns := []sb.Column{
		{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		},
		{
			Name:   "title",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 100,
			Unique: true,
		},
		{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		},
		{
			Name: "price",
			Type: sb.COLUMN_TYPE_DECIMAL,
		},
		{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		},
		{
			Name: "short_description",
			Type: sb.COLUMN_TYPE_TEXT,
		},
		{
			Name: "long_description",
			Type: sb.COLUMN_TYPE_TEXT,
		},
		{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		},
		{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		},
	}

	return columns
}
