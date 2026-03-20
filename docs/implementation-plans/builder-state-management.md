# Builder State Management Implementation Plan

**Priority:** MEDIUM  
**Target:** v0.20.0  
**Estimated Effort:** 1 week  
**Dependencies:** Phase 1 error handling complete

---

## Overview

Implement builder state management functionality to enhance developer experience and enable better query builder reuse patterns. **Note:** Zero-panic error handling system has been completed in Phase 1, providing a solid foundation for all future development.

**Previous Phase 2 components completed:**
- ✅ **Parameterized Queries** - Security enhancement with prepared statements (v0.18.0)
- ✅ **Enhanced Index Support** - Advanced index management capabilities (v0.19.0)

---

## Implementation Goals

### Core Functionality
- **Clone() method** - Create independent copies of builder instances
- **Reset() method** - Clear builder state for reuse
- **State isolation** - Prevent interference between query builders
- **Memory management** - Efficient state handling

### Developer Experience
- **Builder reuse patterns** - Enable safe builder recycling
- **Query composition** - Support complex query building workflows
- **Debugging support** - State inspection capabilities

---

## Implementation Details

### Clone() Method
```go
// Clone creates an independent copy of the builder with all current state
func (b *Builder) Clone() BuilderInterface

// Usage Example
original := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Where(&sb.Where{Column: "active", Operator: "=", Value: true})
clone := original.Clone()

// Independent operations
sql1, _ := original.Select([]string{"name"})
sql2, _ := clone.Where(&sb.Where{Column: "id", Operator: ">", Value: 100}).Select([]string{"*"})
```

### Reset() Method
```go
// Reset clears all builder state while preserving dialect and configuration
func (b *Builder) Reset() BuilderInterface

// Usage Example
builder := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Where(&sb.Where{Column: "active", Operator: "=", Value: true})
sql1, _ := builder.Select([]string{"name"})

// Reset and reuse
builder.Reset()
builder.Table("orders").Where(&sb.Where{Column: "status", Operator: "=", Value: "pending"})
sql2, _ := builder.Select([]string{"*"})
```

### State Management Features
- **Deep copy** of all builder state
- **Independent error collection** per builder instance
- **Parameter isolation** for parameterized queries
- **Configuration preservation** (dialect, settings)

---

## Technical Implementation

### State Copy Strategy
```go
type Builder struct {
    Dialect            string
    sql                map[string]any
    sqlColumns         []Column
    sqlErrors          []error
    sqlGroupBy         []GroupBy
    sqlJoins           []Join
    sqlLimit           int64
    sqlOffset          int64
    sqlOrderBy         []OrderBy
    sqlTable           string
    sqlWhere           []Where
    // ... other fields
}

// Clone implementation
func (b *Builder) Clone() BuilderInterface {
    clone := &Builder{
        Dialect:            b.Dialect,
        sql:                make(map[string]any),
        sqlColumns:         make([]Column, len(b.sqlColumns)),
        sqlErrors:          make([]error, len(b.sqlErrors)),
        sqlGroupBy:         make([]GroupBy, len(b.sqlGroupBy)),
        sqlJoins:           make([]Join, len(b.sqlJoins)),
        sqlLimit:           b.sqlLimit,
        sqlOffset:          b.sqlOffset,
        sqlOrderBy:         make([]OrderBy, len(b.sqlOrderBy)),
        sqlTable:           b.sqlTable,
        sqlWhere:           make([]Where, len(b.sqlWhere)),
        // ... deep copy all fields
    }
    
    // Copy all slices and maps
    copy(clone.sqlColumns, b.sqlColumns)
    copy(clone.sqlErrors, b.sqlErrors)
    // ... copy remaining fields
    
    return clone
}
```

### Reset Implementation
```go
func (b *Builder) Reset() BuilderInterface {
    // Preserve core configuration
    dialect := b.Dialect
    
    // Reset all mutable state
    b.sql = make(map[string]any)
    b.sqlColumns = []Column{}
    b.sqlErrors = []error{}
    b.sqlGroupBy = []GroupBy{}
    b.sqlJoins = []Join{}
    b.sqlLimit = 0
    b.sqlOffset = 0
    b.sqlOrderBy = []OrderBy{}
    b.sqlTable = ""
    b.sqlWhere = []Where{}
    
    // Reset parameter tracking
    b.resetParams()
    
    return b
}
```

