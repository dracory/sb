// Package schema provides DB-execution functions for managing table and column
// schema operations. It complements the github.com/dracory/sb SQL builder by
// executing DDL statements (CREATE, DROP, ALTER) against a live database.
//
// Usage:
//
//	import (
//	    "github.com/dracory/sb"
//	    "github.com/dracory/sb/schema"
//	)
//
//	err := schema.TableColumnAdd(ctx, "users", sb.Column{Name: "email", Type: sb.COLUMN_TYPE_STRING, Length: 255})
package schema
