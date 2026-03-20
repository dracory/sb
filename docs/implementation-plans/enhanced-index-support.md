# Enhanced Index Support Implementation Plan

**Priority:** MEDIUM  
**Target:** v0.19.0  
**Estimated Effort:** 1 week  
**Dependencies:** None (additive only, no breaking changes)

---

## Overview

Extend index management to support UNIQUE indexes, partial indexes, composite indexes with column ordering, and database-specific index types (GIN, FULLTEXT, etc.) while keeping full backward compatibility with the existing `CreateIndex` / `DropIndex` API.

---

## Problem Statement

The current `CreateIndex` method only generates a basic `CREATE INDEX name ON table (col)` statement. It has no support for:

- UNIQUE constraints
- Composite indexes with per-column ASC/DESC ordering
- Partial indexes (WHERE clause)
- PostgreSQL-specific types: GIN, GIST, BRIN, HASH, covering indexes (INCLUDE)
- MySQL-specific types: FULLTEXT, SPATIAL, prefix lengths on text columns
- MSSQL filtered indexes and included columns
- IF NOT EXISTS on index creation

---

## Solution Design

Add `IndexOptions` / `IndexColumn` types and a `CreateIndexWithOptions` method that all convenience helpers delegate to. Dialect-specific methods are kept off `BuilderInterface` to avoid forcing every dialect to implement features it doesn't support.

---

## Step 1: New Types and Constants

Add to `consts.go`:

```go
// Index type constants
const (
    INDEX_TYPE_BTREE   = "BTREE"
    INDEX_TYPE_HASH    = "HASH"
    INDEX_TYPE_GIN     = "GIN"
    INDEX_TYPE_GIST    = "GIST"
    INDEX_TYPE_BRIN    = "BRIN"
    INDEX_TYPE_FULLTEXT = "FULLTEXT"
    INDEX_TYPE_SPATIAL  = "SPATIAL"
)
```

Add new types (e.g. in `builder.go` alongside `Column`, `GroupBy`, etc.):

```go
// IndexColumn defines a single column in an index, with optional direction and prefix length.
type IndexColumn struct {
    Name      string // column name
    Direction string // "ASC" or "DESC" — defaults to ASC if empty
    Length    int    // prefix length for MySQL text/blob columns (ignored on other dialects)
}

// IndexOptions controls advanced index creation behaviour.
type IndexOptions struct {
    Unique      bool          // emit UNIQUE keyword
    IfNotExists bool          // emit IF NOT EXISTS where supported
    Columns     []IndexColumn // columns with optional direction / prefix length
    Using       string        // index type: BTREE, HASH, GIN, GIST, BRIN, FULLTEXT, SPATIAL
    Include     []string      // PostgreSQL / MSSQL covering-index columns (INCLUDE clause)
    Where       string        // partial index predicate (PostgreSQL, SQLite, MSSQL)
    Storage     string        // PostgreSQL WITH (...) storage parameters, e.g. "fillfactor=90"
    Comment     string        // MySQL COMMENT on index
}

// DropIndexOptions controls advanced index removal behaviour.
type DropIndexOptions struct {
    IfExists bool   // emit IF EXISTS where supported
    Schema   string // schema-qualify the index name (PostgreSQL only)
}
```

---

## Step 2: Core `CreateIndexWithOptions`

Add to `builder.go`:

