package joinops

import (
	"github.com/dracory/sb"
)

// InnerJoin returns a SELECT with INNER JOIN
func InnerJoin() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})
}

// LeftJoin returns a SELECT with LEFT JOIN
func LeftJoin() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		LeftJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})
}

// RightJoin returns a SELECT with RIGHT JOIN
func RightJoin() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		RightJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})
}

// JoinWithAlias returns a SELECT with JOIN using table alias
func JoinWithAlias() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
		Select([]string{"orders.*", "p.avatar"})
}

// MultipleJoins returns a SELECT with multiple JOINs
func MultipleJoins() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		LeftJoin("profiles", "users.id = profiles.user_id").
		Select([]string{"orders.total", "users.name", "profiles.avatar"})
}

// JoinWithWhere returns a SELECT with JOIN and WHERE clause
func JoinWithWhere() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		LeftJoin("users", "orders.user_id = users.id").
		Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
		Select([]string{"orders.*", "users.name"})
}

// JoinWithOrderBy returns a SELECT with JOIN and ORDER BY
func JoinWithOrderBy() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		OrderBy("orders.created_at", "DESC").
		Select([]string{"orders.*", "users.name"})
}
