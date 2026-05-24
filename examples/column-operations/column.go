package columnops

import (
	"github.com/dracory/sb"
)

// AddColumn returns an ALTER TABLE ADD COLUMN SQL statement
func AddColumn() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		TableColumnAdd("users", sb.Column{
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
		TableColumnDrop("users", "temp_column")
}

// RenameColumn returns an ALTER TABLE RENAME COLUMN SQL statement
func RenameColumn() (string, error) {
	return sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		TableColumnRename("users", "email", "new_email")
}
