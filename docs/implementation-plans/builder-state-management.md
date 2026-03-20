# Builder State Management Implementation Plan

**Priority:** MEDIUM  
**Target:** v0.19.0  
**Estimated Effort:** 1 week  
**Dependencies:** Phase 1 error handling complete

---

## Overview

Implement builder state management capabilities to enable safe query building patterns, prevent state contamination, and support complex query workflows with proper isolation.

---

## Problem Statement

Current implementation has no way to reset builder state or create independent builder instances. This leads to potential state contamination when reusing builders and limits the ability to build complex queries in parallel.

---

## Solution Design

Add Clone() and Reset() methods for state management, enabling safe builder reuse and independent query building workflows.

---

## Implementation Steps

### Step 1: Design State Management API

```go
// Add to BuilderInterface
Clone() BuilderInterface
Reset() BuilderInterface
GetState() BuilderState
SetState(state BuilderState) BuilderInterface

// State serialization for debugging/inspection
type BuilderState struct {
    Dialect            string
    TableName          string
    Columns            []Column
    Where              []Where
    Joins              []Join
    OrderBy            []OrderBy
    GroupBy            []GroupBy
    Limit              int
    Offset             int
    Errors             []error
    Parameters         []interface{}
    SelectColumns      []string
}
```

### Step 2: Implement Clone Method

```go
func (b *Builder) Clone() BuilderInterface {
    // Create new builder instance
    newBuilder := &Builder{
        Dialect:            b.Dialect,
        sql:                make(map[string]any{}),
        sqlColumns:         make([]Column, len(b.sqlColumns)),
        sqlErrors:          make([]error, len(b.sqlErrors)),
        sqlGroupBy:         make([]GroupBy, len(b.sqlGroupBy)),
        sqlJoins:           make([]Join, len(b.sqlJoins)),
        sqlLimit:           b.sqlLimit,
        sqlOffset:          b.sqlOffset,
        sqlOrderBy:         make([]OrderBy, len(b.sqlOrderBy)),
        sqlTableName:       b.sqlTableName,
        sqlViewName:        b.sqlViewName,
        sqlViewColumns:     make([]string, len(b.sqlViewColumns)),
        sqlViewSQL:         b.sqlViewSQL,
        sqlWhere:           make([]Where, len(b.sqlWhere)),
        columnSQLGenerator: b.columnSQLGenerator,
        params:             make([]interface{}, len(b.params)),
        paramIndex:          b.paramIndex,
        sqlSelectColumns:   make([]string, len(b.sqlSelectColumns)),
    }
    
    // Deep copy slices
    copy(newBuilder.sqlColumns, b.sqlColumns)
    copy(newBuilder.sqlErrors, b.sqlErrors)
    copy(newBuilder.sqlGroupBy, b.sqlGroupBy)
    copy(newBuilder.sqlJoins, b.sqlJoins)
    copy(newBuilder.sqlOrderBy, b.sqlOrderBy)
    copy(newBuilder.sqlWhere, b.sqlWhere)
    copy(newBuilder.params, b.params)
    copy(newBuilder.sqlSelectColumns, b.sqlSelectColumns)
    
    // Deep copy complex structures
    if b.sql != nil {
        for k, v := range b.sql {
            newBuilder.sql[k] = v
        }
    }
    
    // Deep copy sqlViewColumns
    copy(newBuilder.sqlViewColumns, b.sqlViewColumns)
    
    return newBuilder
}
```

### Step 3: Implement Reset Method

```go
func (b *Builder) Reset() BuilderInterface {
    // Clear all accumulated state
    b.sql = make(map[string]any{})
    b.sqlColumns = []Column{}
    b.sqlErrors = []error{}
    b.sqlGroupBy = []GroupBy{}
    b.sqlJoins = []Join{}
    b.sqlLimit = 0
    b.sqlOffset = 0
    b.sqlOrderBy = []OrderBy{}
    b.sqlTableName = ""
    b.sqlViewName = ""
    b.sqlViewColumns = []string{}
    b.sqlViewSQL = ""
    b.sqlWhere = []Where{}
    b.params = []interface{}{}
    b.paramIndex = 0
    b.sqlSelectColumns = []string{}
    
    return b
}
```

### Step 4: Implement State Serialization

