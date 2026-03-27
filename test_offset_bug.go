package sb

import (
	"fmt"
	"testing"
)

func TestMSSQLOffsetBug(t *testing.T) {
	// Test MSSQL with OFFSET - this should demonstrate the bug
	fmt.Println("Testing MSSQL OFFSET bug...")

	// This should generate SQL but might cause the OFFSET syntax error
	sql, _, err := NewBuilder(DIALECT_MSSQL).
		Table("users").
		OrderBy("id", "ASC").
		Limit(10).
		Offset(20).
		Select([]string{"id", "name"})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	fmt.Printf("Generated SQL: %s\n", sql)

	// Test with SQLite to see the difference
	fmt.Println("Testing SQLite OFFSET (should work)...")
	sql2, _, err := NewBuilder(DIALECT_SQLITE).
		Table("users").
		OrderBy("id", "ASC").
		Limit(10).
		Offset(20).
		Select([]string{"id", "name"})

	if err != nil {
		t.Fatalf("Error: %v", err)
	}

	fmt.Printf("Generated SQL: %s\n", sql2)
}
