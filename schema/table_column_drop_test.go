package schema_test

import (
	"context"
	"os"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
)

func TestTableColumnDropMySQL(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnDropMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "temp_column", Type: sb.COLUMN_TYPE_STRING, Length: 50},
	}

	db, err := initMySQLWithTable("test_column_drop", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	exists, err := schema.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist before drop")
	}

	err = schema.TableColumnDrop(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error dropping column: ", err.Error())
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence after drop: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist after drop")
	}
}

func TestTableColumnDropSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "temp_column", Type: sb.COLUMN_TYPE_STRING, Length: 50},
	}

	db, err := initSQLiteWithTable("test_column_drop", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	exists, err := schema.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist before drop")
	}

	err = schema.TableColumnDrop(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error dropping column: ", err.Error())
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_drop", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence after drop: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist after drop")
	}
}

func TestTableColumnDropIfExistsMySQL(t *testing.T) {
	if os.Getenv("GITHUB_ACTIONS") != "true" {
		t.Skip("MySQL integration tests only run in GitHub Actions")
		return
	}

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnDropIfExistsMySQL test")
		return
	}

	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initMySQLWithTable("test_column_drop_if_exists", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	err = schema.TableColumnDropIfExists(ctx, "test_column_drop_if_exists", "nonexistent_column")
	if err != nil {
		t.Fatal("Dropping non-existent column should not error, got: ", err.Error())
	}
}

func TestTableColumnDropIfExistsSQLite(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
	}

	db, err := initSQLiteWithTable("test_column_drop_if_exists", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	err = schema.TableColumnDropIfExists(ctx, "test_column_drop_if_exists", "nonexistent_column")
	if err != nil {
		t.Fatal("Dropping non-existent column should not error, got: ", err.Error())
	}
}

func TestTableColumnDropIfExistsWithActualColumn(t *testing.T) {
	columns := []sb.Column{
		{Name: "id", Type: sb.COLUMN_TYPE_STRING, Length: 40, PrimaryKey: true},
		{Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 100},
		{Name: "temp_column", Type: sb.COLUMN_TYPE_STRING, Length: 50},
	}

	db, err := initSQLiteWithTable("test_column_drop_if_exists_real", columns)
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	defer func() { _ = db.Close() }()

	ctx := database.Context(context.Background(), db)

	exists, err := schema.TableColumnExists(ctx, "test_column_drop_if_exists_real", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence: ", err.Error())
	}
	if !exists {
		t.Fatal("Column should exist before drop")
	}

	err = schema.TableColumnDropIfExists(ctx, "test_column_drop_if_exists_real", "temp_column")
	if err != nil {
		t.Fatal("Error dropping column: ", err.Error())
	}

	exists, err = schema.TableColumnExists(ctx, "test_column_drop_if_exists_real", "temp_column")
	if err != nil {
		t.Fatal("Error checking column existence after drop: ", err.Error())
	}
	if exists {
		t.Fatal("Column should not exist after drop")
	}
}
