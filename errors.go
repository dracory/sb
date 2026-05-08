package sb

import (
	"fmt"
)

// BuilderError represents a structured error from the SQL builder
type BuilderError struct {
	Type    string
	Message string
	Cause   error
}

// Error returns the error message
func (e *BuilderError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying cause
func (e *BuilderError) Unwrap() error {
	return e.Cause
}

// Common error types
var (
	// Validation errors
	ErrEmptyTableName             = &BuilderError{Type: "ValidationError", Message: "table name cannot be empty"}
	ErrEmptyColumnName            = &BuilderError{Type: "ValidationError", Message: "column name cannot be empty"}
	ErrEmptyIndexName             = &BuilderError{Type: "ValidationError", Message: "index name cannot be empty"}
	ErrEmptyColumns               = &BuilderError{Type: "ValidationError", Message: "columns cannot be empty"}
	ErrEmptyOnCondition           = &BuilderError{Type: "ValidationError", Message: "ON condition cannot be empty"}
	ErrInvalidJoinType            = &BuilderError{Type: "ValidationError", Message: "invalid join type"}
	ErrOffsetWithoutLimit         = &BuilderError{Type: "ValidationError", Message: "SQLite requires LIMIT when using OFFSET"}
	ErrMSSQLOffsetRequiresOrderBy = &BuilderError{Type: "ValidationError", Message: "MSSQL requires ORDER BY when using OFFSET"}

	// Configuration errors
	ErrInvalidDialect = &BuilderError{Type: "ConfigurationError", Message: "invalid database dialect"}
	ErrMissingTable   = &BuilderError{Type: "ValidationError", Message: "no table specified"}
	ErrNilQueryable   = &BuilderError{Type: "ArgumentError", Message: "queryable cannot be nil"}

	// Subquery errors
	ErrNilSubquery     = &BuilderError{Type: "ArgumentError", Message: "subquery cannot be nil"}
	ErrSubqueryColumns = &BuilderError{Type: "SubqueryError", Message: "subquery columns validation failed"}
)

// NewValidationError creates a new validation error
func NewValidationError(message string) *BuilderError {
	return &BuilderError{
		Type:    "ValidationError",
		Message: message,
	}
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(message string) *BuilderError {
	return &BuilderError{
		Type:    "ConfigurationError",
		Message: message,
	}
}

// NewSubqueryError creates a new subquery error with optional cause
func NewSubqueryError(message string, cause error) *BuilderError {
	return &BuilderError{
		Type:    "SubqueryError",
		Message: message,
		Cause:   cause,
	}
}
