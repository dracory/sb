# Implementation Guidelines

**Last updated at:** 2026-03-20  
**Status:** Guidelines  
**Priority:** Low

## Overview

This document provides implementation guidelines for contributing to the SB SQL builder library. These guidelines ensure consistency, maintainability, and adherence to the library's core principles.

## Recent Implementation Success: Subquery Support

### Implementation Summary
- **Feature**: Complete subquery support (IN, NOT IN, EXISTS, NOT EXISTS, comparisons)
- **Date Completed**: 2026-03-20
- **Test Coverage**: 16 comprehensive tests across all dialects
- **Backward Compatibility**: 100% maintained
- **Documentation**: Standard Go documentation with examples

### Key Implementation Patterns Demonstrated

#### 1. Interface-First Design
```go
// Interface definition in interfaces.go
Exists(subquery BuilderInterface) BuilderInterface
InSubquery(subquery BuilderInterface) BuilderInterface

// Implementation with type assertion
func (b *Builder) Exists(subquery BuilderInterface) BuilderInterface {
    if subquery == nil {
        panic("subquery cannot be nil")
    }
    where := Where{
        Operator: "EXISTS",
        Subquery: subquery.(*Builder),
        IsNot:    false,
    }
    b.sqlWhere = append(b.sqlWhere, where)
    return b
}
```

#### 2. Enhanced Data Structures
```go
// Extended existing struct with new fields
type Where struct {
    Raw       string
    Column    string
    Operator  string
    Type      string
    Value     string
    Subquery  *Builder  // New field for subquery support
    IsNot     bool      // New field for NOT operators
    Children  []Where
}

// Added field to Builder for subquery support
type Builder struct {
    // ... existing fields ...
    sqlSelectColumns []string // New field to store select columns for subqueries
}
```

#### 3. Centralized SQL Generation
```go
// New method for handling subquery SQL generation
func (b *Builder) whereToSqlSubquery(where Where) string {
    // Get columns and generate subquery SQL
    columns := where.Subquery.sqlSelectColumns
    if len(columns) == 0 {
        columns = []string{"*"}
    }
    
    subquerySQL := where.Subquery.Select(columns)
    subquerySQL = strings.TrimSuffix(subquerySQL, ";")
    
    // Handle different operators with dialect-specific quoting
    switch where.Operator {
    case "EXISTS":
        return "EXISTS (" + subquerySQL + ")"
    case "IN":
        columnQuoted := b.quoteColumn(where.Column)
        return columnQuoted + " IN (" + subquerySQL + ")"
    // ... other cases
    }
}
```

#### 4. Comprehensive Multi-Dialect Testing
```go
func TestBuilderSubqueryInMySQL(t *testing.T) {
    subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("orders").
        Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})
    
    sql := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        InSubquery(subquery).
        Select([]string{"name"})
    
    expected := "SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > \"1000\");"
    if sql != expected {
        t.Fatalf("Expected: %s but found: %s", expected, sql)
    }
}
```

### Lessons Learned for Future Implementations

1. **Always test all dialects early** - Database-specific quoting differences can cause subtle bugs
2. **Interface design matters** - Use interface types in public signatures, handle type assertions safely
3. **Incremental development works** - Start with basic functionality, add complexity gradually
4. **Documentation drives quality** - Write examples before implementation, document edge cases

## Core Principles

### 1. Simplicity First
- Keep APIs simple and intuitive
- Avoid over-engineering solutions
- Prioritize common use cases over edge cases
- Maintain fluent, readable method chaining

### 2. Database Dialect Consistency
- Ensure consistent behavior across MySQL, PostgreSQL, SQLite, and MSSQL
- Use dialect-specific generators when necessary
- Test all implementations against all supported databases

### 3. Backward Compatibility
- Never break existing APIs
- Additive changes only
- Maintain existing method signatures
- Provide migration paths for breaking changes (if absolutely necessary)

