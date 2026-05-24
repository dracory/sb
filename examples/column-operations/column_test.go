package columnops

import (
	"testing"
)

func TestAddColumn(t *testing.T) {
	sql, err := AddColumn()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestDropColumn(t *testing.T) {
	sql, err := DropColumn()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestRenameColumn(t *testing.T) {
	sql, err := RenameColumn()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}
