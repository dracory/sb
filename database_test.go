package sb

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/glebarez/sqlite"
)

func TestNewDatabase(t *testing.T) {
	conn, err := sql.Open("sqlite", "test_new_database.db")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}
	db := NewDatabase(conn, DIALECT_SQLITE)
	if db == nil {
		t.Fatal("Database MUST NOT BE NIL")
	}
	if db.DB() == nil {
		t.Fatal("Database db field MUST NOT BE NIL")
	}
	if db.Tx() != nil {
		t.Fatal("Database tx field MUST BE NIL")
	}
}

func TestDatabase_SqlLog(t *testing.T) {
	db := NewDatabase(nil, DIALECT_SQLITE)

	if !reflect.DeepEqual(db.SqlLog(), []map[string]string{}) {
		t.Fatal("SqlLog must be empty")
	}
}

func TestDatabase_SqlLogEmpty(t *testing.T) {
	db := NewDatabase(nil, DIALECT_SQLITE)

	if len(db.SqlLog()) != 0 {
		t.Fatal("SqlLog must be empty")
	}

	db.SqlLogEmpty()

	if len(db.SqlLog()) != 0 {
		t.Fatal("SqlLog must be empty")
	}
}
