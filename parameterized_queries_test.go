package sb_test

import (
	"testing"

	"github.com/dracory/sb"
)

// TestParameterizedQueriesMySQL tests parameterized query generation for MySQL
func TestParameterizedQueriesMySQL(t *testing.T) {
	// Test SELECT with WHERE clause
	sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "123"}).
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		Select([]string{"name", "email"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT `name`, `email` FROM `users` WHERE `id` = ? AND `status` = ?;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}

	if params[0] != "123" || params[1] != "active" {
		t.Errorf("Expected params [123, active], got %v", params)
	}
}

// TestParameterizedQueriesPostgreSQL tests parameterized query generation for PostgreSQL
func TestParameterizedQueriesPostgreSQL(t *testing.T) {
	// Test SELECT with WHERE clause
	sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{Column: "email", Operator: "=", Value: "user@example.com"}).
		Where(&sb.Where{Column: "age", Operator: ">", Value: "18"}).
		Select([]string{"*"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT * FROM \"users\" WHERE \"email\" = $1 AND \"age\" > $2;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}

	if params[0] != "user@example.com" || params[1] != "18" {
		t.Errorf("Expected params [user@example.com, 18], got %v", params)
	}
}

// TestParameterizedQueriesSQLite tests parameterized query generation for SQLite
func TestParameterizedQueriesSQLite(t *testing.T) {
	// Test INSERT
	sql, params, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Insert(map[string]string{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   "30",
		})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "INSERT INTO \"users\" (\"age\", \"email\", \"name\") VALUES (?, ?, ?);"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}

	// Parameters should be in alphabetical order by column name
	if params[0] != "30" || params[1] != "john@example.com" || params[2] != "John Doe" {
		t.Errorf("Expected params [30, john@example.com, John Doe], got %v", params)
	}
}

// TestParameterizedQueriesMSSQL tests parameterized query generation for MSSQL
func TestParameterizedQueriesMSSQL(t *testing.T) {
	// Test UPDATE
	sql, params, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
		Update(map[string]string{
			"name":   "Jane Doe",
			"status": "inactive",
		})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "UPDATE [users] SET [name]=@p2, [status]=@p3 WHERE [id] = @p1;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}

	// WHERE params are collected first, then UPDATE params (alphabetical)
	if params[0] != "1" || params[1] != "Jane Doe" || params[2] != "inactive" {
		t.Errorf("Expected params [1, Jane Doe, inactive], got %v", params)
	}
}

// TestParameterizedQueriesDelete tests DELETE with parameters
func TestParameterizedQueriesDelete(t *testing.T) {
	sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "deleted"}).
		Where(&sb.Where{Column: "created_at", Operator: "<", Value: "2020-01-01"}).
		Delete()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "DELETE FROM `users` WHERE `status` = ? AND `created_at` < ?;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 2 {
		t.Errorf("Expected 2 parameters, got %d", len(params))
	}

	if params[0] != "deleted" || params[1] != "2020-01-01" {
		t.Errorf("Expected params [deleted, 2020-01-01], got %v", params)
	}
}

// TestParameterizedQueriesNullHandling tests NULL value handling
func TestParameterizedQueriesNullHandling(t *testing.T) {
	sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "deleted_at", Operator: "=", Value: "NULL"}).
		Select([]string{"*"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// NULL should not be parameterized
	expectedSQL := "SELECT * FROM `users` WHERE `deleted_at` IS NULL;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 0 {
		t.Errorf("Expected 0 parameters for NULL, got %d", len(params))
	}
}

