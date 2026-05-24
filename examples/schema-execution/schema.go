package schemaexec

import (
	"github.com/dracory/database"
	"github.com/dracory/sb"
	"github.com/dracory/sb/schema"
)

// TableColumnsExample demonstrates getting table columns
// Usage: columns, err := schema.TableColumns(qctx, "users", true)
func TableColumnsExample(qctx database.QueryableContext) ([]sb.Column, error) {
	return schema.TableColumns(qctx, "users", true)
}

// TableColumnExistsExample demonstrates checking if a column exists
// Usage: exists, err := schema.TableColumnExists(qctx, "users", "email")
func TableColumnExistsExample(qctx database.QueryableContext) (bool, error) {
	return schema.TableColumnExists(qctx, "users", "email")
}

// TableColumnAddExample demonstrates adding a column
// Usage: err := schema.TableColumnAdd(qctx, "users", column)
func TableColumnAddExample(qctx database.QueryableContext) error {
	column := sb.Column{
		Name:     "phone",
		Type:     sb.COLUMN_TYPE_STRING,
		Length:   20,
		Nullable: true,
	}
	return schema.TableColumnAdd(qctx, "users", column)
}

// TableColumnAddIfNotExistsExample demonstrates adding a column if it doesn't exist
// Usage: err := schema.TableColumnAddIfNotExists(qctx, "users", column)
func TableColumnAddIfNotExistsExample(qctx database.QueryableContext) error {
	column := sb.Column{
		Name:     "phone",
		Type:     sb.COLUMN_TYPE_STRING,
		Length:   20,
		Nullable: true,
	}
	return schema.TableColumnAddIfNotExists(qctx, "users", column)
}

// TableColumnDropExample demonstrates dropping a column
// Usage: err := schema.TableColumnDrop(qctx, "users", "temp_column")
func TableColumnDropExample(qctx database.QueryableContext) error {
	return schema.TableColumnDrop(qctx, "users", "temp_column")
}

// TableColumnDropIfExistsExample demonstrates dropping a column if it exists
// Usage: err := schema.TableColumnDropIfExists(qctx, "users", "temp_column")
func TableColumnDropIfExistsExample(qctx database.QueryableContext) error {
	return schema.TableColumnDropIfExists(qctx, "users", "temp_column")
}

// TableColumnRenameExample demonstrates renaming a column
// Usage: err := schema.TableColumnRename(qctx, "users", "email", "new_email")
func TableColumnRenameExample(qctx database.QueryableContext) error {
	return schema.TableColumnRename(qctx, "users", "email", "new_email")
}
