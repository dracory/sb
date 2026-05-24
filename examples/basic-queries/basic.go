package basic

import (
	"github.com/dracory/sb"
)

// SimpleSelect returns a basic SELECT query for active users
func SimpleSelect() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		Select([]string{"id", "name", "email"})
}

// SelectWithOrderBy returns a SELECT query with ORDER BY clause
func SelectWithOrderBy() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		OrderBy("created_at", "DESC").
		Select([]string{"*"})
}

// SelectWithLimit returns a SELECT query with LIMIT clause
func SelectWithLimit() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		OrderBy("created_at", "DESC").
		Limit(10).
		Select([]string{"*"})
}

// SimpleInsert returns a basic INSERT statement
func SimpleInsert() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Insert(map[string]string{
			"name":  "John Doe",
			"email": "john@example.com",
		})
}

// InsertWithMultipleColumns returns an INSERT with multiple columns
func InsertWithMultipleColumns() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Insert(map[string]string{
			"name":       "Jane Doe",
			"email":      "jane@example.com",
			"status":     "active",
			"created_at": "2024-01-01",
		})
}

// SimpleUpdate returns a basic UPDATE statement
func SimpleUpdate() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
		Update(map[string]string{
			"status": "inactive",
		})
}

// UpdateWithMultipleColumns returns an UPDATE with multiple columns
func UpdateWithMultipleColumns() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
		Update(map[string]string{
			"status":     "active",
			"updated_at": "2024-01-01",
		})
}

// SimpleDelete returns a basic DELETE statement
func SimpleDelete() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "1"}).
		Delete()
}

// DeleteWithMultipleConditions returns a DELETE with multiple WHERE conditions
func DeleteWithMultipleConditions() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"}).
		Where(&sb.Where{Column: "created_at", Operator: "<", Value: "2024-01-01"}).
		Delete()
}
