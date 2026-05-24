package builder

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dracory/sb/integration_tests/common"
	"github.com/stretchr/testify/require"
)

// ErrorMissingTable tests error handling for missing table
func ErrorMissingTable(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Try to select from non-existent table
	sql, params, err := sb.NewBuilder(driver.Dialect()).
		Table("non_existent_table").
		Select([]string{"*"})

	require.NoError(t, err)

	// Execute should fail
	_, err = db.Query(sql, params...)
	require.Error(t, err, "Should error when querying non-existent table")
}

// ErrorEmptyJoinCondition tests error handling for empty JOIN condition
func ErrorEmptyJoinCondition(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Setup table
	err = common.CreateTestTableFromStruct(db, common.TestUsersTable, driver.Dialect())
	require.NoError(t, err)
	defer common.DropTestTable(db, common.TestUsersTable.Name)

	// Try to JOIN with empty ON condition
	_, _, err = sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Join(sb.JOIN_INNER, "other_table", "").
		Select([]string{"*"})

	require.Error(t, err, "Should error when JOIN condition is empty")
	require.Contains(t, err.Error(), "ON condition", "Error should mention ON condition")
}

// ErrorNilSubquery tests error handling for nil subquery
func ErrorNilSubquery(t *testing.T, driver common.TestDriver) {
	// Try to use nil subquery
	_, err := sb.NewBuilder(driver.Dialect()).
		Table("test_table").
		Exists(nil)

	require.Error(t, err, "Should error when subquery is nil")
	require.Contains(t, err.Error(), "subquery", "Error should mention subquery")
}

// ErrorEmptyColumnName tests error handling for empty column name
func ErrorEmptyColumnName(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Try to create table with empty column name
	_, err = sb.NewBuilder(driver.Dialect()).
		Table("test_table").
		Column(sb.Column{Name: "", Type: sb.COLUMN_TYPE_STRING}).
		CreateIfNotExists()

	require.Error(t, err, "Should error when column name is empty")
	require.Contains(t, err.Error(), "column name", "Error should mention column name")
}

// ErrorEmptyColumnType tests error handling for empty column type
func ErrorEmptyColumnType(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Try to create table with empty column type
	_, err = sb.NewBuilder(driver.Dialect()).
		Table("test_table").
		Column(sb.Column{Name: "test_col", Type: ""}).
		CreateIfNotExists()

	require.Error(t, err, "Should error when column type is empty")
	require.Contains(t, err.Error(), "column type", "Error should mention column type")
}
