package subqueryops

import (
	"testing"
)

func TestInSubquery(t *testing.T) {
	sql, params, err := InSubquery()
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

func TestNotInSubquery(t *testing.T) {
	sql, params, err := NotInSubquery()
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

func TestExistsSubquery(t *testing.T) {
	sql, params, err := ExistsSubquery()
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

func TestNotExistsSubquery(t *testing.T) {
	sql, params, err := NotExistsSubquery()
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

func TestCorrelatedSubquery(t *testing.T) {
	sql, params, err := CorrelatedSubquery()
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

func TestSubqueryInComparison(t *testing.T) {
	sql, params, err := SubqueryInComparison()
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
