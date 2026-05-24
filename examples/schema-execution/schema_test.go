package schemaexec

import (
	"testing"

	"github.com/dracory/database"
	"github.com/dracory/sb"
)

func TestTableColumnsExample(t *testing.T) {
	// Verify function signature by assignment
	var _ func(database.QueryableContext) ([]sb.Column, error) = TableColumnsExample
}

func TestTableColumnExistsExample(t *testing.T) {
	var _ func(database.QueryableContext) (bool, error) = TableColumnExistsExample
}

func TestTableColumnAddExample(t *testing.T) {
	var _ func(database.QueryableContext) error = TableColumnAddExample
}

func TestTableColumnAddIfNotExistsExample(t *testing.T) {
	var _ func(database.QueryableContext) error = TableColumnAddIfNotExistsExample
}

func TestTableColumnDropExample(t *testing.T) {
	var _ func(database.QueryableContext) error = TableColumnDropExample
}

func TestTableColumnDropIfExistsExample(t *testing.T) {
	var _ func(database.QueryableContext) error = TableColumnDropIfExistsExample
}

func TestTableColumnRenameExample(t *testing.T) {
	var _ func(database.QueryableContext) error = TableColumnRenameExample
}
