package sb_test

import (
	"testing"

	"github.com/dracory/sb/integration_tests/common/builder"
)

// TestMySQLCreateTable tests table creation on MySQL
func TestMySQLCreateTable(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.CreateTable(t, driver)
}

// TestMySQLDropTable tests table dropping on MySQL
func TestMySQLDropTable(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.DropTable(t, driver)
}

// TestMySQLDropTableIfExists tests safe table dropping on MySQL
func TestMySQLDropTableIfExists(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.DropTableIfExists(t, driver)
}

// TestMySQLInsert tests INSERT operations on MySQL
func TestMySQLInsert(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.Insert(t, driver)
}

// TestMySQLSelect tests SELECT operations on MySQL
func TestMySQLSelect(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.Select(t, driver)
}

// TestMySQLUpdate tests UPDATE operations on MySQL
func TestMySQLUpdate(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.Update(t, driver)
}

// TestMySQLDelete tests DELETE operations on MySQL
func TestMySQLDelete(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.Delete(t, driver)
}

// TestMySQLErrorMissingTable tests error handling for missing table on MySQL
func TestMySQLErrorMissingTable(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.ErrorMissingTable(t, driver)
}

// TestMySQLErrorEmptyJoinCondition tests error handling for empty JOIN condition on MySQL
func TestMySQLErrorEmptyJoinCondition(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.ErrorEmptyJoinCondition(t, driver)
}

// TestMySQLErrorNilSubquery tests error handling for nil subquery on MySQL
func TestMySQLErrorNilSubquery(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.ErrorNilSubquery(t, driver)
}

// TestMySQLErrorEmptyColumnName tests error handling for empty column name on MySQL
func TestMySQLErrorEmptyColumnName(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.ErrorEmptyColumnName(t, driver)
}

// TestMySQLErrorEmptyColumnType tests error handling for empty column type on MySQL
func TestMySQLErrorEmptyColumnType(t *testing.T) {
	driver := &MySQLDriver{}
	driver.SkipIfUnavailable(t)
	builder.ErrorEmptyColumnType(t, driver)
}