---

## Testing Strategy

### Unit Tests
- **Clone isolation tests** - Verify independent state changes
- **Reset functionality tests** - Ensure proper state clearing
- **State preservation tests** - Verify configuration kept during reset
- **Deep copy validation** - Ensure complete state duplication

### Integration Tests
- **Query composition patterns** - Test complex builder workflows
- **Parallel operations** - Verify thread safety of cloned builders
- **Memory usage tests** - Validate efficient state management

### Performance Tests
- **Clone performance** - Ensure cloning is efficient
- **Reset performance** - Verify fast state clearing
- **Memory footprint** - Monitor memory usage patterns

---

## Success Criteria

- [ ] Clone() method creates independent builder copies
- [ ] Reset() method clears state while preserving configuration
- [ ] State isolation prevents interference between builders
- [ ] Memory usage remains efficient
- [ ] Thread safety for concurrent operations
- [ ] Backward compatibility maintained

---

## API Design

### Interface Updates
```go
type BuilderInterface interface {
    // ... existing methods
    
    // Clone creates an independent copy of the builder
    Clone() BuilderInterface
    
    // Reset clears all builder state while preserving configuration
    Reset() BuilderInterface
}
```

### Usage Patterns
```go
// Pattern 1: Query composition
base := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users")

activeUsers := base.Clone().Where(&sb.Where{Column: "active", Operator: "=", Value: true})
recentUsers := base.Clone().Where(&sb.Where{Column: "created_at", Operator: ">", Value: "2023-01-01"})

// Pattern 2: Builder reuse
builder := sb.NewBuilder(sb.DIALECT_MYSQL)

// First query
builder.Table("users").Where(&sb.Where{Column: "id", Operator: "=", Value: 1})
sql1, _ := builder.Select([]string{"name"})

// Reset and reuse
builder.Reset()
builder.Table("orders").Where(&sb.Where{Column: "user_id", Operator: "=", Value: 1})
sql2, _ := builder.Select([]string{"total"})

// Pattern 3: Parallel operations
base := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users")

var wg sync.WaitGroup
var results []string

for _, id := range ids {
    wg.Add(1)
    go func(userID int) {
        defer wg.Done()
        
        builder := base.Clone()
        builder.Where(&sb.Where{Column: "id", Operator: "=", Value: userID})
        sql, _ := builder.Select([]string{"name"})
        
        results = append(results, sql)
    }(id)
}

wg.Wait()
```

---

## Error Handling Integration

### Error Collection
- **Independent error collection** per builder instance
- **Error isolation** - errors in clone don't affect original
- **Error propagation** - maintain existing error handling patterns

### Parameter Management
- **Parameter isolation** for parameterized queries
- **Independent parameter tracking** per builder
- **Parameter cleanup** on reset

---

## Implementation Phases

### Phase 1: Core Methods (Week 1)
- Implement Clone() method with deep copy
- Implement Reset() method with state clearing
- Add basic unit tests

### Phase 2: Advanced Features (Week 2)
- Add state serialization support
- Implement debugging helpers
- Add comprehensive integration tests

### Phase 3: Optimization (Week 3)
- Performance optimization for cloning
- Memory usage improvements
- Thread safety validation

---

## Release Notes

### v0.20.0 Minor Release
- **Added:** Clone() method for builder state isolation
- **Added:** Reset() method for builder reuse
- **Enhanced:** Query composition capabilities
- **Improved:** Memory management for builder instances

---

## Rollback Plan

If issues arise:
1. Disable Clone() method (return original builder)
2. Disable Reset() method (no-op implementation)
3. Revert to original state management
4. Maintain backward compatibility

---

## Conclusion

Builder state management will significantly enhance the SB SQL builder library's developer experience by enabling safe query composition, builder reuse, and state isolation. The implementation builds upon the zero-panic foundation established in Phase 1 and maintains consistency with existing error handling patterns.

**Estimated Timeline:** 1 week  
**Priority:** MEDIUM  
**Risk:** LOW  
**Dependencies:** Phase 1 completion ✅ **MET**
