package basic

import (
	"testing"
)

func TestSimpleSelect(t *testing.T) {
	sql, params, err := SimpleSelect()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(params))
	}
}

func TestSelectWithOrderBy(t *testing.T) {
	sql, params, err := SelectWithOrderBy()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(params))
	}
}

func TestSelectWithLimit(t *testing.T) {
	sql, params, err := SelectWithLimit()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	// No parameters expected for ORDER BY + LIMIT
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}

func TestSimpleInsert(t *testing.T) {
	sql, params, err := SimpleInsert()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(params))
	}
}

func TestInsertWithMultipleColumns(t *testing.T) {
	sql, params, err := InsertWithMultipleColumns()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 4 {
		t.Fatalf("Expected 4 parameters, got %d", len(params))
	}
}

func TestSimpleUpdate(t *testing.T) {
	sql, params, err := SimpleUpdate()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(params))
	}
}

func TestUpdateWithMultipleColumns(t *testing.T) {
	sql, params, err := UpdateWithMultipleColumns()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 3 {
		t.Fatalf("Expected 3 parameters, got %d", len(params))
	}
}

func TestSimpleDelete(t *testing.T) {
	sql, params, err := SimpleDelete()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 1 {
		t.Fatalf("Expected 1 parameter, got %d", len(params))
	}
}

func TestDeleteWithMultipleConditions(t *testing.T) {
	sql, params, err := DeleteWithMultipleConditions()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 2 {
		t.Fatalf("Expected 2 parameters, got %d", len(params))
	}
}
