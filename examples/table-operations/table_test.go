package tableops

import (
	"testing"
)

func TestCreateTable(t *testing.T) {
	sql, err := CreateTable()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestDropTable(t *testing.T) {
	sql, err := DropTable()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestDropTableIfExists(t *testing.T) {
	sql, err := DropTableIfExists()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestTruncateTable(t *testing.T) {
	sql, err := TruncateTable()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestTruncateTableWithCascade(t *testing.T) {
	sql, err := TruncateTableWithCascade()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestTruncateTableWithResetIdentity(t *testing.T) {
	sql, err := TruncateTableWithResetIdentity()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}
