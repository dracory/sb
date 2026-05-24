package tableops

import (
	"github.com/dracory/sb"
)

// CreateTable returns a CREATE TABLE SQL statement
func CreateTable() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_INTEGER,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "name",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
		}).
		Column(sb.Column{
			Name: "email",
			Type: sb.COLUMN_TYPE_STRING,
		}).
		Create()
}

// DropTable returns a DROP TABLE SQL statement
func DropTable() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Drop()
}

// DropTableIfExists returns a DROP TABLE IF EXISTS SQL statement
func DropTableIfExists() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIfExists()
}

// TruncateTable returns a TRUNCATE TABLE SQL statement
func TruncateTable() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Truncate()
}

// TruncateTableWithCascade returns a TRUNCATE TABLE with CASCADE (PostgreSQL)
func TruncateTableWithCascade() (string, error) {
	return sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true})
}

// TruncateTableWithResetIdentity returns a TRUNCATE TABLE with identity reset (MSSQL)
func TruncateTableWithResetIdentity() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{ResetIdentity: true})
}