```go
func (b *Builder) GetState() BuilderState {
    return BuilderState{
        Dialect:       b.Dialect,
        TableName:     b.sqlTableName,
        Columns:       b.deepCopyColumns(b.sqlColumns),
        Where:         b.deepCopyWhere(b.sqlWhere),
        Joins:         b.deepCopyJoins(b.sqlJoins),
        OrderBy:       b.deepCopyOrderBy(b.sqlOrderBy),
        GroupBy:       b.deepCopyGroupBy(b.sqlGroupBy),
        Limit:         b.sqlLimit,
        Offset:        b.sqlOffset,
        Errors:        b.deepCopyErrors(b.sqlErrors),
        Parameters:    b.deepCopyParams(b.params),
        SelectColumns: b.deepCopyStrings(b.sqlSelectColumns),
    }
}

func (b *Builder) SetState(state BuilderState) BuilderInterface {
    b.Dialect = state.Dialect
    b.sqlTableName = state.TableName
    b.sqlColumns = b.deepCopyColumns(state.Columns)
    b.sqlWhere = b.deepCopyWhere(state.Where)
    b.sqlJoins = b.deepCopyJoins(state.Joins)
    b.sqlOrderBy = b.deepCopyOrderBy(state.OrderBy)
    b.sqlGroupBy = b.deepCopyGroupBy(state.GroupBy)
    b.sqlLimit = state.Limit
    b.sqlOffset = state.Offset
    b.sqlErrors = b.deepCopyErrors(state.Errors)
    b.params = b.deepCopyParams(state.Parameters)
    b.sqlSelectColumns = b.deepCopyStrings(state.SelectColumns)
    
    return b
}

// Helper methods for deep copying
func (b *Builder) deepCopyColumns(cols []Column) []Column {
    result := make([]Column, len(cols))
    copy(result, cols)
    return result
}

func (b *Builder) deepCopyWhere(wheres []Where) []Where {
    result := make([]Where, len(wheres))
    for i, where := range wheres {
        result[i] = Where{
            Raw:      where.Raw,
            Column:   where.Column,
            Operator: where.Operator,
            Type:     where.Type,
            Value:    where.Value,
            IsNot:    where.IsNot,
        }
        
        // Deep copy children
        if len(where.Children) > 0 {
            result[i].Children = b.deepCopyWhere(where.Children)
        }
        
        // Handle subquery cloning
        if where.Subquery != nil {
            result[i].Subquery = where.Subquery.Clone().(*Builder)
        }
    }
    return result
}

func (b *Builder) deepCopyJoins(joins []Join) []Join {
    result := make([]Join, len(joins))
    copy(result, joins)
    return result
}

func (b *Builder) deepCopyOrderBy(orderBy []OrderBy) []OrderBy {
    result := make([]OrderBy, len(orderBy))
    copy(result, orderBy)
    return result
}

func (b *Builder) deepCopyGroupBy(groupBy []GroupBy) []GroupBy {
    result := make([]GroupBy, len(groupBy))
    copy(result, groupBy)
    return result
}

func (b *Builder) deepCopyErrors(errs []error) []error {
    result := make([]error, len(errs))
    copy(result, errs)
    return result
}

func (b *Builder) deepCopyParams(params []interface{}) []interface{} {
    result := make([]interface{}, len(params))
    copy(result, params)
    return result
}

func (b *Builder) deepCopyStrings(strs []string) []string {
    result := make([]string, len(strs))
    copy(result, strs)
    return result
}
```

### Step 5: Add Convenience Methods

```go
// Create a new builder with same dialect but reset state
func (b *Builder) NewBuilder() BuilderInterface {
    newBuilder := sb.NewBuilder(b.Dialect)
    return newBuilder
}

// Clone and reset in one operation
func (b *Builder) CloneAndReset() BuilderInterface {
    return b.Clone().Reset()
}

// Check if builder has any state
func (b *Builder) IsEmpty() bool {
    return b.sqlTableName == "" && 
           len(b.sqlColumns) == 0 && 
           len(b.sqlWhere) == 0 && 
           len(b.sqlJoins) == 0 && 
           len(b.sqlOrderBy) == 0 &&
           len(b.sqlGroupBy) == 0 &&
           b.sqlLimit == 0 &&
           b.sqlOffset == 0
}

// Get query summary for debugging
func (b *Builder) GetSummary() string {
    var summary strings.Builder
    
    if b.sqlTableName != "" {
        summary.WriteString("Table: ")
        summary.WriteString(b.sqlTableName)
        summary.WriteString("\n")
    }
    
    if len(b.sqlColumns) > 0 {
        summary.WriteString("Columns: ")
        for i, col := range b.sqlColumns {
            if i > 0 {
                summary.WriteString(", ")
            }
            summary.WriteString(col.Name)
        }
        summary.WriteString("\n")
    }
    
    if len(b.sqlWhere) > 0 {
        summary.WriteString("Where: ")
        summary.WriteString(strconv.Itoa(len(b.sqlWhere)))
        summary.WriteString(" conditions\n")
    }
    
    if len(b.sqlJoins) > 0 {
        summary.WriteString("Joins: ")
        summary.WriteString(strconv.Itoa(len(b.sqlJoins)))
        summary.WriteString("\n")
    }
    
    if len(b.sqlOrderBy) > 0 {
        summary.WriteString("OrderBy: ")
        summary.WriteString(strconv.Itoa(len(b.sqlOrderBy)))
        summary.WriteString("\n")
    }
    
    if b.sqlLimit > 0 {
        summary.WriteString("Limit: ")
        summary.WriteString(strconv.Itoa(b.sqlLimit))
        summary.WriteString("\n")
    }
    
    if len(b.sqlErrors) > 0 {
        summary.WriteString("Errors: ")
        summary.WriteString(strconv.Itoa(len(b.sqlErrors)))
        summary.WriteString("\n")
    }
    
    return summary.String()
}
```