```go
// CreateIndexWithOptions generates a CREATE INDEX statement with full dialect-aware options.
//
// Example (PostgreSQL partial unique index):
//
//   sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
//       Table("users").
//       CreateIndexWithOptions("idx_users_active_email", sb.IndexOptions{
//           Unique: true,
//           Columns: []sb.IndexColumn{{Name: "email"}},
//           Where:  "deleted_at IS NULL",
//       })
//   // CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_active_email"
//   //   ON "users" ("email") WHERE deleted_at IS NULL;
func (b *Builder) CreateIndexWithOptions(name string, opts IndexOptions) (string, error) {
    if err := b.validateAndReturnError(); err != nil {
        return "", err
    }
    if name == "" {
        return "", ErrEmptyIndexName
    }
    if b.sqlTableName == "" {
        return "", ErrMissingTable
    }
    if len(opts.Columns) == 0 {
        return "", ErrEmptyColumns
    }

    var sb strings.Builder

    sb.WriteString("CREATE ")
    if opts.Unique {
        sb.WriteString("UNIQUE ")
    }
    sb.WriteString("INDEX ")

    // IF NOT EXISTS — supported by PostgreSQL, SQLite; not by MySQL or MSSQL
    if opts.IfNotExists && (b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE) {
        sb.WriteString("IF NOT EXISTS ")
    }

    sb.WriteString(b.quoteTable(name))
    sb.WriteString(" ON ")
    sb.WriteString(b.quoteTable(b.sqlTableName))

    // USING clause — PostgreSQL only for standard index types; MySQL uses it for FULLTEXT/SPATIAL
    if opts.Using != "" {
        switch b.Dialect {
        case DIALECT_POSTGRES:
            sb.WriteString(" USING ")
            sb.WriteString(opts.Using)
        case DIALECT_MYSQL:
            // MySQL uses USING inside the column list for BTREE/HASH,
            // but FULLTEXT/SPATIAL are keywords before the column list.
            if opts.Using == INDEX_TYPE_FULLTEXT || opts.Using == INDEX_TYPE_SPATIAL {
                // handled below — rewrite the CREATE line
            }
        }
    }

    // MySQL FULLTEXT / SPATIAL require the keyword before the column list,
    // replacing "INDEX" entirely. Rebuild from scratch for these cases.
    if b.Dialect == DIALECT_MYSQL &&
        (opts.Using == INDEX_TYPE_FULLTEXT || opts.Using == INDEX_TYPE_SPATIAL) {
        sb.Reset()
        sb.WriteString("CREATE ")
        sb.WriteString(opts.Using) // FULLTEXT or SPATIAL
        sb.WriteString(" INDEX ")
        sb.WriteString(b.quoteTable(name))
        sb.WriteString(" ON ")
        sb.WriteString(b.quoteTable(b.sqlTableName))
    }

    // Column list
    sb.WriteString(" (")
    for i, col := range opts.Columns {
        if i > 0 {
            sb.WriteString(", ")
        }
        sb.WriteString(b.quoteColumn(col.Name))

        // MySQL prefix length must come immediately after the column name, before direction
        if b.Dialect == DIALECT_MYSQL && col.Length > 0 {
            sb.WriteString("(")
            sb.WriteString(strconv.Itoa(col.Length))
            sb.WriteString(")")
        }

        // Direction (not meaningful for FULLTEXT/SPATIAL but harmless to omit)
        if col.Direction != "" &&
            !(b.Dialect == DIALECT_MYSQL &&
                (opts.Using == INDEX_TYPE_FULLTEXT || opts.Using == INDEX_TYPE_SPATIAL)) {
            sb.WriteString(" ")
            sb.WriteString(strings.ToUpper(col.Direction))
        }
    }

    // MySQL BTREE/HASH USING goes inside the column list parentheses
    if b.Dialect == DIALECT_MYSQL && opts.Using != "" &&
        opts.Using != INDEX_TYPE_FULLTEXT && opts.Using != INDEX_TYPE_SPATIAL {
        sb.WriteString(" USING ")
        sb.WriteString(opts.Using)
    }

    sb.WriteString(")")

    // INCLUDE clause — PostgreSQL 11+ and MSSQL
    if len(opts.Include) > 0 &&
        (b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_MSSQL) {
        sb.WriteString(" INCLUDE (")
        for i, col := range opts.Include {
            if i > 0 {
                sb.WriteString(", ")
            }
            sb.WriteString(b.quoteColumn(col))
        }
        sb.WriteString(")")
    }

    // Partial index WHERE clause — PostgreSQL, SQLite, MSSQL
    if opts.Where != "" &&
        (b.Dialect == DIALECT_POSTGRES ||
            b.Dialect == DIALECT_SQLITE ||
            b.Dialect == DIALECT_MSSQL) {
        sb.WriteString(" WHERE ")
        sb.WriteString(opts.Where)
    }

    // PostgreSQL storage parameters
    if b.Dialect == DIALECT_POSTGRES && opts.Storage != "" {
        sb.WriteString(" WITH (")
        sb.WriteString(opts.Storage)
        sb.WriteString(")")
    }

    // MySQL index comment
    if b.Dialect == DIALECT_MYSQL && opts.Comment != "" {
        sb.WriteString(" COMMENT '")
        sb.WriteString(strings.ReplaceAll(opts.Comment, "'", "''"))
        sb.WriteString("'")
    }

    sb.WriteString(";")
    return sb.String(), nil
}
```

---

## Step 3: Convenience Methods

These all delegate to `CreateIndexWithOptions` and are added to `BuilderInterface`.

