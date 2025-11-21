package sb

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

type OrderBy struct {
	Column    string
	Direction string
}

type Column struct {
	Name          string
	Type          string
	Length        int
	Decimals      int
	AutoIncrement bool
	PrimaryKey    bool
	Nullable      bool
	Unique        bool
	Default       string
}

type GroupBy struct {
	Column string
}

type Builder struct {
	Dialect            string
	sql                map[string]any
	sqlColumns         []Column
	sqlGroupBy         []GroupBy
	sqlLimit           int64
	sqlOffset          int64
	sqlOrderBy         []OrderBy
	sqlTableName       string
	sqlViewName        string
	sqlViewColumns     []string
	sqlViewSQL         string
	sqlWhere           []Where
	columnSQLGenerator ColumnSQLGenerator
}

var _ BuilderInterface = (*Builder)(nil)

func NewBuilder(dialect string) *Builder {
	var columnSQLGenerator ColumnSQLGenerator
	switch dialect {
	case DIALECT_MYSQL:
		columnSQLGenerator = MySQLColumnSQLGenerator{}
	case DIALECT_POSTGRES:
		columnSQLGenerator = PostgreSQLColumnSQLGenerator{}
	case DIALECT_SQLITE:
		columnSQLGenerator = SQLiteColumnSQLGenerator{}
	case DIALECT_MSSQL:
		columnSQLGenerator = MSSQLColumnSQLGenerator{}
	default:
		panic("unsupported dialect: " + dialect)
	}

	return &Builder{
		Dialect:            dialect,
		sql:                map[string]any{},
		sqlColumns:         []Column{},
		sqlGroupBy:         []GroupBy{},
		sqlLimit:           0,
		sqlOffset:          0,
		sqlOrderBy:         []OrderBy{},
		sqlTableName:       "",
		sqlViewName:        "",
		sqlViewColumns:     []string{},
		sqlViewSQL:         "",
		sqlWhere:           []Where{},
		columnSQLGenerator: columnSQLGenerator,
	}
}

func (b *Builder) Table(tableName string) BuilderInterface {
	b.sqlTableName = tableName
	return b
}

func (b *Builder) View(viewName string) BuilderInterface {
	b.sqlViewName = viewName
	return b
}

func (b *Builder) ViewSQL(sql string) BuilderInterface {
	b.sqlViewSQL = sql
	return b
}

func (b *Builder) ViewColumns(columns []string) BuilderInterface {
	b.sqlViewColumns = columns
	return b
}

func (b *Builder) Column(column Column) BuilderInterface {
	if column.Name == "" {
		panic("column name is required")
	}

	if column.Type == "" {
		panic("column type is required")
	}

	b.sqlColumns = append(b.sqlColumns, column)

	return b
}

/**
 * The create method creates new database or table.
 * If the database or table can not be created it will return false.
 * False will be returned if the database or table already exist.
 * <code>
 * // Creating a new database
 * $database->create();
 *
 * // Creating a new table
 * $database->table("STATES")
 *     ->column("STATE_NAME","STRING")
 *     ->create();
 * </code>
 * @return boolean true, on success, false, otherwise
 * @access public
 */
func (b *Builder) Create() string {
	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	sql := ""

	if isTable {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = `CREATE TABLE ` + b.quoteTable(b.sqlTableName) + `(` + b.columnsToSQL(b.sqlColumns) + `);`
		}
		if b.Dialect == DIALECT_MSSQL {
			sql = `CREATE TABLE [` + b.sqlTableName + `] (` + b.columnsToSQL(b.sqlColumns) + `);`
		}
	}

	if isView {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			viewColumnsToSQL := strings.Join(lo.Map(b.sqlViewColumns, func(columnName string, _ int) string {
				return b.quoteColumn(columnName)
			}), ", ")
			viewColumns := lo.If(len(b.sqlViewColumns) > 0, ` (`+viewColumnsToSQL+`)`).Else(``)

			sql = `CREATE VIEW ` + b.quoteTable(b.sqlViewName) + viewColumns + " AS " + b.sqlViewSQL
		}
	}

	return sql
}

