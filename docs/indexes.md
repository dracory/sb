# Index Management

## Overview

SB SQL builder provides comprehensive index management capabilities with database-specific features while maintaining simplicity and consistency across all supported databases.

## Enhanced Index Support (v0.19.0+)

The enhanced index support provides advanced index creation and management capabilities with database-specific features:

### Key Features
- **Advanced index types**: UNIQUE, partial, composite, covering indexes
- **Database-specific features**: FULLTEXT, GIN, partial indexes, etc.
- **Convenience methods**: Simplified API for common patterns
- **Comprehensive testing**: 32 tests covering all dialects and features

### Quick Examples

#### Basic Index Creation
```go
// Simple unique index
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("users").
    CreateUniqueIndex("idx_users_email", "email")
```

#### Advanced PostgreSQL Index
```go
// GIN index with covering columns and partial filter
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("documents").
    CreateIndexWithOptions("idx_search", sb.IndexOptions{
        Unique:      true,
        IfNotExists: true,
        Using:       sb.INDEX_TYPE_GIN,
        Columns:     []sb.IndexColumn{{Name: "search_vector"}},
        Include:     []string{"title", "content"},
        Where:       "published = true",
        Storage:     "fillfactor=90",
    })
```

#### MySQL FULLTEXT Index
```go
// Advanced MySQL FULLTEXT index with prefix and comment
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("articles").
    CreateIndexWithOptions("idx_content", sb.IndexOptions{
        Using:   sb.INDEX_TYPE_FULLTEXT,
        Columns: []sb.IndexColumn{
            {Name: "title", Length: 100},
            {Name: "content", Length: 255},
        },
        Comment: "Full-text search index",
    })
```

---

## API Reference

### Core Methods

#### CreateIndexWithOptions()
```go
CreateIndexWithOptions(name string, opts IndexOptions) (string, error)
```
Creates an index with advanced options. Supports all database-specific features.

**Parameters:**
- `name` - Index name
- `opts` - IndexOptions struct with advanced configuration

**IndexOptions Fields:**
- `Unique` - Emit UNIQUE keyword
- `IfNotExists` - Add IF NOT EXISTS where supported
- `Columns` - Array of IndexColumn with direction and prefix length
- `Using` - Index type (BTREE, HASH, GIN, GIST, BRIN, FULLTEXT, SPATIAL)
- `Include` - Covering index columns (PostgreSQL, MSSQL)
- `Where` - Partial index predicate (PostgreSQL, SQLite, MSSQL)
- `Storage` - Storage parameters (PostgreSQL WITH clause)
- `Comment` - Index comment (MySQL)

#### Convenience Methods

```go
CreateUniqueIndex(name string, columns ...string) (string, error)
CreateCompositeIndex(name string, columns []IndexColumn) (string, error)
CreatePartialIndex(name string, where string, columns ...string) (string, error)
CreateCoveringIndex(name string, include []string, columns ...string) (string, error)
DropIndexWithOptions(name string, opts DropIndexOptions) (string, error)
```

---

## Database-Specific Features

### MySQL Support
- ✅ **FULLTEXT indexes** - Full-text search capabilities
- ✅ **SPATIAL indexes** - Geospatial indexing
- ✅ **Prefix lengths** - Optimized for text/blob columns
- ✅ **Index comments** - Documentation and metadata
- ✅ **BTREE/HASH types** - Using clause inside column list

### PostgreSQL Support
- ✅ **Advanced types** - GIN, GIST, BRIN, HASH indexes
- ✅ **Partial indexes** - WHERE predicates for filtering
- ✅ **Covering indexes** - INCLUDE clause for additional columns
- ✅ **Storage parameters** - WITH clause for optimization
- ✅ **Schema qualification** - Schema.index_name syntax
- ✅ **IF NOT EXISTS** - Safe index creation

### SQLite Support
- ✅ **Partial indexes** - WHERE predicates supported
- ✅ **IF NOT EXISTS** - Safe index creation
- ✅ **Basic indexing** - Standard index functionality

### MSSQL Support
- ✅ **Covering indexes** - INCLUDE clause for additional columns
- ✅ **Partial indexes** - WHERE predicates for filtering
- ✅ **IF EXISTS** - Safe index creation (SQL Server 2016+)
- ✅ **Column direction** - ASC/DESC ordering