### Step 6: Update Interface

```go
// Add to BuilderInterface
Clone() BuilderInterface
Reset() BuilderInterface
GetState() BuilderState
SetState(state BuilderState) BuilderInterface
NewBuilder() BuilderInterface
CloneAndReset() BuilderInterface
IsEmpty() bool
GetSummary() string
```

---

## Usage Patterns

### Pattern 1: Builder Reuse with Reset
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users")

// First query
sql1, _ := builder.Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    Select([]string{"name", "email"})

// Reset for second query
builder.Reset()
sql2, _ := builder.Where(&sb.Where{Column: "age", Operator: ">", Value: 18}).
    Select([]string{"name", "age"})
```

### Pattern 2: Parallel Query Building
```go
baseBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users")

// Clone for different queries
activeUsers := baseBuilder.Clone().Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})
recentUsers := baseBuilder.Clone().Where(&sb.Where{Column: "created_at", Operator: ">", Value: "2023-01-01"})

// Build queries independently
sql1, _ := activeUsers.Select([]string{"name"})
sql2, _ := recentUsers.Select([]string{"name"})
```

### Pattern 3: Complex Query Construction
```go
// Build base query
baseQuery := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    Join(sb.JOIN_INNER, "users", "orders.user_id = users.id")

// Create variations
baseQuery1 := baseQuery.Clone().Where(&sb.Where{Column: "status", Operator: "=", Value: "pending"})
baseQuery2 := baseQuery.Clone().Where(&sb.Where{Column: "status", Operator: "=", Value: "completed"})

// Add different ordering
sql1, _ := baseQuery1.OrderBy("created_at", "DESC").Select([]string{"*"})
sql2, _ := baseQuery2.OrderBy("total", "ASC").Select([]string{"*"})
```

### Pattern 4: State Inspection and Debugging
```go
builder := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
    OrderBy("name", "ASC")

// Inspect state
state := builder.GetState()
fmt.Printf("Table: %s\n", state.TableName)
fmt.Printf("Where conditions: %d\n", len(state.Where))

// Get summary
fmt.Println(builder.GetSummary())

// Clone and modify
clone := builder.Clone()
clone.Reset()
```

---

## Testing Strategy

### Unit Tests
- Clone method creates independent copies
- Reset method clears all state
- State serialization/deserialization
- Deep copy accuracy
- Memory leak prevention

### Integration Tests
- Complex query building patterns
- Parallel query construction
- Builder reuse scenarios
- State contamination prevention

### Performance Tests
- Clone performance with large queries
- Reset performance
- Memory usage patterns
- Thread safety validation

---

## Success Criteria

- [ ] Builder state management implemented
- [ ] Clone creates independent copies
- [ ] Reset clears all state
- [ ] No state contamination
- [ ] Thread safety maintained
- [ ] Comprehensive test coverage
- [ ] Documentation complete

---

## Release Notes

### v0.19.0 Major Release
- **Added:** Builder state management with Clone() and Reset()
- **Added:** State serialization for debugging
- **Added:** Convenience methods for common patterns
- **Improved:** Builder reuse and parallel query building
- **Breaking Changes:** None - additive only

---

## Rollback Plan

If issues arise:
1. Remove state management methods
2. Keep existing builder behavior
3. Revert to original API

---

## Timeline

- **Week 1:** Core Clone and Reset implementation
- **Week 2:** State serialization and convenience methods
- **Week 3:** Testing and documentation
- **Total:** 1 week for core implementation

---

## Resources Required

### Development
- **1 developer** for 1 week
- **Testing expertise** for thread safety validation
- **Performance testing** for optimization

### Testing
- **CI/CD pipeline** for automated testing
- **Performance testing** tools
- **Thread safety** testing frameworks

---

## Conclusion

Builder state management will significantly improve the developer experience by enabling safe builder reuse, parallel query construction, and better debugging capabilities while maintaining the library's simplicity and performance.