func (b *Builder) CreateIfNotExists() string {
	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	sql := ""

	if isTable {
		if b.Dialect == DIALECT_MYSQL {
			sql = "CREATE TABLE IF NOT EXISTS " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
		}
		if b.Dialect == DIALECT_POSTGRES {
			sql = `CREATE TABLE IF NOT EXISTS ` + b.quoteTable(b.sqlTableName) + `(` + b.columnsToSQL(b.sqlColumns) + `);`
		}
		if b.Dialect == DIALECT_SQLITE {
			sql = "CREATE TABLE IF NOT EXISTS " + b.quoteTable(b.sqlTableName) + "(" + b.columnsToSQL(b.sqlColumns) + ");"
		}
	}

	if isView {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			viewColumnsToSQL := strings.Join(lo.Map(b.sqlViewColumns, func(columnName string, _ int) string {
				return b.quoteColumn(columnName)
			}), ", ")
			viewColumns := lo.If(len(b.sqlViewColumns) > 0, ` (`+viewColumnsToSQL+`)`).Else(``)

			sqlStart := "CREATE VIEW IF NOT EXISTS"
			if b.Dialect == DIALECT_MYSQL {
				sqlStart = "CREATE OR REPLACE VIEW"
			}

			sql = sqlStart + ` ` + b.quoteTable(b.sqlViewName) + viewColumns + " AS " + b.sqlViewSQL
		}
	}

	return sql
}

func (b *Builder) CreateIndex(indexName string, columnName ...string) string {
	if b.sqlTableName == "" {
		panic("In method CreateIndex() no table specified to create index on!")
	}

	columns := lo.Map(columnName, func(columnName string, i int) string {
		return b.quoteColumn(columnName)
	})

	sql := `CREATE INDEX ` + b.quoteTable(indexName) + ` ON ` + b.quoteTable(b.sqlTableName) + ` (` + strings.Join(columns, `,`) + `);`

	return sql
}

/**
 * The delete method deletes a row in a table. For deleting a database
 * or table use the drop method.
 * <code>
 * // Deleting a row
 * sql := builder.Table("STATES").Where("STATE_NAME","=","Alabama").Delete();
 * </code>
 * @return string
 * @access public
 */
// Drop deletes a table
func (b *Builder) Delete() string {
	if b.sqlTableName == "" {
		panic("In method Delete() no table specified to delete from!")
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		where = b.whereToSql(b.sqlWhere)
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	sql := ""
	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
		sql = "DELETE FROM " + b.quoteTable(b.sqlTableName) + where + orderBy + limit + offset + ";"
	}
	return sql
}

// Drop deletes a table or a view
func (b *Builder) Drop() string {
	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	sql := ""

	if isTable {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP TABLE " + b.quoteTable(b.sqlTableName) + ";"
		}
	}

	if isView {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP VIEW " + b.quoteTable(b.sqlViewName) + ";"
		}
	}

	return sql
}

func (b *Builder) DropIfExists() string {
	isView := b.sqlViewName != ""
	isTable := b.sqlTableName != ""

	sql := ""

	if isTable {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP TABLE IF EXISTS " + b.quoteTable(b.sqlTableName) + ";"
		}
	}

	if isView {
		if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
			sql = "DROP VIEW IF EXISTS " + b.quoteTable(b.sqlViewName) + ";"
		}
	}

	return sql
}

func (b *Builder) Limit(limit int64) BuilderInterface {
	b.sqlLimit = limit
	return b
}

func (b *Builder) Offset(offset int64) BuilderInterface {
	b.sqlOffset = offset
	return b
}

func (b *Builder) GroupBy(groupBy GroupBy) BuilderInterface {
	b.sqlGroupBy = append(b.sqlGroupBy, groupBy)
	return b
}

func (b *Builder) OrderBy(columnName, direction string) BuilderInterface {
	if strings.EqualFold(direction, "desc") || strings.EqualFold(direction, "descending") {
		direction = "DESC"
	} else {
		direction = "ASC"
	}

	b.sqlOrderBy = append(b.sqlOrderBy, OrderBy{
		Column:    columnName,
		Direction: direction,
	})

	return b
}

// Rename renames a table or a view
func (b *Builder) TableRename(oldTableName, newTableName string) (sql string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sql = "EXEC sp_rename " + b.quoteTable(oldTableName) + ", " + b.quoteTable(newTableName) + ", 'OBJECT';"
		return sql, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sql = "ALTER TABLE " + b.quoteTable(oldTableName) + " RENAME TO " + b.quoteTable(newTableName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sql = "ALTER TABLE " + b.quoteTable(oldTableName) + " RENAME " + b.quoteTable(newTableName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sql = "ALTER TABLE " + b.quoteTable(oldTableName) + " RENAME TO " + b.quoteTable(newTableName) + ";"
		return sql, nil
	}

	return "", errors.New("renaming a table is not supported for driver " + b.Dialect + "")
}

// TableColumnAdd adds a column to the table
func (b *Builder) TableColumnAdd(tableName string, column Column) (sql string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " ADD " + b.columnsToSQL([]Column{column}) + ";"
		return sql, nil
	}

	return "", errors.New("adding a column is not supported for driver " + b.Dialect + "")
}

