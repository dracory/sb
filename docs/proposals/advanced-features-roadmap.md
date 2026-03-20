# Advanced SQL Features Roadmap

**Last updated at:** 2026-03-20  
**Status:** Roadmap  
**Priority:** Low

## Overview

This roadmap outlines advanced SQL features that could be added to SB beyond the current scope. These features would enhance the library's capabilities while maintaining its simplicity philosophy.

## Current Scope vs Future Enhancements

### Current Scope (Maintained)
- ✅ Basic CRUD operations
- ✅ Table and column management
- ✅ Simple index operations
- ✅ View creation
- ✅ Transaction support
- ✅ Multi-dialect support
- ✅ JOIN operations
- ✅ Subquery support

### Future Scope (This Roadmap)
- 🔄 Advanced query building
- 🔄 Complex JOIN operations
- 🔄 Advanced index types
- 🔄 Schema introspection
- 🔄 Migration tools

## Feature Categories

### 1. Advanced Query Building

#### 1.1 Subquery Support ✅ COMPLETED
**Priority:** High  
**Complexity:** Medium  
**Status:** ✅ Completed 2026-03-20

```go
// Subquery in WHERE clause
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    InSubquery(subquery).
    Select([]string{"name"})
// Result: SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > "1000");

// EXISTS subquery
activeOrders := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    Exists(activeOrders).
    Select([]string{"name", "email"})
// Result: SELECT "name", "email" FROM "users" WHERE EXISTS (SELECT * FROM "orders" WHERE "status" = "active");

// Correlated subquery
subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"})

sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{
        Column:   "id",
        Operator: ">",
        Subquery: subquery.(*sb.Builder),
    }).
    Select([]string{"name"})
// Result: SELECT `name` FROM `users` WHERE `id` > (SELECT * FROM `orders` WHERE `user_id` = "users.id");
```

**Implementation Completed:**
- ✅ Extend Where struct to support subqueries
- ✅ Add subquery SQL generation with dialect-specific quoting
- ✅ Handle IN, NOT IN, EXISTS, NOT EXISTS, and comparison operators
- ✅ Support correlated subqueries
- ✅ Test across all dialects (16 comprehensive tests)

#### 1.2 UNION Operations
**Priority:** Medium  
**Complexity:** Low

```go
builder1 := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("active_users").
    Select([]string{"name", "email"})

builder2 := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("inactive_users").
    Select([]string{"name", "email"})

sql := builder1.Union(builder2)
// Result: SELECT name, email FROM active_users UNION SELECT name, email FROM inactive_users
```

#### 1.3 Window Functions
**Priority:** Low  
**Complexity:** High

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("sales").
    Select([]string{
        "product_id",
        "sale_amount",
        "ROW_NUMBER() OVER (PARTITION BY product_id ORDER BY sale_date DESC) as rank",
    }).
    Where(&sb.Where{
        Column: "rank",
        Operator: "<=",
        Value: 3,
    })
```

### 2. Advanced JOIN Operations

#### 2.1 Complex JOIN Conditions
**Priority:** High  
**Complexity:** Medium

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    InnerJoin("users", "orders.user_id = users.id AND users.status = ?", "active").
    LeftJoin("profiles", "users.id = profiles.user_id").
    Select([]string{"orders.*", "users.name", "profiles.avatar"})
```

#### 2.2 Self-Joins
**Priority:** Medium  
**Complexity:** Low

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("employees").
    InnerJoin("employees AS manager", "employees.manager_id = manager.id").
    Select([]string{"employees.name", "manager.name as manager_name"})
```

#### 2.3 USING Clause Support
**Priority:** Low  
**Complexity:** Low

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    InnerJoinUsing("users", "user_id").
    Select([]string{"orders.total", "users.name"})
// Result: SELECT orders.total, users.name FROM orders INNER JOIN users USING (user_id)
```

### 3. Advanced Index Management

#### 3.1 Partial Indexes
**Priority:** Medium  
**Complexity:** Medium

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("orders").
    CreatePartialIndex("idx_active_orders", "user_id", "status = 'active'")
// Result: CREATE INDEX "idx_active_orders" ON "orders" ("user_id") WHERE (status = 'active')
```

#### 3.2 Functional Indexes
**Priority:** Low  
**Complexity:** High

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    CreateFunctionalIndex("idx_users_lower_email", "LOWER(email)")
// Result: CREATE INDEX "idx_users_lower_email" ON "users" (LOWER(email))
```

#### 3.3 Composite Indexes with Ordering
**Priority:** Medium  
**Complexity:** Low

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("events").
    CreateIndexWithDirection("idx_events_user_time", []IndexColumn{
        {Name: "user_id", Direction: "ASC"},
        {Name: "created_at", Direction: "DESC"},
    })
```

### 4. Schema Introspection

#### 4.1 Table Information
**Priority:** Medium  
**Complexity:** Low

```go
func (b *Builder) GetTableInfo(tableName string) (*TableInfo, error) {
    // Return table metadata: columns, indexes, constraints, etc.
}
```

#### 4.2 Column Information
**Priority:** Medium  
**Complexity:** Low

```go
func (b *Builder) GetColumnInfo(tableName string) ([]ColumnInfo, error) {
    // Return detailed column information
}
```

#### 4.3 Index Information
**Priority:** Low  
**Complexity:** Medium

```go
func (b *Builder) GetIndexInfo(tableName string) ([]IndexInfo, error) {
    // Return index metadata
}
```

### 5. Migration Tools

#### 5.1 Schema Migration
**Priority:** Low  
**Complexity:** High

```go
type Migration struct {
    Version     int
    Description string
    Up          func(*Builder) []string
    Down        func(*Builder) []string
}

