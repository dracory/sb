package columnops

import (
	"github.com/dracory/sb"
)

// AddColumn returns an ALTER TABLE ADD COLUMN SQL statement
func AddColumn() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		ColumnAdd(sb.Column{
			Name:     "phone",
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   20,
			Nullable: true,
		})
}

// DropColumn returns an ALTER TABLE DROP COLUMN SQL statement
func DropColumn() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		ColumnDrop("temp_column")
}

// RenameColumn returns an ALTER TABLE RENAME COLUMN SQL statement
func RenameColumn() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		ColumnRename("email", "new_email")
}

// ColumnExists returns the SQL query and parameters for checking if a column exists
func ColumnExists() (string, []interface{}, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		ColumnExists("email")
}