// TableColumnChange changes a column in the table
func (b *Builder) TableColumnChange(tableName string, column Column) (sqlString string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " ALTER COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " ALTER COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " MODIFY COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " ALTER COLUMN " + b.columnsToSQL([]Column{column}) + ";"
		return sqlString, nil
	}

	return "", errors.New("modifying a column is not supported for driver " + b.Dialect + "")
}

// TableColumnDrop drops a column from the table
func (b *Builder) TableColumnDrop(tableName, columnName string) (sqlString string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sqlString = "ALTER TABLE " + b.quoteTable(tableName) + " DROP COLUMN " + b.quoteColumn(columnName) + ";"
		return sqlString, nil
	}

	return "", errors.New("dropping a column is not supported for driver " + b.Dialect + "")
}

// TableColumnExists checks if a column exists in a table for various database types
//
//	Example:
//	b := NewBuilder(DIALECT_MYSQL)
//	sqlString, sqlParams, err := b.TableColumnExists("test_table", "test_column")
//
// Params:
// - tableName: The name of the table to check.
// - columnName: The name of the column to check.
//
// Returns:
// - sql: The SQL statement to check for the existence of the column.
// - params: An array of parameters to be bound to the statement.
// - err: An error object, if any.
func (b *Builder) TableColumnExists(tableName, columnName string) (sql string, params []interface{}, err error) {
	switch b.Dialect {
	case DIALECT_MYSQL:
		return "SELECT 1 FROM information_schema.COLUMNS WHERE TABLE_NAME = ? AND COLUMN_NAME = ?", []interface{}{tableName, columnName}, nil
	case DIALECT_POSTGRES:
		return "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name = $1 AND column_name = $2)", []interface{}{tableName, columnName}, nil
	case DIALECT_SQLITE:
		return "SELECT 1 FROM pragma_table_info(?) WHERE name = ?", []interface{}{tableName, columnName}, nil
	default:
		return "", nil, fmt.Errorf("database type '%s' not supported", b.Dialect)
	}
}