```go
// CreateUniqueIndex creates a UNIQUE index on one or more columns.
func (b *Builder) CreateUniqueIndex(name string, columns ...string) (string, error) {
    return b.CreateIndexWithOptions(name, IndexOptions{
        Unique:   true,
        Columns:  indexColumnsFromNames(columns),
    })
}

// CreateCompositeIndex creates an index on multiple columns with explicit ordering.
func (b *Builder) CreateCompositeIndex(name string, columns []IndexColumn) (string, error) {
    return b.CreateIndexWithOptions(name, IndexOptions{Columns: columns})
}

// CreatePartialIndex creates an index with a WHERE predicate (PostgreSQL, SQLite, MSSQL).
func (b *Builder) CreatePartialIndex(name string, where string, columns ...string) (string, error) {
    return b.CreateIndexWithOptions(name, IndexOptions{
        Columns: indexColumnsFromNames(columns),
        Where:   where,
    })
}

// CreateCoveringIndex creates a covering index using the INCLUDE clause (PostgreSQL, MSSQL).
func (b *Builder) CreateCoveringIndex(name string, include []string, columns ...string) (string, error) {
    return b.CreateIndexWithOptions(name, IndexOptions{
        Columns: indexColumnsFromNames(columns),
        Include: include,
    })
}

// indexColumnsFromNames is a package-level helper that converts plain column names
// to []IndexColumn with default ASC direction.
func indexColumnsFromNames(names []string) []IndexColumn {
    cols := make([]IndexColumn, len(names))
    for i, n := range names {
        cols[i] = IndexColumn{Name: n, Direction: "ASC"}
    }
    return cols
}
```

---

## Step 4: Enhanced `DropIndexWithOptions`

Replace the existing `DropIndexWithSchema` with a more general method. Keep `DropIndex`, `DropIndexIfExists`, and `DropIndexWithSchema` unchanged for backward compatibility.

```go
// DropIndexWithOptions generates a DROP INDEX statement with full dialect-aware options.
func (b *Builder) DropIndexWithOptions(name string, opts DropIndexOptions) (string, error) {
    if err := b.validateAndReturnError(); err != nil {
        return "", err
    }
    if name == "" {
        return "", ErrEmptyIndexName
    }

    var sb strings.Builder
    sb.WriteString("DROP INDEX ")

    switch b.Dialect {
    case DIALECT_POSTGRES:
        if opts.IfExists {
            sb.WriteString("IF EXISTS ")
        }
        if opts.Schema != "" {
            sb.WriteString(b.quoteTable(opts.Schema))
            sb.WriteString(".")
        }
        sb.WriteString(b.quoteTable(name))

    case DIALECT_SQLITE:
        if opts.IfExists {
            sb.WriteString("IF EXISTS ")
        }
        sb.WriteString(b.quoteTable(name))

    case DIALECT_MSSQL:
        if opts.IfExists {
            sb.WriteString("IF EXISTS ")
        }
        sb.WriteString(b.quoteTable(name))
        if b.sqlTableName != "" {
            sb.WriteString(" ON ")
            sb.WriteString(b.quoteTable(b.sqlTableName))
        }

    case DIALECT_MYSQL:
        // MySQL has no IF EXISTS for DROP INDEX
        sb.WriteString(b.quoteTable(name))
        if b.sqlTableName != "" {
            sb.WriteString(" ON ")
            sb.WriteString(b.quoteTable(b.sqlTableName))
        }

    default:
        return "", ErrInvalidDialect
    }

    sb.WriteString(";")
    return sb.String(), nil
}
```

---

## Step 5: Interface Updates

Add to `BuilderInterface` in `interfaces.go`:

```go
// CreateIndexWithOptions creates an index with advanced options.
CreateIndexWithOptions(name string, opts IndexOptions) (string, error)

// CreateUniqueIndex creates a UNIQUE index on one or more columns.
CreateUniqueIndex(name string, columns ...string) (string, error)

// CreateCompositeIndex creates an index on multiple columns with explicit ordering.
CreateCompositeIndex(name string, columns []IndexColumn) (string, error)

// CreatePartialIndex creates an index with a WHERE predicate (PostgreSQL, SQLite, MSSQL).
CreatePartialIndex(name string, where string, columns ...string) (string, error)

// CreateCoveringIndex creates a covering index using the INCLUDE clause (PostgreSQL, MSSQL).
CreateCoveringIndex(name string, include []string, columns ...string) (string, error)

// DropIndexWithOptions drops an index with advanced options.
DropIndexWithOptions(name string, opts DropIndexOptions) (string, error)
```

