# Table Operations

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

## Table Existence Check

Check if a table exists in the database.

```go
exists := sb.TableColumnExists(db, "users", "id")
```
