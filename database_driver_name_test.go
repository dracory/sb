package sb_test

import (
	"database/sql"
	"testing"

	"github.com/dracory/sb"
	_ "github.com/glebarez/sqlite"
)

func TestDatabaseDriverName(t *testing.T) {
	conn, err := sql.Open("sqlite", "test_newdatabase.db")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	name := sb.DatabaseDriverName(conn)

	if name != "sqlite" {
		t.Fatal(`Error must be "sqlite" but got: `, name)
	}
}
