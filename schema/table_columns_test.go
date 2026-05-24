package schema_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
	"github.com/samber/lo"
)

func TestTableColumnsMySQL(t *testing.T) {
	columns := testTableColumns_columns()
	db, err := initMySQLWithTable("test_table_columns", columns)

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnsMySQL test")
		return
	}

	defer func() { _ = db.Close() }()

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	columns, err = schema.TableColumns(database.Context(context.Background(), db), "test_table_columns", true)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if len(columns) != 9 {
		t.Fatal("Error must be 9 but got: ", len(columns))
	}

	expecteds := []struct {
		columnName string
		columnType string
	}{
		{"id", sb.COLUMN_TYPE_STRING},
		{"title", sb.COLUMN_TYPE_STRING},
		{"image", sb.COLUMN_TYPE_BLOB},
		{"price", sb.COLUMN_TYPE_DECIMAL},
		{"price_custom", sb.COLUMN_TYPE_DECIMAL},
		{"short_description", sb.COLUMN_TYPE_TEXT},
		{"long_description", sb.COLUMN_TYPE_TEXT},
		{"created_at", sb.COLUMN_TYPE_DATETIME},
		{"deleted_at", sb.COLUMN_TYPE_DATETIME},
	}

	for _, expected := range expecteds {
		column, found := lo.Find(columns, func(c sb.Column) bool {
			return c.Name == expected.columnName
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
	columns := testTableColumns_columns()
	db, err := initSQLiteWithTable("test_table_columns", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	defer func() { _ = db.Close() }()

	columns, err = schema.TableColumns(database.Context(context.Background(), db), "test_table_columns", true)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if len(columns) != 9 {
		t.Fatal("Error must be 9 but got: ", len(columns))
	}

	expecteds := []struct {
		columnName string
		columnType string
	}{
		{"id", sb.COLUMN_TYPE_STRING},
		{"title", sb.COLUMN_TYPE_STRING},
		{"image", sb.COLUMN_TYPE_BLOB},
		{"price", sb.COLUMN_TYPE_DECIMAL},
		{"price_custom", sb.COLUMN_TYPE_DECIMAL},
		{"short_description", sb.COLUMN_TYPE_TEXT},
		{"long_description", sb.COLUMN_TYPE_TEXT},
		{"created_at", sb.COLUMN_TYPE_DATETIME},
		{"deleted_at", sb.COLUMN_TYPE_DATETIME},
	}

	for _, expected := range expecteds {
		column, found := lo.Find(columns, func(c sb.Column) bool {
			return c.Name == expected.columnName
		})

		if !found {
			t.Fatal("Error column '"+expected.columnName+"' must be found but got: ", found)
		}

		if column.Type != expected.columnType {
			t.Fatal("Error column '"+expected.columnName+"' type must be '"+expected.columnType+"' but got: ", column.Type)
		}
	}
}
