package indexops

import (
	"testing"
)

func TestCreateSimpleIndex(t *testing.T) {
	sql, err := CreateSimpleIndex()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestCreateUniqueIndex(t *testing.T) {
	sql, err := CreateUniqueIndex()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestCreateCompositeIndex(t *testing.T) {
	sql, err := CreateCompositeIndex()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestCreatePartialIndex(t *testing.T) {
	sql, err := CreatePartialIndex()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestDropIndex(t *testing.T) {
	sql, err := DropIndex()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestDropIndexIfExists(t *testing.T) {
	sql, err := DropIndexIfExists()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}

func TestDropIndexWithSchema(t *testing.T) {
	sql, err := DropIndexWithSchema()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
}
