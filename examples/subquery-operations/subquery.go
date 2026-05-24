package subqueryops

import (
	"github.com/dracory/sb"
)

// InSubquery returns a SELECT with IN subquery
func InSubquery() (string, []interface{}, error) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(subquery)
	if err != nil {
		return "", nil, err
	}

	return builder.Select([]string{"name"})
}

// NotInSubquery returns a SELECT with NOT IN subquery
func NotInSubquery() (string, []interface{}, error) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "cancelled"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		NotInSubquery(subquery)
	if err != nil {
		return "", nil, err
	}

	return builder.Select([]string{"name"})
}

// ExistsSubquery returns a SELECT with EXISTS subquery
func ExistsSubquery() (string, []interface{}, error) {
	activeOrders := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

	builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Exists(activeOrders)
	if err != nil {
		return "", nil, err
	}

	return builder.Select([]string{"name"})
}

// NotExistsSubquery returns a SELECT with NOT EXISTS subquery
func NotExistsSubquery() (string, []interface{}, error) {
	activeOrders := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

	builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		NotExists(activeOrders)
	if err != nil {
		return "", nil, err
	}

	return builder.Select([]string{"name"})
}

// CorrelatedSubquery returns a SELECT with correlated subquery
func CorrelatedSubquery() (string, []interface{}, error) {
	subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"}).
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

	builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		NotExists(subquery)
	if err != nil {
		return "", nil, err
	}

	return builder.Select([]string{"name", "email"})
}

// SubqueryInComparison returns a SELECT with subquery in comparison
func SubqueryInComparison() (string, []interface{}, error) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"})

	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{
			Column:   "id",
			Operator: ">",
			Subquery: subquery.(*sb.Builder),
		}).
		Select([]string{"name"})
}
