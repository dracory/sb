package schema_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
)

func TestTableColumnRenameSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initSQLiteWithTable("test_column_rename", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	exists, err := schema.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Old column should exist before rename")
	}

	err = schema.TableColumnRename(ctx, "test_column_rename", "email", "new_email")
	if err != nil {
		t.Fatal("Error renaming column: ", err.Error())
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking old column existence after rename: ", err.Error())
	}
	if exists {
		t.Fatal("Old column should not exist after rename")
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_rename", "new_email")
	if err != nil {
		t.Fatal("Error checking new column existence after rename: ", err.Error())
	}
	if !exists {
		t.Fatal("New column should exist after rename")
	}
}

func TestTableColumnRenameMySQL(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnRenameMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initMySQLWithTable("test_column_rename", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	exists, err := schema.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Old column should exist before rename")
	}

	err = schema.TableColumnRename(ctx, "test_column_rename", "email", "new_email")
	if err != nil {
		t.Fatal("Error renaming column: ", err.Error())
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_rename", "email")
	if err != nil {
		t.Fatal("Error checking old column existence after rename: ", err.Error())
	}
	if exists {
		t.Fatal("Old column should not exist after rename")
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_rename", "new_email")
	if err != nil {
		t.Fatal("Error checking new column existence after rename: ", err.Error())
	}
	if !exists {
		t.Fatal("New column should exist after rename")
	}
}
