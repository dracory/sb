package sb_test

import (
	"testing"

	"github.com/dracory/sb/integration_tests/common/builder"
)

// TestSQLiteCreateTable tests table creation on SQLite
func TestSQLiteCreateTable(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.CreateTable(t, driver)
}

// TestSQLiteDropTable tests table dropping on SQLite
func TestSQLiteDropTable(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.DropTable(t, driver)
}

// TestSQLiteDropTableIfExists tests safe table dropping on SQLite
func TestSQLiteDropTableIfExists(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.DropTableIfExists(t, driver)
}

// TestSQLiteInsert tests INSERT operations on SQLite
func TestSQLiteInsert(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.Insert(t, driver)
}

// TestSQLiteSelect tests SELECT operations on SQLite
func TestSQLiteSelect(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.Select(t, driver)
}

// TestSQLiteUpdate tests UPDATE operations on SQLite
func TestSQLiteUpdate(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.Update(t, driver)
}

// TestSQLiteDelete tests DELETE operations on SQLite
func TestSQLiteDelete(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.Delete(t, driver)
}

// TestSQLiteErrorMissingTable tests error handling for missing table on SQLite
func TestSQLiteErrorMissingTable(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.ErrorMissingTable(t, driver)
}

// TestSQLiteErrorEmptyJoinCondition tests error handling for empty JOIN condition on SQLite
func TestSQLiteErrorEmptyJoinCondition(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.ErrorEmptyJoinCondition(t, driver)
}

// TestSQLiteErrorNilSubquery tests error handling for nil subquery on SQLite
func TestSQLiteErrorNilSubquery(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.ErrorNilSubquery(t, driver)
}

// TestSQLiteErrorEmptyColumnName tests error handling for empty column name on SQLite
func TestSQLiteErrorEmptyColumnName(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.ErrorEmptyColumnName(t, driver)
}

// TestSQLiteErrorEmptyColumnType tests error handling for empty column type on SQLite
func TestSQLiteErrorEmptyColumnType(t *testing.T) {
	driver := &SQLiteDriver{}
	builder.ErrorEmptyColumnType(t, driver)
}
