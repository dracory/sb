# Unified Driver-Agnostic Integration Test Suite

This plan creates a comprehensive, driver-agnostic integration test suite that tests both the SB builder and schema packages with actual SQL execution against real databases.

## Current State Analysis

**Existing Tests:**
- SQLite: 2 test files (sqlite_test.go, error_handling_test.go) - very light coverage
- MySQL: 1 test file (mysql_test.go) - basic SELECT only
- PostgreSQL: 1 test file (postgres_test.go) - basic SELECT only
- Schema package: Has tests but they're driver-specific (separate test functions per driver)

**Gaps Identified:**
- No unified test structure across drivers
- Minimal builder coverage (only basic SELECT tested)
- No schema package integration tests
- No actual SQL execution verification for most operations
- Tests are duplicated per driver instead of shared

## Proposed Architecture

### 1. Unified Test Structure

Keep existing driver-specific folder structure. Create unified test files that are shared across all drivers:

```
integration_tests/
  common/
    setup.go                 # Common test setup/helpers (enhanced)
    builder/
      table_test.go          # Shared builder tests for table operations
      crud_test.go           # Shared builder tests for CRUD operations
      join_test.go           # Shared builder tests for JOIN operations
      subquery_test.go       # Shared builder tests for subquery operations
      index_test.go          # Shared builder tests for index operations
      column_test.go         # Shared builder tests for column operations
      truncate_test.go       # Shared builder tests for TRUNCATE operations
      view_test.go           # Shared builder tests for view operations
      error_test.go          # Shared builder tests for error handling
    schema/
      table_test.go          # Shared schema tests for table functions
      column_test.go         # Shared schema tests for column functions
      helpers_test.go        # Shared schema tests for helper functions
  sqlite/
    sqlite_test.go          # SQLite-specific test runner
  mysql/
    mysql_test.go           # MySQL-specific test runner
  postgres/
    postgres_test.go        # PostgreSQL-specific test runner
  mssql/
    mssql_test.go           # MSSQL-specific test runner (create if needed)
```

### 2. Driver Interface

Define a `TestDriver` interface for database-specific operations:

```go
type TestDriver interface {
    Name() string
    Dialect() string
    Connect() (*sql.DB, error)
    SkipIfUnavailable(t *testing.T)
    Cleanup(t *testing.T, db *sql.DB)
}
```

### 3. Test Categories

#### A. Builder Tests (organized by feature file)

**table_test.go - Table Operations:**
- Create() / CreateIfNotExists()
- Drop() / DropIfExists()
- Rename()

**crud_test.go - CRUD Operations:**
- Insert() with various data types
- Select() with WHERE, ORDER BY, LIMIT, OFFSET, GROUP BY
- Update() with WHERE conditions
- Delete() with WHERE conditions

**join_test.go - JOIN Operations:**
- JOIN operations (INNER, LEFT, RIGHT, FULL, CROSS)
- JOIN with aliases
- Multiple JOINs in single query
- JOIN with WHERE and ORDER BY

**subquery_test.go - Subquery Operations:**
- Subqueries (IN, NOT IN, EXISTS, NOT EXISTS)
- Correlated subqueries
- Subqueries in comparison operators
- Multiple subquery conditions

**index_test.go - Index Operations:**
- CreateIndex() / CreateUniqueIndex()
- CreateIndexWithOptions() (advanced features)
- CreateCompositeIndex() / CreatePartialIndex() / CreateCoveringIndex()
- DropIndex() / DropIndexIfExists()
- DropIndexWithOptions()

**column_test.go - Column Operations:**
- ColumnAdd() / ColumnChange() / ColumnDrop() / ColumnRename()
- TableColumnAdd() / TableColumnChange() / TableColumnDrop() / TableColumnRename() (deprecated)
- ColumnExists()

**truncate_test.go - TRUNCATE Operations:**
- Truncate() / TruncateWithOptions()
- Database-specific options (CASCADE, ResetIdentity)

**view_test.go - View Operations:**
- View() / ViewSQL() / ViewColumns()
- View creation and querying

**error_test.go - Error Handling:**
- Missing table errors
- Empty JOIN condition errors
- Invalid dialect errors
- Nil subquery errors
- Empty column name/type errors

#### B. Schema Tests (organized by feature file)

**table_test.go - Table Schema Functions:**
- TableCreate() - create tables with various column types
- TableDrop() / TableDropIfExists()
- TableColumns() - retrieve and validate column metadata

**column_test.go - Column Schema Functions:**
- TableColumnAdd() / TableColumnAddIfNotExists()
- TableColumnDrop() / TableColumnDropIfExists()
- TableColumnRename()
- TableColumnExists() - check column existence

**helpers_test.go - Schema Helper Functions:**
- rawColumnProcess() - column type processing
- Helper function validation

### 4. Test Data Strategy

**Test Tables:**
- `test_users` - basic user table (id, name, email, status, created_at)
- `test_orders` - orders table for JOIN testing (id, user_id, total, status)
- `test_products` - products table for complex queries (id, name, price, category)

**Test Data Helpers:**
- `InsertTestUsers()` - populate users table
- `InsertTestOrders()` - populate orders table
- `InsertTestProducts()` - populate products table
- `CleanupAllTables()` - drop all test tables

### 5. Test Execution Pattern

**Shared test functions** in common/ accept a TestDriver parameter:

```go
// common/builder/table_test.go
func TestBuilderCreateTable(t *testing.T, driver TestDriver) {
    db, err := driver.Connect()
    require.NoError(t, err)
    defer driver.Cleanup(t, db)

    // Test logic using driver.Dialect() for SB builder
    sql, err := sb.NewBuilder(driver.Dialect()).
        Table("test_users").
        Column(...).
        CreateIfNotExists()
    
    _, err = db.Exec(sql)
    require.NoError(t, err)
}
```

