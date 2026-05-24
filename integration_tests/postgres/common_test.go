package sb_test

import (
	"database/sql"
	"fmt"

	"github.com/dracory/sb"
)

// CreateTestTable creates a simple test table for integration testing using SB library
func CreateTestTable(db *sql.DB, tableName string, dialect string) error {
	// Create table using SB library with proper dialect-specific syntax
	createSQL, err := sb.NewBuilder(dialect).
		Table(tableName).
		Column(sb.Column{
			Name:          "id",
			Type:          sb.COLUMN_TYPE_INTEGER,
			AutoIncrement: true,
			PrimaryKey:    true,
		}).
		Column(sb.Column{
			Name:     "name",
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   100,
			Nullable: false, // NOT NULL
		}).
		Column(sb.Column{
			Name:     "email",
			Type:     sb.COLUMN_TYPE_STRING,
			Length:   100,
			Nullable: false, // NOT NULL
		}).
		Column(sb.Column{
			Name:    "status",
			Type:    sb.COLUMN_TYPE_STRING,
			Length:  20,
			Default: "active",
		}).
		Column(sb.Column{
			Name:    "created_at",
			Type:    sb.COLUMN_TYPE_DATETIME,
			Default: "CURRENT_TIMESTAMP",
		}).
		CreateIfNotExists()

	if err != nil {
		return fmt.Errorf("failed to generate CREATE TABLE SQL: %w", err)
	}

	// Log the SQL for debugging
	fmt.Printf("DEBUG: Generated SQL for %s:\n%s\n", dialect, createSQL)

	_, err = db.Exec(createSQL)
	if err != nil {
		return fmt.Errorf("failed to execute CREATE TABLE: %w", err)
	}
	return nil
}

// DropTestTable drops a test table after testing
func DropTestTable(db *sql.DB, tableName string) error {
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)
	_, err := db.Exec(dropSQL)
	if err != nil {
		return fmt.Errorf("failed to drop test table %s: %w", tableName, err)
	}
	return nil
}
