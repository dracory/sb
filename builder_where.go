package sb

import (
	"errors"
	"strings"
)

// Where represents a WHERE clause condition.
//
// Security Note: The Raw field allows arbitrary SQL and can be dangerous with user input.
// Prefer using Column/Operator/Value for safe parameter handling.
//
// Example:
//
//	// Safe - values are properly quoted
//	&sb.Where{Column: "email", Operator: "=", Value: userEmail}
//
//	// Dangerous - allows arbitrary SQL
//	&sb.Where{Raw: "email = '" + userInput + "'"}
type Where struct {
	Raw      string   // Raw SQL - use with caution, allows arbitrary SQL
	Column   string   // Column name - safe for validated column names
	Operator string   // SQL operator - safe for validated operators
	Type     string   // Logic type (AND/OR) - safe
	Value    string   // Column value - safely quoted
	Subquery *Builder // Subquery - safely handled with validation
	IsNot    bool     // NOT operator flag - safe
	Children []Where  // Nested conditions - safe
}

/**
 * Converts wheres to SQL
 * @param array $wheres
 * @return string
 */
func (b *Builder) whereToSql(wheres []Where) (string, error) {
	sql := []string{}
	for _, where := range wheres {
		// Validate subquery for nil case when operator indicates subquery usage
		if where.Subquery == nil && (where.Operator == "IN" || where.Operator == "NOT IN" ||
			where.Operator == "EXISTS" || where.Operator == "NOT EXISTS" ||
			(where.Operator != "" && where.Column != "" && where.Value == "")) {
			if err := b.validateSubqueryColumns(where); err != nil {
				return "", err
			}
		}

		if where.Raw != "" {
			sql = append(sql, where.Raw)
			continue
		}

		if where.Type == "" {
			where.Type = "AND"
		}

		if where.Subquery != nil {
			sqlSingle, err := b.whereToSqlSubquery(where)
			if err != nil {
				return "", err
			}
			if len(sql) > 0 {
				sql = append(sql, where.Type+" "+sqlSingle)
			} else {
				sql = append(sql, sqlSingle)
			}
			continue
		}

		if where.Column != "" {
			sqlSingle := b.whereToSqlSingle(where.Column, where.Operator, where.Value)

			if len(sql) > 0 {
				sql = append(sql, where.Type+" "+sqlSingle)
			} else {
				sql = append(sql, sqlSingle)
			}

		}
		// 			$_sql[] = $w['TYPE'] . " " . $sqlSingle;
		// 		}
		// 	}
		// 	$_sql = (count($_sql) > 0) ? " (" . implode(" ", $_sql) . ")" : "";

		// 	if ($i == 0) {
		// 		$sql[] = $_sql;
		// 	} else {
		// 		$sql[] = $where['TYPE'] . " " . $_sql;
		// 	}
		// }
	}

	if len(sql) > 0 {
		return " WHERE " + strings.Join(sql, " "), nil
	}

	return "", nil
}