**Driver-specific test runners** call shared tests with their driver:

```go
// sqlite/sqlite_test.go
func TestSQLiteTableOperations(t *testing.T) {
    TestBuilderCreateTable(t, SQLiteDriver{})
    TestBuilderDropTable(t, SQLiteDriver{})
}

// mysql/mysql_test.go
func TestMySQLTableOperations(t *testing.T) {
    TestBuilderCreateTable(t, MySQLDriver{})
    TestBuilderDropTable(t, MySQLDriver{})
}
```

**Key points:**
- Shared test logic in common/ uses driver.Dialect() for SB builder
- SB library handles dialect differences internally (quoting, syntax)
- Driver-specific code is only in TestDriver implementations (connection setup)
- Same test runs on all drivers with minimal code duplication

### 6. Implementation Phases

**Phase 1: Infrastructure (Priority: High)**
- Enhance common/common.go with TestDriver interface and helpers
- Create common/builder/ and common/schema/ directories
- Create driver-specific test runners in existing driver folders
- Update common.CreateTestTable to support multiple tables

**Phase 2: Builder Core Tests (Priority: High)**
- Create common/builder/table_test.go - Table operations (Create, Drop, Rename)
- Create common/builder/crud_test.go - CRUD operations (Insert, Select, Update, Delete)
- Create common/builder/error_test.go - Error handling tests
- Update driver-specific test runners to call shared tests
- Verify all tests pass on SQLite (in-memory)

**Phase 3: Builder Advanced Tests (Priority: Medium)**
- Create common/builder/join_test.go - JOIN operations
- Create common/builder/subquery_test.go - Subquery operations
- Create common/builder/index_test.go - Index operations
- Create common/builder/column_test.go - Column operations
- Create common/builder/truncate_test.go - TRUNCATE operations
- Create common/builder/view_test.go - View operations

**Phase 4: Schema Tests (Priority: High)**
- Create common/schema/table_test.go - Table schema functions
- Create common/schema/column_test.go - Column schema functions
- Create common/schema/helpers_test.go - Schema helper functions

**Phase 5: Cross-Driver Validation (Priority: Medium)**
- Run all tests on MySQL (in CI)
- Run all tests on PostgreSQL (in CI)
- Run all tests on MSSQL (in CI)
- Fix any driver-specific issues
- Document driver-specific limitations

### 7. File Changes Summary

**New Files:**
- `integration_tests/common/builder/table_test.go` - Shared table operations tests
- `integration_tests/common/builder/crud_test.go` - Shared CRUD operations tests
- `integration_tests/common/builder/join_test.go` - Shared JOIN operations tests
- `integration_tests/common/builder/subquery_test.go` - Shared subquery operations tests
- `integration_tests/common/builder/index_test.go` - Shared index operations tests
- `integration_tests/common/builder/column_test.go` - Shared column operations tests
- `integration_tests/common/builder/truncate_test.go` - Shared TRUNCATE operations tests
- `integration_tests/common/builder/view_test.go` - Shared view operations tests
- `integration_tests/common/builder/error_test.go` - Shared error handling tests
- `integration_tests/common/schema/table_test.go` - Shared table schema tests
- `integration_tests/common/schema/column_test.go` - Shared column schema tests
- `integration_tests/common/schema/helpers_test.go` - Shared schema helper tests
- `integration_tests/mssql/mssql_test.go` - MSSQL test runner (if not exists)

**Modified Files:**
- `integration_tests/common/common.go` - Enhance with TestDriver interface and multi-table support
- `integration_tests/sqlite/sqlite_test.go` - Update to call shared tests
- `integration_tests/mysql/mysql_test.go` - Update to call shared tests
- `integration_tests/postgres/postgres_test.go` - Update to call shared tests

**Deprecated Files (Phase 5):**
- `integration_tests/sqlite/error_handling_test.go` - Merge into common/builder/error_test.go

### 8. Test Coverage Goals

**Builder Coverage:**
- Table operations: 100%
- CRUD operations: 100%
- JOIN operations: 100%
- Subqueries: 100%
- Index operations: 100%
- Column operations: 100%
- Error handling: 100%

**Schema Coverage:**
- Table functions: 100%
- Column functions: 100%
- Error scenarios: 100%

**Driver Coverage:**
- SQLite: 100% (always available)
- MySQL: 100% (CI only)
- PostgreSQL: 100% (CI only)
- MSSQL: 100% (CI only)

### 9. Success Criteria

- ✅ All builder methods tested with actual SQL execution
- ✅ All schema functions tested with actual SQL execution
- ✅ Same test files run across all drivers
- ✅ Only driver-specific code is in driver implementations
- ✅ All tests pass on SQLite (local development)
- ✅ All tests pass on MySQL/PostgreSQL/MSSQL (CI)
- ✅ Old driver-specific tests deprecated
- ✅ Test execution time < 5 minutes for full suite

### 10. CI/CD Integration

Update `.github/workflows/tests.yml` to:
- Run unified integration tests on all drivers
- Skip MySQL/PostgreSQL/MSSQL tests if services unavailable
- Report test coverage per driver
- Fail build if any driver's tests fail

## Estimated Effort

- Phase 1: 4-6 hours
- Phase 2: 6-8 hours
- Phase 3: 8-10 hours
- Phase 4: 4-6 hours
- Phase 5: 4-6 hours

**Total: 26-36 hours**