### 4. Type Safety
- Use strong typing for SQL components
- Avoid string concatenation for SQL building
- Validate inputs at compile time when possible

## Code Standards

### 1. File Organization

#### Core Files
```
builder.go           # Main builder implementation
interfaces.go        # Interface definitions
database.go          # Database wrapper implementation
consts.go           # Constants and enums
builder_*.go         # Feature-specific implementations
```

#### Naming Conventions
```go
// Types: PascalCase
type Builder struct { ... }
type Column struct { ... }
type Where struct { ... }

// Methods: PascalCase (public), camelCase (private)
func (b *Builder) Table(name string) BuilderInterface { ... }
func (b *Builder) quoteColumn(name string) string { ... }

// Constants: UPPER_SNAKE_CASE
const DIALECT_MYSQL = "mysql"
const COLUMN_TYPE_STRING = "string"

// Variables: camelCase
var sqlTableName string
var columnSQLGenerator ColumnSQLGenerator
```

### 2. Method Implementation Patterns

#### Fluent Interface Pattern
```go
func (b *Builder) TableName(name string) BuilderInterface {
    b.sqlTableName = name
    return b
}
```

#### Error Handling Pattern
```go
func (b *Builder) SomeOperation() (result string, err error) {
    if b.sqlTableName == "" {
        return "", errors.New("no table specified")
    }
    
    switch b.Dialect {
    case DIALECT_MYSQL:
        return b.generateMySQLSQL(), nil
    case DIALECT_POSTGRES:
        return b.generatePostgreSQLSQL(), nil
    default:
        return "", errors.New("unsupported dialect: " + b.Dialect)
    }
}
```

#### Validation Pattern
```go
func (b *Builder) validateInputs() error {
    if b.sqlTableName == "" {
        return errors.New("table name required")
    }
    if len(b.sqlColumns) == 0 {
        return errors.New("at least one column required")
    }
    return nil
}
```

### 3. SQL Generation Patterns

#### Dialect-Specific Generation
```go
func (b *Builder) generateCreateTableSQL() string {
    switch b.Dialect {
    case DIALECT_MYSQL:
        return b.generateMySQLCreateTable()
    case DIALECT_POSTGRES:
        return b.generatePostgreSQLCreateTable()
    case DIALECT_SQLITE:
        return b.generateSQLiteCreateTable()
    case DIALECT_MSSQL:
        return b.generateMSSQLCreateTable()
    default:
        panic("unsupported dialect: " + b.Dialect)
    }
}
```

#### Quote Handling
```go
func (b *Builder) quoteTable(name string) string {
    return b.quote(name, "table")
}

func (b *Builder) quoteColumn(name string) string {
    return b.quote(name, "column")
}

func (b *Builder) quote(name string, quoteType string) string {
    switch b.Dialect {
    case DIALECT_MYSQL:
        return "`" + name + "`"
    case DIALECT_POSTGRES, DIALECT_SQLITE:
        return `"` + name + `"`
    case DIALECT_MSSQL:
        return "[" + name + "]"
    default:
        return name
    }
}
```

## Testing Guidelines

### 1. Test Structure

#### Unit Test Pattern
```go
func TestBuilderCreateTableMySQL(t *testing.T) {
    // Arrange
    builder := sb.NewBuilder(sb.DIALECT_MYSQL).
        Table("users").
        Column(sb.Column{
            Name: "id",
            Type: sb.COLUMN_TYPE_STRING,
            Length: 40,
            PrimaryKey: true,
        })
    
    // Act
    sql := builder.Create()
    
    // Assert
    expected := "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL);"
    if sql != expected {
        t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
    }
}
```

