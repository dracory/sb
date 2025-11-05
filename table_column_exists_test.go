package sb_test

import (
	"context"
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

func TestTableColumnExistsMySQL(t *testing.T) {
	columns := _TestTableColumns_columns()

	db, err := initMySQLWithTable("test_table_columns", columns)

	if TestsWithMySQL == false {
		t.Log("TestsWithMySQL is false. Skipping TestTableColumnExistsMySQL test")
		return
	}

	defer db.Close()

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	exists, err := sb.TableColumnExists(database.Context(context.Background(), db), "test_table_columns", "id")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if !exists {
		t.Fatal("Error must be true but got: ", exists)
	}
}

func TestTableColumnsExistsSQLite(t *testing.T) {
	columns := _TestTableColumns_columns()

	db, err := initSQLiteWithTable("test_table_columns", columns)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	defer db.Close()

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	exists, err := sb.TableColumnExists(database.Context(context.Background(), db), "test_table_columns", "id")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if !exists {
		t.Fatal("Error must be true but got: ", exists)
	}
}
