package joinops

import (
	"testing"
)

func TestInnerJoin(t *testing.T) {
	sql, params, err := InnerJoin()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}

func TestLeftJoin(t *testing.T) {
	sql, params, err := LeftJoin()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}

func TestRightJoin(t *testing.T) {
	sql, params, err := RightJoin()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}

func TestJoinWithAlias(t *testing.T) {
	sql, params, err := JoinWithAlias()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}

func TestMultipleJoins(t *testing.T) {
	sql, params, err := MultipleJoins()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}

func TestJoinWithWhere(t *testing.T) {
	sql, params, err := JoinWithWhere()
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

func TestJoinWithOrderBy(t *testing.T) {
	sql, params, err := JoinWithOrderBy()
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}
	if sql == "" {
		t.Fatal("SQL should not be empty")
	}
	if len(params) != 0 {
		t.Fatalf("Expected 0 parameters, got %d", len(params))
	}
}
