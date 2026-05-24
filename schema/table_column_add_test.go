package schema_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
)

func TestTableColumnAddSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initSQLiteWithTable("test_column_add", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	newCol := sb.Column{Name: "new_col", Type: sb.COLUMN_TYPE_STRING, Length: 50}

	exists, err := schema.TableColumnExists(ctx, "test_column_add", "new_col")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist before add")
	}

	err = schema.TableColumnAdd(ctx, "test_column_add", newCol)
	if err != nil {
		t.Fatal("Error adding column: ", err.Error())
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_add", "new_col")
	if err != nil {
		t.Fatal("Error checking column existence after add: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist after add")
	}
}

func TestTableColumnAddMySQL(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnAddMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initMySQLWithTable("test_column_add", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	newCol := sb.Column{Name: "new_col", Type: sb.COLUMN_TYPE_STRING, Length: 50}

	err = schema.TableColumnAdd(ctx, "test_column_add", newCol)
	if err != nil {
		t.Fatal("Error adding column: ", err.Error())
	}

	exists, err := schema.TableColumnExists(ctx, "test_column_add", "new_col")
	if err != nil {
		t.Fatal("Error checking column existence after add: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist after add")
	}
}

func TestTableColumnAddIfNotExistsSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "existing_col", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initSQLiteWithTable("test_column_add_if_not_exists", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	existingCol := sb.Column{Name: "existing_col", Type: sb.COLUMN_TYPE_STRING, Length: 100}
	err = schema.TableColumnAddIfNotExists(ctx, "test_column_add_if_not_exists", existingCol)
	if err != nil {
		t.Fatal("Adding existing column with IfNotExists should not error, got: ", err.Error())
	}

	newCol := sb.Column{Name: "brand_new_col", Type: sb.COLUMN_TYPE_STRING, Length: 50}
	err = schema.TableColumnAddIfNotExists(ctx, "test_column_add_if_not_exists", newCol)
	if err != nil {
		t.Fatal("Error adding new column: ", err.Error())
	}

	exists, err := schema.TableColumnExists(ctx, "test_column_add_if_not_exists", "brand_new_col")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("New column should exist after add")
	}
}