func (b *Builder) whereToSqlSingle(column, operator, value string) string {
	if operator == "==" || operator == "===" {
		operator = "="
	}
	if operator == "!=" || operator == "!==" {
		operator = "<>"
	}
	columnQuoted := b.quoteColumn(column)

	// Use parameterized queries by default, unless interpolatedValues mode is enabled
	var valueQuoted string
	if b.interpolatedValues {
		valueQuoted = b.quoteValue(value)
	} else {
		// For parameterized queries, use placeholder
		if value == "NULL" {
			// NULL is a special case, not a parameter
			valueQuoted = "NULL"
		} else {
			valueQuoted = b.addParam(value)
		}
	}

	sql := ""
	if b.Dialect == DIALECT_MYSQL {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	if b.Dialect == DIALECT_POSTGRES {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	if b.Dialect == DIALECT_SQLITE {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	if b.Dialect == DIALECT_MSSQL {
		if value == "NULL" && operator == "=" {
			sql = columnQuoted + " IS NULL"
		} else if value == "NULL" && operator == "<>" {
			sql = columnQuoted + " IS NOT NULL"
		} else {
			sql = columnQuoted + " " + operator + " " + valueQuoted
		}
	}
	return sql
}

// validateSubqueryColumns validates that subquery returns appropriate number of columns for the operator
func (b *Builder) validateSubqueryColumns(where Where) error {
	if where.Subquery == nil {
		return errors.New("subquery cannot be nil")
	}

	columns := where.Subquery.sqlSelectColumns
	if len(columns) == 0 {
		// No columns explicitly set, check if this is a problem based on operator
		switch where.Operator {
		case "IN", "NOT IN":
			// For IN operations without explicit columns, we need to validate
			// Since we can't determine the actual columns, we'll assume single column for safety
			return nil // Allow for now, but this could be enhanced
		case "=", "!=", ">", "<", ">=", "<=":
			// For comparison operations without explicit columns, assume single column
			return nil
		case "EXISTS", "NOT EXISTS":
			// EXISTS doesn't care about column count
			return nil
		default:
			// Unknown operator, allow
			return nil
		}
	}

	switch where.Operator {
	case "IN", "NOT IN":
		if len(columns) > 1 {
			return errors.New("IN/NOT IN subquery must select exactly one column")
		}
	case "=", "!=", ">", "<", ">=", "<=":
		if len(columns) > 1 {
			return errors.New("comparison subquery must select exactly one column")
		}
	case "EXISTS", "NOT EXISTS":
		// Column count doesn't matter for EXISTS
		return nil
	default:
		// Unknown operator, allow for flexibility
		return nil
	}

	return nil
}

// whereToSqlSubquery converts a subquery WHERE condition to SQL
func (b *Builder) whereToSqlSubquery(where Where) (string, error) {
	// Validate subquery columns
	if err := b.validateSubqueryColumns(where); err != nil {
		return "", err
	}

	// Get the columns from the subquery builder
	columns := where.Subquery.sqlSelectColumns
	if len(columns) == 0 {
		columns = []string{"*"} // Default to all columns
	}

	// Generate subquery SQL without the trailing semicolon
	// Note: Select now returns (sql, params, error) but we only need the SQL for subqueries
	subquerySQL, _, err := where.Subquery.Select(columns)
	if err != nil {
		return "", err
	}
	// Remove the trailing semicolon from subquery
	subquerySQL = strings.TrimSuffix(subquerySQL, ";")

	// Handle different subquery operators
	switch where.Operator {
	case "EXISTS":
		if where.IsNot {
			return "NOT EXISTS (" + subquerySQL + ")", nil
		}
		return "EXISTS (" + subquerySQL + ")", nil
	case "IN":
		if where.Column != "" {
			columnQuoted := b.quoteColumn(where.Column)
			if where.IsNot {
				return columnQuoted + " NOT IN (" + subquerySQL + ")", nil
			}
			return columnQuoted + " IN (" + subquerySQL + ")", nil
		} else {
			// IN without column (used in WHERE clauses without explicit column)
			// Default to "id" for IN operations without explicit column
			defaultColumn := "id"
			if where.IsNot {
				return b.quoteColumn(defaultColumn) + " NOT IN (" + subquerySQL + ")", nil
			}
			return b.quoteColumn(defaultColumn) + " IN (" + subquerySQL + ")", nil
		}
	default:
		// For comparison operators (=, >, <, etc.)
		if where.Column != "" {
			columnQuoted := b.quoteColumn(where.Column)
			if where.IsNot {
				return columnQuoted + " NOT " + where.Operator + " (" + subquerySQL + ")", nil
			}
			return columnQuoted + " " + where.Operator + " (" + subquerySQL + ")", nil
		}
		// Fallback for edge cases
		return "(" + subquerySQL + ")", nil
	}
}