func (b *Builder) TableColumnRename(tableName, oldColumnName, newColumnName string) (sql string, err error) {
	if b.Dialect == DIALECT_MSSQL {
		sql = "EXEC sp_rename " + b.quoteTable(tableName) + "." + b.quoteTable(oldColumnName) + ", " + b.quoteTable(newColumnName) + ", 'COLUMN';"
		return sql, nil
	}

	if b.Dialect == DIALECT_SQLITE {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " RENAME COLUMN " + b.quoteTable(oldColumnName) + " TO " + b.quoteTable(newColumnName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_MYSQL {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " RENAME COLUMN " + b.quoteTable(oldColumnName) + " TO " + b.quoteTable(newColumnName) + ";"
		return sql, nil
	}

	if b.Dialect == DIALECT_POSTGRES {
		sql = "ALTER TABLE " + b.quoteTable(tableName) + " RENAME COLUMN " + b.quoteTable(oldColumnName) + " TO " + b.quoteTable(newColumnName) + ";"
		return sql, nil
	}

	return "", errors.New("renaming a column is not supported for driver " + b.Dialect + "")
}

/** The <b>select</b> method selects rows from a table, based on criteria.
 * <code>
 * // Selects all the rows from the table
 * $db->table("USERS")->select();
 *
 * // Selects the rows where the column NAME is different from Peter, in descending order
 * $db->table("USERS")
 *     ->where("NAME","!=","Peter")
 *     ->orderby("NAME","desc")
 *     ->select();
 * </code>
 * @return mixed rows as associative array, false on error
 * @access public
 */
func (b *Builder) Select(columns []string) string {
	if b.sqlTableName == "" {
		panic("In method Select() no table specified to select from!")
	}

	join := "" // TODO add support for joins

	groupBy := ""
	if len(b.sqlGroupBy) > 0 {
		groupBy = b.groupByToSql(b.sqlGroupBy)
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		where = b.whereToSql(b.sqlWhere)
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	columnsStr := "*"

	if len(columns) > 0 {
		for index, column := range columns {
			if strings.Contains(column, "(") {
				columns[index] = column // Do not quote function calls
			} else {
				columns[index] = b.quoteColumn(column)
			}
		}
		columnsStr = strings.Join(columns, ", ")
	}

	sql := ""

	if b.Dialect == DIALECT_MYSQL || b.Dialect == DIALECT_POSTGRES || b.Dialect == DIALECT_SQLITE {
		sql = "SELECT " + columnsStr + " FROM " + b.quoteTable(b.sqlTableName) + join + where + groupBy + orderBy + limit + offset + ";"
	}

	return sql
}

/**
 * The <b>update</b> method updates the values of a row in a table.
 * <code>
 * $updated_user = array("USER_MANE"=>"Mike");
 * $database->table("USERS")->where("USER_NAME","==","Peter")->update($updated_user);
 * </code>
 * @param Array an associative array, where keys are the column names of the table
 * @return int 0 or 1, on success, false, otherwise
 * @access public
 */
func (b *Builder) Insert(columnValuesMap map[string]string) string {
	if b.sqlTableName == "" {
		panic("In method Insert() no table specified to insert in!")
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	columnNames := []string{}
	columnValues := []string{}

	// Order keys
	keys := make([]string, 0, len(columnValuesMap))
	for k := range columnValuesMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, columnName := range keys {
		columnValue := columnValuesMap[columnName]
		columnNames = append(columnNames, b.quoteColumn(columnName))
		columnValues = append(columnValues, b.quoteValue(columnValue))
	}

	return "INSERT INTO " + b.quoteTable(b.sqlTableName) + " (" + strings.Join(columnNames, ", ") + ") VALUES (" + strings.Join(columnValues, ", ") + ")" + limit + offset + ";"
}

func (b *Builder) Truncate() string {
	// TODO: implement
	return ""
}

/**
 * The <b>update</b> method updates the values of a row in a table.
 * <code>
 * $updated_user = array("USER_MANE"=>"Mike");
 * $database->table("USERS")->where("USER_NAME","==","Peter")->update($updated_user);
 * </code>
 * @param Array an associative array, where keys are the column names of the table
 * @return int 0 or 1, on success, false, otherwise
 * @access public
 */
func (b *Builder) Update(columnValues map[string]string) string {
	if b.sqlTableName == "" {
		panic("In method Update() no table specified to update!")
	}

	join := "" // TODO add support for joins

	groupBy := ""
	if len(b.sqlGroupBy) > 0 {
		groupBy = b.groupByToSql(b.sqlGroupBy)
	}

	where := ""
	if len(b.sqlWhere) > 0 {
		where = b.whereToSql(b.sqlWhere)
	}

	orderBy := ""
	if len(b.sqlOrderBy) > 0 {
		orderBy = b.orderByToSql(b.sqlOrderBy)
	}

	limit := ""
	if b.sqlLimit > 0 {
		limit = " LIMIT " + strconv.FormatInt(b.sqlLimit, 10)
	}

	offset := ""
	if b.sqlOffset > 0 {
		offset = " OFFSET " + strconv.FormatInt(b.sqlOffset, 10)
	}

	// Order keys
	keys := make([]string, 0, len(columnValues))
	for k := range columnValues {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	updateSql := []string{}
	for _, columnName := range keys {
		columnValue := columnValues[columnName]
		updateSql = append(updateSql, b.quoteColumn(columnName)+"="+b.quoteValue(columnValue))
	}

	return "UPDATE " + b.quoteTable(b.sqlTableName) + " SET " + strings.Join(updateSql, ", ") + join + where + groupBy + orderBy + limit + offset + ";"
}

func (b *Builder) Where(where *Where) BuilderInterface {
	if where == nil {
		return b
	}

	b.sqlWhere = append(b.sqlWhere, *where)
	return b
}

// columnsToSQL converts the columns statements to SQL.
func (b *Builder) columnsToSQL(columns []Column) string {
	columnSQLs := []string{}

	for i := 0; i < len(columns); i++ {
		column := columns[i]
		columnSQLs = append(columnSQLs, b.columnSQLGenerator.GenerateSQL(column))
	}

	return strings.Join(columnSQLs, ", ")
}

func (b *Builder) groupByToSql(groupBys []GroupBy) string {
	sql := []string{}
	for _, groupBy := range groupBys {
		sql = append(sql, b.quoteColumn(groupBy.Column))
	}

	if len(sql) > 0 {
		return " GROUP BY " + strings.Join(sql, ",")
	}

	return ""
}

// /**
//      * Joins tables to SQL.
//      * @return String the join SQL string
//      * @access private
//      */
// 	 private function join_to_sql($join, $table_name)
// 	 {
// 		 $sql = '';
// 		 // MySQL
// 		 if ($this->database_type == 'mysql') {
// 			 foreach ($join as $what) {
// 				 $type = $what[3] ?? '';
// 				 $alias = $what[4] ?? '';
// 				 $sql .= ' ' . $type . ' JOIN `' . $what[0] . '`';
// 				 if ($alias != "") {
// 					 $sql .= ' AS ' . $alias . '';
// 					 $what[0] = $alias;
// 				 }
// 				 if ($what[1] == $what[2]) {
// 					 $sql .= ' USING (`' . $what[1] . '`)';
// 				 } else {
// 					 $sql .= ' ON ' . $table_name . '.' . $what[1] . '=' . $what[0] . '.' . $what[2];
// 				 }
// 			 }
// 		 }
// 		 // SQLite
// 		 if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
// 			 foreach ($join as $what) {
// 				 $type = $what[3] ?? '';
// 				 $alias = $what[4] ?? '';
// 				 $sql .= " $type JOIN '" . $what[0] . "'";
// 				 if ($alias != "") {
// 					 $sql .= " AS '$alias'";
// 					 $what[0] = $alias;
// 				 }
// 				 $sql .= ' ON ' . $table_name . '.' . $what[1] . '=' . $what[0] . '.' . $what[2];
// 			 }
// 		 }

// 		 return $sql;
// 	 }

func (b *Builder) orderByToSql(orderBys []OrderBy) string {
	sql := []string{}

	if b.Dialect == DIALECT_MYSQL {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if b.Dialect == DIALECT_POSTGRES {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if b.Dialect == DIALECT_SQLITE {
		for _, orderBy := range orderBys {
			sql = append(sql, b.quoteColumn(orderBy.Column)+" "+orderBy.Direction)
		}
	}

	if len(sql) > 0 {
		return ` ORDER BY ` + strings.Join(sql, `,`)
	}

	return ""
}

func (b *Builder) quoteColumn(columnName string) string {
	columnSplit := strings.Split(columnName, ".")
	columnQuoted := []string{}

	for _, columnPart := range columnSplit {
		if columnPart == "*" {
			columnQuoted = append(columnQuoted, columnPart)
			continue
		}

		if strings.Contains(columnPart, "(") {
			columnQuoted = append(columnQuoted, columnPart)
		}

		columnQuoted = append(columnQuoted, b.quote(columnPart, "column"))
	}

	return strings.Join(columnQuoted, ".")
}

func (b *Builder) quoteTable(tableName string) string {
	tableSplit := strings.Split(tableName, ".")
	tableQuoted := []string{}

	for _, tablePart := range tableSplit {
		tableQuoted = append(tableQuoted, b.quote(tablePart, "table"))
	}

	return strings.Join(tableQuoted, ".")
}

/**
 * The <b>tables</b> method returns the names of all the tables, that
 * exist in the database.
 * <code>
 * foreach($database->tables() as $table){
 *     echo $table;
 * }
 * </code>
 * @param String the name of the table
 * @return array the names of the tables
 * @access public
 */
//  func (b *Builder) Tables(value string)
//  {
// 	 $tables = array();

// 	 if ($this->database_type == 'mysql') {
// 		 //$sql = "SHOW TABLES";
// 		 $sql = "SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_TYPE='BASE TABLE' AND TABLE_SCHEMA='" . $this->database_name . "'";
// 		 $result = $this->executeQuery($sql);
// 		 if ($result === false)
// 			 return false;
// 		 foreach ($result as $row) {
// 			 $tables[] = $row['TABLE_NAME'];
// 		 }
// 		 return $tables;
// 	 }

// 	 if ($this->database_type == 'sqlite' or $this->database_type == 'sqlitedb') {
// 		 $sql = "SELECT * FROM 'SQLITE_MASTER' WHERE type='table' ORDER BY NAME ASC";
// 		 $result = $this->executeQuery($sql);
// 		 if ($result === false) {
// 			 return false;
// 		 }
// 		 foreach ($result as $row) {
// 			 $tables[] = $row['name'];
// 		 }
// 		 return $tables;
// 	 }
// 	 return false;
//  }