Note: dialect-specific methods (`CreateMySQLFullTextIndex`, `CreatePostgreSQLGINIndex`, etc.) are NOT added to the interface. Callers that need them work directly with `*Builder` or use `CreateIndexWithOptions` with the appropriate `Using` value.

---

## Dialect Feature Matrix

| Feature              | MySQL | PostgreSQL | SQLite | MSSQL |
|----------------------|-------|------------|--------|-------|
| UNIQUE index         | ✓     | ✓          | ✓      | ✓     |
| Composite index      | ✓     | ✓          | ✓      | ✓     |
| Column direction     | ✓     | ✓          | ✓      | ✓     |
| IF NOT EXISTS        | ✗     | ✓          | ✓      | ✗     |
| Partial index (WHERE)| ✗     | ✓          | ✓      | ✓     |
| INCLUDE columns      | ✗     | ✓ (v11+)   | ✗      | ✓     |
| USING clause         | ✓*    | ✓          | ✗      | ✗     |
| FULLTEXT index       | ✓     | ✗          | ✗      | ✗     |
| SPATIAL index        | ✓     | ✗          | ✗      | ✗     |
| GIN / GIST / BRIN    | ✗     | ✓          | ✗      | ✗     |
| Prefix length        | ✓     | ✗          | ✗      | ✗     |
| Storage params (WITH)| ✗     | ✓          | ✗      | ✗     |
| Index COMMENT        | ✓     | ✗          | ✗      | ✗     |

*MySQL USING goes inside the column list parentheses, not before it.

---

## Known Dialect Quirks (Implementation Notes)

- **MySQL IF NOT EXISTS** — not supported for `CREATE INDEX`; silently omitted.
- **MySQL FULLTEXT / SPATIAL** — the keyword replaces `INDEX` in the statement (`CREATE FULLTEXT INDEX ...`), not a `USING` modifier.
- **MySQL prefix length** — must appear immediately after the column name and before any direction: `` `col`(255) ASC ``.
- **MySQL USING BTREE/HASH** — goes inside the column list: `(col ASC USING BTREE)`.
- **PostgreSQL INCLUDE** — requires PostgreSQL 11+; silently omitted on older versions (caller's responsibility).
- **MSSQL IF EXISTS** — requires SQL Server 2016+; document this in release notes.
- **MSSQL DROP INDEX** — always requires `ON table_name`; `DropIndexWithOptions` enforces this.
- **SQLite** — no INCLUDE, no USING, no storage params; unsupported options are silently ignored.

---

## Testing Strategy

### Unit Tests (`builder_test.go` or a new `builder_index_test.go`)

Cover every dialect × feature combination from the matrix above:

- Basic index (regression — existing `CreateIndex` behaviour unchanged)
- UNIQUE index on single and multiple columns
- Composite index with mixed ASC/DESC columns
- Partial index (PostgreSQL, SQLite, MSSQL)
- Covering index with INCLUDE (PostgreSQL, MSSQL)
- FULLTEXT index (MySQL only)
- SPATIAL index (MySQL only)
- GIN / GIST index (PostgreSQL only)
- Prefix length on text column (MySQL only)
- IF NOT EXISTS respected / silently omitted per dialect
- `DropIndexWithOptions` with IfExists and Schema

### Error Cases

- Empty index name → `ErrEmptyIndexName`
- Empty table name → `ErrMissingTable`
- Empty columns slice → `ErrEmptyColumns`

### Integration Tests (`integration_test.go`)

- Create and drop each index type against a live database container
- Verify UNIQUE constraint is enforced at the database level
- Verify partial index filters rows correctly

---

## Success Criteria

- [ ] All existing `CreateIndex` / `DropIndex` tests still pass (no regression)
- [ ] New unit tests cover every cell in the dialect feature matrix
- [ ] `CreateIndexWithOptions` and convenience methods added to `BuilderInterface`
- [ ] `DropIndexWithOptions` added to `BuilderInterface`
- [ ] Dialect quirks handled correctly (MySQL FULLTEXT keyword, prefix length position, etc.)
- [ ] No dialect-specific methods on `BuilderInterface`

---

## Timeline

- Day 1–2: Types, constants, `CreateIndexWithOptions` core logic
- Day 3: Convenience methods + `DropIndexWithOptions`
- Day 4: Unit tests for all dialects
- Day 5: Integration tests + documentation update

---

## Backward Compatibility

All existing methods (`CreateIndex`, `DropIndex`, `DropIndexIfExists`, `DropIndexWithSchema`) remain unchanged. New methods are purely additive.
