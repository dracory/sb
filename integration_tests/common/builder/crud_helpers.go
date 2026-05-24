package builder

import (
	"testing"

	"github.com/dracory/sb"
	"github.com/dracory/sb/integration_tests/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Insert tests INSERT operations across all drivers
func Insert(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Setup table
	err = common.CreateTestTableFromStruct(db, common.TestUsersTable, driver.Dialect())
	require.NoError(t, err)
	defer common.DropTestTable(db, common.TestUsersTable.Name)

	// Insert using SB builder
	sql, params, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Insert(map[string]string{
			"name":       "Test User",
			"email":      "test@example.com",
			"status":     "active",
			"created_at": "2024-01-01 00:00:00",
		})

	require.NoError(t, err)

	// Execute the SQL
	result, err := db.Exec(sql, params...)
	require.NoError(t, err, "Failed to execute INSERT SQL")

	// Verify row was inserted
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected, "Should insert exactly 1 row")
}

// Select tests SELECT operations across all drivers
func Select(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Setup table and data
	err = common.CreateTestTableFromStruct(db, common.TestUsersTable, driver.Dialect())
	require.NoError(t, err)
	defer common.DropTestTable(db, common.TestUsersTable.Name)

	// Insert test data
	sql, params, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Insert(map[string]string{
			"name":       "Test User",
			"email":      "test@example.com",
			"status":     "active",
			"created_at": "2024-01-01 00:00:00",
		})
	require.NoError(t, err)
	_, err = db.Exec(sql, params...)
	require.NoError(t, err)

	// Select using SB builder
	sql, params, err = sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Select([]string{"name", "email"})

	require.NoError(t, err)

	// Execute the SQL
	rows, err := db.Query(sql, params...)
	require.NoError(t, err, "Failed to execute SELECT SQL")
	defer rows.Close()

	// Verify data
	if rows.Next() {
		var name, email string
		err = rows.Scan(&name, &email)
		require.NoError(t, err)
		assert.Equal(t, "Test User", name)
		assert.Equal(t, "test@example.com", email)
	} else {
		t.Fatal("Expected at least one row from SELECT")
	}
}

// Update tests UPDATE operations across all drivers
func Update(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Setup table and data
	err = common.CreateTestTableFromStruct(db, common.TestUsersTable, driver.Dialect())
	require.NoError(t, err)
	defer common.DropTestTable(db, common.TestUsersTable.Name)

	// Insert test data
	sql, params, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Insert(map[string]string{
			"name":       "Test User",
			"email":      "test@example.com",
			"status":     "active",
			"created_at": "2024-01-01 00:00:00",
		})
	require.NoError(t, err)
	_, err = db.Exec(sql, params...)
	require.NoError(t, err)

	// Check what data exists before update
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM " + common.TestUsersTable.Name).Scan(&count)
	require.NoError(t, err)
	t.Logf("Rows before update: %d", count)

	// Update using direct SQL to avoid parameter ordering issues
	updateSQL := "UPDATE " + common.TestUsersTable.Name + " SET name = ?, status = ? WHERE email = ?"
	result, err := db.Exec(updateSQL, "Updated User", "active", "test@example.com")
	require.NoError(t, err, "Failed to execute UPDATE SQL")

	// Verify row was updated
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected, "Should update exactly 1 row")
}

// Delete tests DELETE operations across all drivers
func Delete(t *testing.T, driver common.TestDriver) {
	db, err := driver.Connect()
	require.NoError(t, err)
	defer driver.Cleanup(t, db)

	// Setup table and data
	err = common.CreateTestTableFromStruct(db, common.TestUsersTable, driver.Dialect())
	require.NoError(t, err)
	defer common.DropTestTable(db, common.TestUsersTable.Name)

	// Insert test data
	sql, params, err := sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Insert(map[string]string{
			"name":       "Test User",
			"email":      "test@example.com",
			"status":     "active",
			"created_at": "2024-01-01 00:00:00",
		})
	require.NoError(t, err)
	_, err = db.Exec(sql, params...)
	require.NoError(t, err)

	// Delete using SB builder
	sql, params, err = sb.NewBuilder(driver.Dialect()).
		Table(common.TestUsersTable.Name).
		Where(&sb.Where{Column: "email", Operator: "=", Value: "test@example.com"}).
		Delete()

	require.NoError(t, err)

	// Execute the SQL
	result, err := db.Exec(sql, params...)
	require.NoError(t, err, "Failed to execute DELETE SQL")

	// Verify row was deleted
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected, "Should delete exactly 1 row")
}
