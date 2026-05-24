package schema_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
)

func TestTableColumnExistsMySQL(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	columns := testTableColumns_columns()
	db, err := initMySQLWithTable("test_table_columns", columns)

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnExistsMySQL test")
		return
	}

	defer func() { _ = db.Close() }()

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	exists, err := schema.TableColumnExists(database.Context(context.Background(), db), "test_table_columns", "id")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if !exists {
		t.Fatal("Error must be true but got: ", exists)
	}
}

func TestTableColumnsExistsSQLite(t *testing.T) {
	columns := testTableColumns_columns()
	db, err := initSQLiteWithTable("test_table_columns", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	defer func() { _ = db.Close() }()

	exists, err := schema.TableColumnExists(database.Context(context.Background(), db), "test_table_columns", "id")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if !exists {
		t.Fatal("Error must be true but got: ", exists)
	}
}

func TestTableColumnExistsNilQueryable(t *testing.T) {
	ctx := database.Context(context.Background(), nil)
	_, err := schema.TableColumnExists(ctx, "users", "id")
	if err == nil {
		t.Fatal("Expected error for nil queryable")
	}
}

func TestTableColumnExistsEmptyNames(t *testing.T) {
	db, err := initSQLiteWithTable("test_empty", []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
	})
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	_, err = schema.TableColumnExists(ctx, "", "id")
	if err == nil {
		t.Fatal("Expected error for empty table name")
	}

	_, err = schema.TableColumnExists(ctx, "test_empty", "")
	if err == nil {
		t.Fatal("Expected error for empty column name")
	}
}