// TestInterpolatedValuesMode tests backward compatibility mode
func TestInterpolatedValuesMode(t *testing.T) {
	sql, params, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "123"}).
		WithInterpolatedValues().
		Select([]string{"name"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// In interpolated mode, values should be quoted in SQL
	expectedSQL := "SELECT `name` FROM `users` WHERE `id` = \"123\";"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	// No parameters in interpolated mode
	if len(params) != 0 {
		t.Errorf("Expected 0 parameters in interpolated mode, got %d", len(params))
	}
}

// TestParameterizedQueriesComplexWhere tests complex WHERE conditions
func TestParameterizedQueriesComplexWhere(t *testing.T) {
	sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "pending"}).
		Where(&sb.Where{Column: "total", Operator: ">", Value: "100"}).
		Where(&sb.Where{Column: "created_at", Operator: ">=", Value: "2024-01-01"}).
		OrderBy("created_at", "DESC").
		Limit(10).
		Select([]string{"id", "total", "status"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT \"id\", \"total\", \"status\" FROM \"orders\" WHERE \"status\" = $1 AND \"total\" > $2 AND \"created_at\" >= $3 ORDER BY \"created_at\" DESC LIMIT 10;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 3 {
		t.Errorf("Expected 3 parameters, got %d", len(params))
	}

	if params[0] != "pending" || params[1] != "100" || params[2] != "2024-01-01" {
		t.Errorf("Expected params [pending, 100, 2024-01-01], got %v", params)
	}
}

// TestParameterizedQueriesAllDialects tests all dialects produce correct placeholders
func TestParameterizedQueriesAllDialects(t *testing.T) {
	tests := []struct {
		dialect        string
		expectedSQL    string
		expectedParams []interface{}
	}{
		{
			dialect:        sb.DIALECT_MYSQL,
			expectedSQL:    "SELECT * FROM `users` WHERE `id` = ?;",
			expectedParams: []interface{}{"1"},
		},
		{
			dialect:        sb.DIALECT_POSTGRES,
			expectedSQL:    "SELECT * FROM \"users\" WHERE \"id\" = $1;",
			expectedParams: []interface{}{"1"},
		},
		{
			dialect:        sb.DIALECT_SQLITE,
			expectedSQL:    "SELECT * FROM \"users\" WHERE \"id\" = ?;",
			expectedParams: []interface{}{"1"},
		},
		{
			dialect:        sb.DIALECT_MSSQL,
			expectedSQL:    "SELECT * FROM [users] WHERE [id] = @p1;",
			expectedParams: []interface{}{"1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.dialect, func(t *testing.T) {
			sql, params, err := sb.NewBuilder(tt.dialect).
				Table("users").
				Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
				Select([]string{"*"})

			if err != nil {
				t.Fatalf("Unexpected error for %s: %v", tt.dialect, err)
			}

			if sql != tt.expectedSQL {
				t.Errorf("Dialect %s:\nExpected SQL: %s\nGot SQL: %s", tt.dialect, tt.expectedSQL, sql)
			}

			if len(params) != len(tt.expectedParams) {
				t.Errorf("Dialect %s: Expected %d params, got %d", tt.dialect, len(tt.expectedParams), len(params))
			}

			for i, expected := range tt.expectedParams {
				if params[i] != expected {
					t.Errorf("Dialect %s: Expected param[%d]=%v, got %v", tt.dialect, i, expected, params[i])
				}
			}
		})
	}
}

// TestParameterizedQueriesMultipleParameters tests queries with many parameters
func TestParameterizedQueriesMultipleParameters(t *testing.T) {
	sql, params, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("products").
		Where(&sb.Where{Column: "category", Operator: "=", Value: "electronics"}).
		Where(&sb.Where{Column: "price", Operator: ">=", Value: "100"}).
		Where(&sb.Where{Column: "price", Operator: "<=", Value: "500"}).
		Where(&sb.Where{Column: "in_stock", Operator: "=", Value: "true"}).
		Where(&sb.Where{Column: "brand", Operator: "=", Value: "Apple"}).
		Select([]string{"name", "price"})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	expectedSQL := "SELECT \"name\", \"price\" FROM \"products\" WHERE \"category\" = $1 AND \"price\" >= $2 AND \"price\" <= $3 AND \"in_stock\" = $4 AND \"brand\" = $5;"
	if sql != expectedSQL {
		t.Errorf("Expected SQL:\n%s\nGot:\n%s", expectedSQL, sql)
	}

	if len(params) != 5 {
		t.Errorf("Expected 5 parameters, got %d", len(params))
	}

	expectedParams := []interface{}{"electronics", "100", "500", "true", "Apple"}
	for i, expected := range expectedParams {
		if params[i] != expected {
			t.Errorf("Expected param[%d]=%v, got %v", i, expected, params[i])
		}
	}
}
