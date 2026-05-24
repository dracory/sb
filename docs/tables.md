# Table Operations

This document covers SQL generation for table operations (using the `sb` package). For **database execution** functions (CREATE, DROP, column operations), see the [Schema Execution](#schema-execution) section below, which uses the `schema` sub-package.

## CREATE TABLE

Create tables with comprehensive column definitions and database-specific dialect support.

### Basic Example

```go
import "github.com/dracory/sb"

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Column(sb.Column{
		Name:       "id",
		Type:       sb.COLUMN_TYPE_STRING,
		Length:     40,
		PrimaryKey: true,
	}).
	Column(sb.Column{
		Name: "image",
		Type: sb.COLUMN_TYPE_BLOB,
	}).
	Column(sb.Column{
		Name: "price_default",
		Type: sb.COLUMN_TYPE_DECIMAL,
	}).
	Column(sb.Column{
		Name:     "price_custom",
		Type:     sb.COLUMN_TYPE_DECIMAL,
		Length:   12,
		Decimals: 10,
	}).
	Column(sb.Column{
		Name: "created_at",
		Type: sb.COLUMN_TYPE_DATETIME,
	}).
	Column(sb.Column{
		Name:     "deleted_at",
		Type:     sb.COLUMN_TYPE_DATETIME,
		Nullable: true,
	}).
	Create()
```

### Column Types

- `COLUMN_TYPE_STRING` - String/varchar columns
- `COLUMN_TYPE_INT` - Integer columns
- `COLUMN_TYPE_DECIMAL` - Decimal/numeric columns
- `COLUMN_TYPE_BLOB` - Binary data columns
- `COLUMN_TYPE_DATETIME` - Date/time columns

### Column Options

- `Name` - Column name
- `Type` - Column data type
- `Length` - Column length (for strings, decimals)
- `Decimals` - Decimal places (for decimal columns)
- `PrimaryKey` - Set as primary key
- `Nullable` - Allow NULL values

## DROP TABLE

Remove tables from the database.

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Drop()
```

## TRUNCATE TABLE

Remove all data from a table while preserving structure.

### Basic Truncate

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
	Table("users").
	Truncate()
```

### Truncate with Options

#### PostgreSQL with CASCADE

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
	Table("orders").
	TruncateWithOptions(sb.TruncateOptions{Cascade: true})
// Result: TRUNCATE TABLE "orders" CASCADE;
```

#### MSSQL with Identity Reset

```go
sql := sb.NewBuilder(sb.DIALECT_MSSQL).
	Table("users").
	TruncateWithOptions(sb.TruncateOptions{ResetIdentity: true})
// Result: TRUNCATE TABLE [users]; DBCC CHECKIDENT ('users', RESEED, 0)
```

### Database-Specific Behavior

- **MySQL**: `TRUNCATE TABLE table_name;`
- **PostgreSQL**: `TRUNCATE TABLE "table_name";` (with CASCADE option)
- **SQLite**: Falls back to `DELETE FROM "table_name";`
- **MSSQL**: `TRUNCATE TABLE [table_name];` (with identity reset option)

### TruncateOptions

```go
type TruncateOptions struct {
    Cascade        bool // PostgreSQL: CASCADE for foreign key constraints
    ResetIdentity  bool // MSSQL: Reset identity column seed
}
```

## DROP COLUMN

Remove columns from an existing table.

### Builder Method (SQL Generation)

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    TableColumnDrop("users", "temp_column")
// Result: ALTER TABLE `users` DROP COLUMN `temp_column`;
```

### Standalone Function (Execution)

```go
import "github.com/dracory/database"

// Drop column directly using database context
err := sb.TableColumnDrop(database.Context(ctx, db), "users", "temp_column")
if err != nil {
    log.Fatal("Error dropping column:", err)
}
```

### Drop Column If Exists

```go
// Safe drop - no error if column doesn't exist
err := sb.TableColumnDropIfExists(database.Context(ctx, db), "users", "temp_column")
if err != nil {
    log.Fatal("Error dropping column:", err)
}
```

### Database-Specific SQL

```sql
-- MySQL: ALTER TABLE `users` DROP COLUMN `temp_column`;
-- PostgreSQL: ALTER TABLE "users" DROP COLUMN "temp_column";
-- SQLite: ALTER TABLE "users" DROP COLUMN "temp_column";
-- MSSQL: ALTER TABLE [users] DROP COLUMN [temp_column];
```

## RENAME COLUMN

Rename columns in an existing table.

### Builder Method (SQL Generation)

```go
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    TableColumnRename("users", "email", "new_email")
// Result: ALTER TABLE `users` RENAME COLUMN `email` TO `new_email`;
```

### Standalone Function (Execution)

```go
import "github.com/dracory/database"

// Rename column directly using database context
err := sb.TableColumnRename(database.Context(ctx, db), "users", "email", "new_email")
if err != nil {
    log.Fatal("Error renaming column:", err)
}
```

### Database-Specific SQL

```sql
-- MySQL: ALTER TABLE `users` RENAME COLUMN `email` TO `new_email`;
-- PostgreSQL: ALTER TABLE "users" RENAME COLUMN "email" TO "new_email";
-- SQLite: ALTER TABLE "users" RENAME COLUMN "email" TO "new_email";
-- MSSQL: EXEC sp_rename [users].[email], [new_email], 'COLUMN';
```

## Table Existence Check

Check if a table exists in the database.

```go
exists := sb.TableColumnExists(db, "users", "id")
```

---

# Schema Execution

The `schema` sub-package (`github.com/dracory/sb/schema`) provides database execution functions for table and column operations. These functions execute SQL directly against the database, as opposed to the Builder methods which only generate SQL strings.

## Import

```go
import (
    "github.com/dracory/sb"
    "github.com/dracory/sb/schema"
)
```

## CREATE TABLE

Execute table creation directly against the database.

```go
ctx := context.Background()
db := sb.NewDatabaseFromDriver("sqlite3", "test.db")

columns := []sb.Column{
    {Name: "id", Type: sb.COLUMN_TYPE_INT, PrimaryKey: true},
    {Name: "name", Type: sb.COLUMN_TYPE_STRING, Length: 255},
}

err := schema.TableCreate(ctx, db, "users", columns)
if err != nil {
    log.Fatal("Error creating table:", err)
}
```

## DROP TABLE

Execute table drop directly against the database.

```go
err := schema.TableDrop(ctx, db, "users")
if err != nil {
    log.Fatal("Error dropping table:", err)
}

// Safe drop - no error if table doesn't exist
err = schema.TableDropIfExists(ctx, db, "users")
```

## COLUMN OPERATIONS

### Add Column

```go
column := sb.Column{
    Name:     "email",
    Type:     sb.COLUMN_TYPE_STRING,
    Length:   255,
    Nullable: true,
}

err := schema.TableColumnAdd(ctx, db, "users", column)
if err != nil {
    log.Fatal("Error adding column:", err)
}

// Add column if it doesn't exist
err = schema.TableColumnAddIfNotExists(ctx, db, "users", column)
```

### Drop Column

```go
err := schema.TableColumnDrop(ctx, db, "users", "temp_column")
if err != nil {
    log.Fatal("Error dropping column:", err)
}

// Safe drop - no error if column doesn't exist
err = schema.TableColumnDropIfExists(ctx, db, "users", "temp_column")
```

### Rename Column

```go
err := schema.TableColumnRename(ctx, db, "users", "email", "new_email")
if err != nil {
    log.Fatal("Error renaming column:", err)
}
```

### Check Column Exists

```go
exists, err := schema.TableColumnExists(ctx, db, "users", "email")
if err != nil {
    log.Fatal("Error checking column:", err)
}
if exists {
    fmt.Println("Column exists")
}
```

### Get Table Columns

```go
columns, err := schema.TableColumns(ctx, db, "users", true) // true = commonize types
if err != nil {
    log.Fatal("Error getting columns:", err)
}

for _, col := range columns {
    fmt.Printf("Column: %s, Type: %s\n", col.Name, col.Type)
}
```

## SQL Generation vs Execution

| Operation | SQL Generation (sb) | Execution (schema) |
|-----------|---------------------|-------------------|
| CREATE TABLE | `builder.Create()` | `schema.TableCreate()` |
| DROP TABLE | `builder.Drop()` | `schema.TableDrop()` |
| ADD COLUMN | `builder.TableColumnAdd()` | `schema.TableColumnAdd()` |
| DROP COLUMN | `builder.TableColumnDrop()` | `schema.TableColumnDrop()` |
| RENAME COLUMN | `builder.TableColumnRename()` | `schema.TableColumnRename()` |
| COLUMN EXISTS | — | `schema.TableColumnExists()` |
| GET COLUMNS | — | `schema.TableColumns()` |

**Use Builder methods** when you need to:
- Generate SQL strings for logging/debugging
- Build complex queries before execution
- Use the SQL in other contexts

**Use schema functions** when you need to:
- Execute schema operations directly against the database
- Perform existence checks
- Retrieve metadata from the database
