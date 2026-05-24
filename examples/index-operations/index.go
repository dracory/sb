package indexops

import (
	"github.com/dracory/sb"
)

// CreateSimpleIndex returns a CREATE INDEX SQL statement
func CreateSimpleIndex() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndex("idx_users_email", "email")
}

// CreateUniqueIndex returns a CREATE UNIQUE INDEX SQL statement
func CreateUniqueIndex() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateUniqueIndex("idx_users_email", "email")
}

// CreateCompositeIndex returns a CREATE INDEX with multiple columns
func CreateCompositeIndex() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		CreateCompositeIndex("idx_orders_user_status", []sb.IndexColumn{
			{Name: "user_id"},
			{Name: "status"},
		})
}

// CreatePartialIndex returns a CREATE INDEX with WHERE clause
func CreatePartialIndex() (string, error) {
	return sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		CreatePartialIndex("idx_orders_active", "status = 'active'", "user_id")
}

// DropIndex returns a DROP INDEX SQL statement
func DropIndex() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndex("idx_users_email")
}

// DropIndexIfExists returns a DROP INDEX IF EXISTS SQL statement
func DropIndexIfExists() (string, error) {
	return sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("cache").
		DropIndexIfExists("idx_cache_key")
}

// DropIndexWithSchema returns a DROP INDEX with schema (PostgreSQL)
func DropIndexWithSchema() (string, error) {
	return sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")
}
