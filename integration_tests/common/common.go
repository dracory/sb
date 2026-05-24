package common

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/dracory/sb"
)

// TestDriver defines the interface for database-specific test operations
type TestDriver interface {
	Name() string
	Dialect() string
	Connect() (*sql.DB, error)
	SkipIfUnavailable(t *testing.T)
	Cleanup(t *testing.T, db *sql.DB)
}

// TestTable represents a test table structure
type TestTable struct {
	Name    string
	Columns []sb.Column
}

// Common test tables
var (
	TestUsersTable = TestTable{
		Name: "test_users",
		Columns: []sb.Column{
			{
				Name:          "id",
				Type:          sb.COLUMN_TYPE_INTEGER,
				AutoIncrement: true,
				PrimaryKey:    true,
			},
			{
				Name:     "name",
				Type:     sb.COLUMN_TYPE_STRING,
				Length:   100,
				Nullable: false,
			},
			{
				Name:     "email",
				Type:     sb.COLUMN_TYPE_STRING,
				Length:   100,
				Nullable: false,
			},
			{
				Name:    "status",
				Type:    sb.COLUMN_TYPE_STRING,
				Length:  20,
				Default: "active",
			},
			{
				Name:    "created_at",
				Type:    sb.COLUMN_TYPE_DATETIME,
				Default: "CURRENT_TIMESTAMP",
			},
		},
	}

	TestOrdersTable = TestTable{
		Name: "test_orders",
		Columns: []sb.Column{
			{
				Name:          "id",
				Type:          sb.COLUMN_TYPE_INTEGER,
				AutoIncrement: true,
				PrimaryKey:    true,
			},
			{
				Name:     "user_id",
				Type:     sb.COLUMN_TYPE_INTEGER,
				Nullable: false,
			},
			{
				Name:     "total",
				Type:     sb.COLUMN_TYPE_DECIMAL,
				Length:   10,
				Decimals: 2,
				Nullable: false,
			},
			{
				Name:    "status",
				Type:    sb.COLUMN_TYPE_STRING,
				Length:  20,
				Default: "pending",
			},
		},
	}

	TestProductsTable = TestTable{
		Name: "test_products",
		Columns: []sb.Column{
			{
				Name:          "id",
				Type:          sb.COLUMN_TYPE_INTEGER,
				AutoIncrement: true,
				PrimaryKey:    true,
			},
			{
				Name:     "name",
				Type:     sb.COLUMN_TYPE_STRING,
				Length:   200,
				Nullable: false,
			},
			{
				Name:     "price",
				Type:     sb.COLUMN_TYPE_DECIMAL,
				Length:   10,
				Decimals: 2,
				Nullable: false,
			},
			{
				Name:     "category",
				Type:     sb.COLUMN_TYPE_STRING,
				Length:   50,
				Nullable: false,
			},
		},
	}
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

// CreateTestTableFromStruct creates a test table from a TestTable structure
func CreateTestTableFromStruct(db *sql.DB, table TestTable, dialect string) error {
	builder := sb.NewBuilder(dialect).Table(table.Name)
	for _, col := range table.Columns {
		builder.Column(col)
	}

	createSQL, err := builder.CreateIfNotExists()
	if err != nil {
		return fmt.Errorf("failed to generate CREATE TABLE SQL: %w", err)
	}

	_, err = db.Exec(createSQL)
	return err
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

// CleanupAllTables drops all test tables
func CleanupAllTables(db *sql.DB) error {
	tables := []string{
		TestUsersTable.Name,
		TestOrdersTable.Name,
		TestProductsTable.Name,
	}

	for _, table := range tables {
		if err := DropTestTable(db, table); err != nil {
			return err
		}
	}
	return nil
}