func (m *MigrationManager) RunMigration(db DatabaseInterface, migration Migration) error {
    // Execute migration with transaction support
}
```

#### 5.2 Schema Diff
**Priority:** Low  
**Complexity:** High

```go
func (b *Builder) SchemaDiff(targetSchema Schema) ([]string, error) {
    // Generate SQL to transform current schema to target schema
}
```

### 6. Advanced Data Types

#### 6.1 JSON Support
**Priority:** Medium  
**Complexity:** Medium

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    Where(&sb.Where{
        Column: "metadata->>'role'",
        Operator: "=",
        Value: "admin",
    }).
    Select([]string{"name"})
```

#### 6.2 Array Support
**Priority:** Low  
**Complexity:** Medium

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("posts").
    Where(&sb.Where{
        Column: "tags",
        Operator: "@>",
        Value: []string{"golang", "database"},
    }).
    Select([]string{"title"})
```

#### 6.3 Spatial Types
**Priority:** Low  
**Complexity:** High

```go
sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("locations").
    Where(&sb.Where{
        Column: "coordinates",
        Operator: "<@>",
        Value: "POINT(0 0)",
    }).
    Select([]string{"name"})
```

### 7. Performance and Optimization

#### 7.1 Query Explain
**Priority:** Medium  
**Complexity:** Low

```go
func (b *Builder) Explain() string {
    switch b.Dialect {
    case DIALECT_MYSQL:
        return "EXPLAIN " + b.Select([]string{"*"})
    case DIALECT_POSTGRES:
        return "EXPLAIN ANALYZE " + b.Select([]string{"*"})
    // ... other dialects
    }
}
```

#### 7.2 Query Hints
**Priority:** Low  
**Complexity:** Medium

```go
sql := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    WithHint("USE INDEX (idx_orders_user_id)").
    Select([]string{"*"})
```

### 8. Advanced Transaction Management

#### 8.1 Savepoints
**Priority:** Low  
**Complexity:** Medium

```go
func (db *Database) CreateSavepoint(name string) error
func (db *Database) RollbackToSavepoint(name string) error
func (db *Database) ReleaseSavepoint(name string) error
```

#### 8.2 Transaction Isolation Levels
**Priority:** Low  
**Complexity:** Low

```go
func (db *Database) BeginTransactionWithIsolation(level IsolationLevel) error
```

## Implementation Strategy

### Phase 1: Essential Enhancements (3-6 months) ✅ PARTIALLY COMPLETED
1. ✅ **JOIN Support** - Critical for real-world queries (Completed 2026-03-19)
2. ✅ **TRUNCATE Implementation** - Complete basic operations (Completed 2026-03-19)
3. ✅ **DropIndex Implementation** - Complete index lifecycle (Completed 2026-03-19)
4. ✅ **Basic Subqueries** - Enhance query capabilities (Completed 2026-03-20)

### Phase 2: Advanced Features (6-12 months) 🔄 NEXT
1. **UNION Operations** - Expand query building
2. **Schema Introspection** - Better database awareness
3. **Advanced Index Types** - Performance optimization
4. **Query Explain** - Performance debugging

### Phase 3: Specialized Features (12-18 months) 🔮 FUTURE
1. **Migration Tools** - Schema management
2. **Window Functions** - Advanced analytics
3. **JSON/Array Support** - Modern data types
4. **Advanced Transactions** - Complex workflows

## Decision Criteria

### Include Feature If:
- Aligns with simplicity philosophy
- Has broad applicability
- Maintains backward compatibility
- Works across all supported dialects
- Has clear performance benefits

### Exclude Feature If:
- Overly complex for limited use cases
- Breaks existing API
- Only works with specific databases
- Duplicates existing functionality
- Requires heavy dependencies

## Resource Requirements

### Development Resources
- **Senior Go Developer**: 0.5 FTE for Phase 1
- **Database Expert**: 0.25 FTE for dialect-specific features
- **QA Engineer**: 0.25 FTE for comprehensive testing

### Infrastructure
- **Test Databases**: MySQL, PostgreSQL, SQLite, MSSQL instances
- **CI/CD**: Enhanced testing matrix
- **Documentation**: Updated API docs and examples

### Timeline Estimates
- **Phase 1**: 3-6 months (parallel development possible)
- **Phase 2**: 6-12 months (depends on Phase 1 completion)
- **Phase 3**: 12-18 months (resource-dependent)

## Success Metrics

### Technical Metrics
- **Test Coverage**: >90% for new features
- **Performance**: <5% overhead over existing features
- **Compatibility**: 100% backward compatibility
- **Dialect Support**: Consistent across all databases

### User Metrics
- **Adoption Rate**: Usage of new features in community
- **Bug Reports**: <5 critical bugs per feature
- **Documentation**: Clear examples and guides
- **Community Feedback**: Positive response to enhancements

## Risk Assessment

### Technical Risks
- **Complexity Creep**: Features becoming too complex
- **Dialect Inconsistencies**: Different database behaviors
- **Performance Regression**: Impact on existing operations
- **Maintenance Burden**: Increased code complexity

### Mitigation Strategies
- **Feature Flags**: Enable/disable advanced features
- **Phased Rollout**: Gradual feature introduction
- **Comprehensive Testing**: Prevent regressions
- **Community Feedback**: Early and often user input

## Conclusion

This roadmap provides a structured approach to enhancing SB while maintaining its core simplicity philosophy. The phased implementation allows for incremental delivery and continuous user feedback. Each feature is evaluated against the library's core principles to ensure consistent user experience and maintainable codebase.

The focus remains on practical, broadly useful features that enhance the library's capabilities without compromising its simplicity and performance characteristics.
