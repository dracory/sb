package sb

import (
	"strings"
)

type Where struct {
	Raw      string
	Column   string
	Operator string
	Type     string
	Value    string
	Subquery *Builder
	IsNot    bool
	Children []Where
}

/**
 * Converts wheres to SQL
 * @param array $wheres
 * @return string
 */
func (b *Builder) whereToSql(wheres []Where) string {
	sql := []string{}
	for _, where := range wheres {
		if where.Raw != "" {
			sql = append(sql, where.Raw)
			continue
		}

		if where.Type == "" {
			where.Type = "AND"
		}

		if where.Subquery != nil {
			sqlSingle := b.whereToSqlSubquery(where)
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
		// } else {
		// 	$_sql = array();
		// 	$all = $where['WHERE'];
		// 	for ($k = 0; k < count($all); k++) {
		// 		$w = $all[$k];
		// 		$sqlSingle = $this->whereToSqlSingle($w['COLUMN'], $w['OPERATOR'], $w['VALUE']);
		// 		if ($k == 0) {
		// 			$_sql[] = $sqlSingle;
		// 		} else {
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
		return " WHERE " + strings.Join(sql, " ")
	}

	return ""
}

func (b *Builder) whereToSqlSingle(column, operator, value string) string {
	if operator == "==" || operator == "===" {
		operator = "="
	}
	if operator == "!=" || operator == "!==" {
		operator = "<>"
	}
	columnQuoted := b.quoteColumn(column)
	valueQuoted := b.quoteValue(value)

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

// whereToSqlSubquery converts a subquery WHERE condition to SQL
func (b *Builder) whereToSqlSubquery(where Where) string {
	// Get the columns from the subquery builder
	columns := where.Subquery.sqlSelectColumns
	if len(columns) == 0 {
		columns = []string{"*"} // Default to all columns
	}

	// Generate subquery SQL without the trailing semicolon
	subquerySQL := where.Subquery.Select(columns)
	// Remove the trailing semicolon from subquery
	subquerySQL = strings.TrimSuffix(subquerySQL, ";")

	// Handle different subquery operators
	switch where.Operator {
	case "EXISTS":
		if where.IsNot {
			return "NOT EXISTS (" + subquerySQL + ")"
		}
		return "EXISTS (" + subquerySQL + ")"
	case "IN":
		if where.Column != "" {
			columnQuoted := b.quoteColumn(where.Column)
			if where.IsNot {
				return columnQuoted + " NOT IN (" + subquerySQL + ")"
			}
			return columnQuoted + " IN (" + subquerySQL + ")"
		} else {
			// IN without column (used in WHERE clauses without explicit column)
			// Default to "id" for IN operations without explicit column
			defaultColumn := "id"
			if where.IsNot {
				return b.quoteColumn(defaultColumn) + " NOT IN (" + subquerySQL + ")"
			}
			return b.quoteColumn(defaultColumn) + " IN (" + subquerySQL + ")"
		}
	default:
		// For comparison operators (=, >, <, etc.)
		if where.Column != "" {
			columnQuoted := b.quoteColumn(where.Column)
			if where.IsNot {
				return columnQuoted + " NOT " + where.Operator + " (" + subquerySQL + ")"
			}
			return columnQuoted + " " + where.Operator + " (" + subquerySQL + ")"
		}
		// Fallback for edge cases
		return "(" + subquerySQL + ")"
	}
}
