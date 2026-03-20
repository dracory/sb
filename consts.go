package sb

// Dialects
const DIALECT_MSSQL = "mssql"
const DIALECT_MYSQL = "mysql"
const DIALECT_POSTGRES = "postgres"
const DIALECT_SQLITE = "sqlite"

// Column Attributes
const COLUMN_ATTRIBUTE_AUTO = "auto"
const COLUMN_ATTRIBUTE_DECIMALS = "decimals"
const COLUMN_ATTRIBUTE_LENGTH = "length"
const COLUMN_ATTRIBUTE_NULLABLE = "nullable"
const COLUMN_ATTRIBUTE_PRIMARY = "primary"

// Column Types
const COLUMN_TYPE_BLOB = "blob"
const COLUMN_TYPE_DATE = "date"
const COLUMN_TYPE_DATETIME = "datetime"
const COLUMN_TYPE_DECIMAL = "decimal"
const COLUMN_TYPE_FLOAT = "float"
const COLUMN_TYPE_INTEGER = "integer"
const COLUMN_TYPE_STRING = "string"
const COLUMN_TYPE_TEXT = "text"
const COLUMN_TYPE_LONGTEXT = "longtext"

// Common
const YES = "yes"
const NO = "no"

// Null time (earliest valid date in Gregorian calendar is 1AD, no year 0)
const NULL_DATE = "0002-01-01"
const NULL_DATETIME = "0002-01-01 00:00:00"

const MAX_DATE = "9999-12-31"
const MAX_DATETIME = "9999-12-31 23:59:59"

// Sortable
const ASC = "asc"
const DESC = "desc"

// Index Types
const INDEX_TYPE_BTREE = "BTREE"
const INDEX_TYPE_HASH = "HASH"
const INDEX_TYPE_GIN = "GIN"
const INDEX_TYPE_GIST = "GIST"
const INDEX_TYPE_BRIN = "BRIN"
const INDEX_TYPE_FULLTEXT = "FULLTEXT"
const INDEX_TYPE_SPATIAL = "SPATIAL"

// Join Types
const JOIN_INNER = "INNER"
const JOIN_LEFT = "LEFT"
const JOIN_RIGHT = "RIGHT"
const JOIN_FULL = "FULL"
const JOIN_CROSS = "CROSS"