#### Multi-Dialect Test Pattern
```go
func TestBuilderCreateTableAllDialects(t *testing.T) {
    tests := []struct {
        name     string
        dialect  string
        expected string
    }{
        {"MySQL", sb.DIALECT_MYSQL, "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL);"},
        {"PostgreSQL", sb.DIALECT_POSTGRES, `CREATE TABLE "users"("id" TEXT PRIMARY KEY NOT NULL);`},
        {"SQLite", sb.DIALECT_SQLITE, `CREATE TABLE "users"("id" TEXT(40) PRIMARY KEY NOT NULL);`},
        {"MSSQL", sb.DIALECT_MSSQL, "CREATE TABLE [users]([id] NVARCHAR(40) PRIMARY KEY NOT NULL);"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            sql := sb.NewBuilder(tt.dialect).
                Table("users").
                Column(sb.Column{
                    Name: "id",
                    Type: sb.COLUMN_TYPE_STRING,
                    Length: 40,
                    PrimaryKey: true,
                }).
                Create()
            
            if sql != tt.expected {
                t.Fatal("Expected:\n", tt.expected, "\nbut found:\n", sql)
            }
        })
    }
}
```

### 2. Test Coverage Requirements

#### Minimum Coverage
- **Core functionality**: 100% coverage
- **Dialect-specific code**: 100% coverage
- **Error handling**: 100% coverage
- **Edge cases**: >90% coverage

#### Test Categories
```go
// Basic functionality tests
func TestBuilderBasicOperation(t *testing.T) { ... }

// Dialect-specific tests
func TestBuilderMySQLSpecific(t *testing.T) { ... }
func TestBuilderPostgreSQLSpecific(t *testing.T) { ... }

// Error handling tests
func TestBuilderErrorHandling(t *testing.T) { ... }

// Integration tests
func TestBuilderIntegration(t *testing.T) { ... }
```

### 3. Integration Testing

#### Database Setup
```go
func initTestDB(dialect string) (*sb.Database, func()) {
    // Setup test database based on dialect
    // Return cleanup function
}

func TestBuilderWithRealDatabase(t *testing.T) {
    db, cleanup := initTestDB(sb.DIALECT_MYSQL)
    defer cleanup()
    
    // Test actual SQL execution
}
```

## Documentation Guidelines

### 1. Code Documentation

#### Function Documentation
```go
// CreateIndex creates an index on the table with the specified columns.
//
// Example:
//   sql := builder.Table("users").CreateIndex("idx_users_email", "email")
//
// Parameters:
//   - indexName: The name of the index to create
//   - columnName: One or more column names to include in the index
//
// Returns:
//   - string: The CREATE INDEX SQL statement
func (b *Builder) CreateIndex(indexName string, columnName ...string) string {
    // Implementation
}
```

#### Type Documentation
```go
// Column represents a database column with its properties.
//
// Fields:
//   - Name: The column name
//   - Type: The column data type (e.g., COLUMN_TYPE_STRING)
//   - Length: The column length/precision (optional)
//   - Decimals: The number of decimal places (for numeric types)
//   - AutoIncrement: Whether the column auto-increments
//   - PrimaryKey: Whether the column is a primary key
//   - Nullable: Whether the column allows NULL values
//   - Unique: Whether the column has a UNIQUE constraint
//   - Default: The default value for the column
type Column struct {
    // Field definitions...
}
```

### 2. Example Documentation

#### Usage Examples
```go
// Example: Create a table with multiple columns
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Column(sb.Column{
        Name: "id",
        Type: sb.COLUMN_TYPE_STRING,
        Length: 40,
        PrimaryKey: true,
    }).
    Column(sb.Column{
        Name: "email",
        Type: sb.COLUMN_TYPE_STRING,
        Length: 255,
        Unique: true,
    }).
    Create()
// Result: CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `email` VARCHAR(255) NOT NULL UNIQUE);
```

### 3. README Updates

When adding new features:
1. Update the feature list
2. Add usage examples
3. Document any breaking changes
4. Update installation instructions if needed

## Performance Guidelines

### 1. SQL Generation Performance

