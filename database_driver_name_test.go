package sb

import (
	"database/sql"
	"testing"

	_ "github.com/glebarez/sqlite"
)

func TestDatabaseDriverName(t *testing.T) {
	conn, err := sql.Open("sqlite", "test_newdatabase.db")
	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	name := DatabaseDriverName(conn)

	if name != "sqlite" {
		t.Fatal(`Error must be "sqlite" but got: `, name)
	}
}