---

## Index Types

### Available Types
```go
const (
    INDEX_TYPE_BTREE   = "BTREE"
    INDEX_TYPE_HASH    = "HASH"
    INDEX_TYPE_GIN     = "INDEX_TYPE_GIN"
    INDEX_TYPE_GIST    = "INDEX_TYPE_GIST"
    INDEX_TYPE_BRIN    = "INDEX_TYPE_BRIN"
    INDEX_TYPE_FULLTEXT = "INDEX_TYPE_FULLTEXT"
    INDEX_TYPE_SPATIAL  = "INDEX_TYPE_SPATIAL"
)
```

### Type Usage by Database

| Type | MySQL | PostgreSQL | SQLite | MSSQL |
|------|-------|------------|--------|-------|
| BTREE | ✅ | ✅ | ✗ | ✅ |
| HASH | ✅ | ✅ | ✗ | ✗ |
| GIN | ✗ | ✅ | ✗ | ✗ |
| GIST | ✗ | ✅ | ✗ | ✗ |
| BRIN | ✗ | ✅ | ✗ | ✗ |
| FULLTEXT | ✅ | ✗ | ✗ | ✗ |
| SPATIAL | ✅ | ✗ | ✗ | ✗ |

---

## Usage Patterns

### Composite Indexes
```go
// Multi-column index with explicit ordering
sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
    Table("orders").
    CreateCompositeIndex("idx_orders_status_date", []sb.IndexColumn{
        {Name: "status", Direction: "ASC"},
        {Name: "created_at", Direction: "DESC"},
    })
```

### Covering Indexes
```go
// PostgreSQL covering index
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("documents").
    CreateCoveringIndex("idx_docs_search", []string{"title", "content"}, "search_vector")
```

### Partial Indexes
```go
// PostgreSQL partial index
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    CreatePartialIndex("idx_active_users", "deleted_at IS NULL", "email")
```

### Schema-Aware Operations
```go
// PostgreSQL schema-aware drop
sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
    Table("users").
    DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
        IfExists: true,
        Schema:   "public",
    })
```

---

## Error Handling

All index methods return `(string, error)` and use structured error types:

- `ErrEmptyIndexName` - Empty index name provided
- `ErrMissingTable` - No table specified
- `ErrEmptyColumns` - No columns specified

```go
sql, err := builder.CreateIndexWithOptions("", opts)
// Returns: "", ErrEmptyIndexName

sql, err := builder.CreateIndexWithOptions("idx_name", sb.IndexOptions{})
// Returns: "", ErrEmptyColumns
```

---

## Testing

The enhanced index support includes **32 comprehensive tests** covering:

- All database dialects (MySQL, PostgreSQL, SQLite, MSSQL)
- All index types and features
- Error handling scenarios
- Database-specific syntax validation
- Edge cases and complex scenarios

Run tests:
```bash
go test -v -run "TestBuilder.*Index.*"
```

---

## Migration from Basic Indexes

### Before (Basic API)
```go
// Old way - basic index only
sql := builder.CreateIndex("idx_name", "column")
```

### After (Enhanced API)
```go
// New way - enhanced with options
sql, err := builder.CreateIndexWithOptions("idx_name", sb.IndexOptions{
    Columns: []sb.IndexColumn{{Name: "column"}},
})
```

The enhanced API is **backward compatible** - existing code continues to work unchanged.

---

## Performance Considerations

### Index Selection
- **Choose appropriate index types** for your use case
- **Consider partial indexes** for large tables with selective queries
- **Use covering indexes** to avoid table lookups
- **Monitor index effectiveness** with database statistics

### Database-Specific Notes
- **MySQL**: FULLTEXT indexes require special syntax
- **PostgreSQL**: GIN indexes are ideal for array/json columns
- **SQLite**: Limited index types, focus on basic and partial indexes
- **MSSQL**: Consider filtered indexes for large tables

---

## Conclusion

The enhanced index support provides comprehensive database-agnostic index management while maintaining SB's philosophy of simplicity and consistency. The implementation supports advanced database features while ensuring full backward compatibility for existing code.

For detailed implementation information, see the memory bank documentation for the enhanced index support implementation.