#### Efficient String Building
```go
// Good: Use strings.Builder for complex SQL
func (b *Builder) generateComplexSQL() string {
    var sql strings.Builder
    sql.WriteString("SELECT ")
    sql.WriteString(strings.Join(columns, ", "))
    sql.WriteString(" FROM ")
    sql.WriteString(b.quoteTable(tableName))
    return sql.String()
}

// Avoid: Repeated string concatenation
func (b *Builder) generateInefficientSQL() string {
    sql := "SELECT "
    sql += strings.Join(columns, ", ")
    sql += " FROM "
    sql += b.quoteTable(tableName)
    return sql
}
```

#### Memory Efficiency
```go
// Good: Reuse builders when possible
func (b *Builder) BatchGenerate(queries []Query) []string {
    results := make([]string, 0, len(queries))
    for _, query := range queries {
        // Reuse builder state efficiently
        sql := b.generateSQL(query)
        results = append(results, sql)
    }
    return results
}
```

### 2. Database Performance

#### Connection Management
```go
// Good: Use connection pooling
db := sb.NewDatabaseFromDriver("mysql", dsn)
defer db.Close()

// Good: Use transactions for multiple operations
err := db.ExecInTransaction(ctx, func(tx *sb.Database) error {
    // Multiple operations in single transaction
    return nil
})
```

## Security Guidelines

### 1. SQL Injection Prevention

#### Input Validation
```go
func (b *Builder) validateTableName(name string) error {
    // Check for dangerous characters
    if strings.ContainsAny(name, ";'\"\\") {
        return errors.New("invalid table name")
    }
    return nil
}
```

#### Proper Quoting
```go
// Always quote identifiers to prevent injection
func (b *Builder) quoteIdentifier(name string) string {
    // Implement dialect-specific quoting
    return b.quote(name, "identifier")
}
```

### 2. Parameter Binding

#### Support for Prepared Statements
```go
// Future enhancement: Support parameterized queries
type ParameterizedQuery struct {
    SQL    string
    Params []interface{}
}

func (b *Builder) SelectWithParams(columns []string) *ParameterizedQuery {
    // Generate parameterized SQL
    return &ParameterizedQuery{
        SQL:    sql,
        Params: params,
    }
}
```

## Release Process

### 1. Version Management

#### Semantic Versioning
- **Major (X.0.0)**: Breaking changes
- **Minor (X.Y.0)**: New features, backward compatible
- **Patch (X.Y.Z)**: Bug fixes only

#### Release Checklist
- [ ] All tests pass
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
- [ ] Version bumped in go.mod
- [ ] Tag created in Git
- [ ] GitHub release created

### 2. Quality Assurance

#### Pre-Release Testing
```bash
# Run all tests
go test ./...

# Test with race detector
go test -race ./...

# Test coverage
go test -cover ./...

# Benchmark performance
go test -bench=. ./...
```

#### Multi-Database Testing
```bash
# Test against all supported databases
make test-mysql
make test-postgres
make test-sqlite
make test-mssql
```

## Contribution Process

### 1. Development Workflow

#### Branch Strategy
```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes
# ... implement feature ...

# Test changes
go test ./...

# Submit pull request
git push origin feature/new-feature
```

#### Code Review Checklist
- [ ] Code follows style guidelines
- [ ] Tests are comprehensive
- [ ] Documentation is updated
- [ ] No breaking changes
- [ ] Performance impact assessed

### 2. Pull Request Requirements

#### PR Template
```markdown
## Description
Brief description of the change

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] All tests pass

## Checklist
- [ ] Code follows guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] CHANGELOG.md updated
```

## Conclusion

These guidelines ensure that SB maintains its high quality standards while continuing to evolve. Following these principles helps create a consistent, maintainable, and reliable SQL builder library that serves the Go community effectively.

Remember: **Simplicity is the ultimate sophistication** - always prioritize the user experience and the library's core philosophy of providing clean, intuitive SQL building capabilities.
