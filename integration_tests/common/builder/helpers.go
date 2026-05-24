package builder

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dracory/sb/integration_tests/common"
	"github.com/stretchr/testify/require"
)

// CreateTable tests table creation across all drivers
func CreateTable(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Create table using SB builder
	createSQL, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Column(common.TestUsersTable.Columns[0]).
		Column(common.TestUsersTable.Columns[1]).
		Column(common.TestUsersTable.Columns[2]).
		Column(common.TestUsersTable.Columns[3]).
		Column(common.TestUsersTable.Columns[4]).
		CreateIfNotExists()

	require.NoError(t, err)

	// Execute the SQL
	_, err = db.Exec(createSQL)
	require.NoError(t, err, "Failed to execute CREATE TABLE SQL")

	// Verify table exists by attempting to query it
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + common.TestUsersTable.Name).Scan(&count)
	require.NoError(t, err, "Table should exist after creation")
}

// DropTable tests table dropping across all drivers
func DropTable(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// First create the table
	err = common.CreateTestTableFromStruct(db, common.TestUsersTable, driver.Dialect())
	require.NoError(t, err)

	// Drop table using SB builder
	dropSQL, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Drop()

	require.NoError(t, err)

	// Execute the SQL
	_, err = db.Exec(dropSQL)
	require.NoError(t, err, "Failed to execute DROP TABLE SQL")

	// Verify table no longer exists
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + common.TestUsersTable.Name).Scan(&count)
	require.Error(t, err, "Table should not exist after drop")
}

// DropTableIfExists tests safe table dropping across all drivers
func DropTableIfExists(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Drop table that doesn't exist - should not error
	dropSQL, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		DropIfExists()

	require.NoError(t, err)

	// Execute the SQL - should succeed even if table doesn't exist
	_, err = db.Exec(dropSQL)
	require.NoError(t, err, "DROP IF EXISTS should not error on missing table")
}
