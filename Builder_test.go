package sb

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	// _ "github.com/glebarez/go-sqlite"
	_ "github.com/mattn/go-sqlite3"
)

func initSqliteDB(filepath string) (DatabaseInterface, error) {
	if filepath == "" {
		return nil, errors.New("filepath is required")
	}

	err := os.Remove(filepath) // remove database

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	sqlDB, err := sql.Open("sqlite3", filepath)

	if err != nil {
		return nil, err
	}

	return NewDatabase(sqlDB, DIALECT_SQLITE), nil
}

func TestBuilderTableCreateMssql(t *testing.T) {
	sql := NewBuilder(DIALECT_MSSQL).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name:   "email",
			Type:   COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := `CREATE TABLE [users] ("id" NVARCHAR(40) PRIMARY KEY NOT NULL, "email" NVARCHAR(255) NOT NULL UNIQUE, "image" VARBINARY(MAX) NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" DATETIME2 NOT NULL, "deleted_at" DATETIME2);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFull(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Where(Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		Select([]string{"id", "first_name", "last_name"})

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> 'Jane' ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name:   "email",
			Type:   COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "short_description",
			Type: COLUMN_TYPE_TEXT,
		}).
		Column(Column{
			Name: "long_description",
			Type: COLUMN_TYPE_LONGTEXT,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `email` VARCHAR(255) NOT NULL UNIQUE, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `short_description` LONGTEXT NOT NULL, `long_description` LONGTEXT NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreatePostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name:   "email",
			Type:   COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "short_description",
			Type: COLUMN_TYPE_TEXT,
		}).
		Column(Column{
			Name: "long_description",
			Type: COLUMN_TYPE_LONGTEXT,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT PRIMARY KEY NOT NULL, "email" TEXT NOT NULL UNIQUE, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "short_description" TEXT NOT NULL, "long_description" TEXT NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name:   "email",
			Type:   COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "short_description",
			Type: COLUMN_TYPE_TEXT,
		}).
		Column(Column{
			Name: "long_description",
			Type: COLUMN_TYPE_LONGTEXT,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT(40) PRIMARY KEY NOT NULL, "email" TEXT(255) NOT NULL UNIQUE, "image" BLOB NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "short_description" TEXT NOT NULL, "long_description" TEXT NOT NULL, "created_at" DATETIME NOT NULL, "deleted_at" DATETIME);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	expected := "CREATE TABLE IF NOT EXISTS `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	expected := `CREATE TABLE IF NOT EXISTS "users"("id" TEXT PRIMARY KEY NOT NULL, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name: "image",
			Type: COLUMN_TYPE_BLOB,
		}).
		Column(Column{
			Name: "price_default",
			Type: COLUMN_TYPE_DECIMAL,
		}).
		Column(Column{
			Name:     "price_custom",
			Type:     COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	expected := `CREATE TABLE IF NOT EXISTS "users"("id" TEXT(40) PRIMARY KEY NOT NULL, "image" BLOB NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" DATETIME NOT NULL, "deleted_at" DATETIME);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateMysql(t *testing.T) {
	selectSQL := NewBuilder(DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})

	sql := NewBuilder(DIALECT_MYSQL).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		Create()

	expected := "CREATE VIEW `v_users` (`first_name`, `last_name`) AS SELECT `FirstName`, `LastName` FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreatePostgresql(t *testing.T) {
	selectSQL := NewBuilder(DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})

	sql := NewBuilder(DIALECT_POSTGRES).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		Create()

	expected := `CREATE VIEW "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateSqlite(t *testing.T) {
	selectSQL := NewBuilder(DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})

	sql := NewBuilder(DIALECT_SQLITE).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		Create()

	expected := `CREATE VIEW "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateIfNotExistsMysql(t *testing.T) {
	selectSQL := NewBuilder(DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})

	sql := NewBuilder(DIALECT_MYSQL).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		CreateIfNotExists()

	expected := "CREATE OR REPLACE VIEW `v_users` (`first_name`, `last_name`) AS SELECT `FirstName`, `LastName` FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateIfNotExistsPostgresql(t *testing.T) {
	selectSQL := NewBuilder(DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})

	sql := NewBuilder(DIALECT_POSTGRES).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		CreateIfNotExists()

	expected := `CREATE VIEW IF NOT EXISTS "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateIfNotExistsSqlite(t *testing.T) {
	selectSQL := NewBuilder(DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})

	sql := NewBuilder(DIALECT_SQLITE).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		CreateIfNotExists()

	expected := `CREATE VIEW IF NOT EXISTS "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		CreateIndex("idx_users_id", "id")

	expected := "CREATE INDEX `idx_users_id` ON `users` (`id`);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		CreateIndex("idx_users_id", "id")

	expected := `CREATE INDEX "idx_users_id" ON "users" ("id");`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		CreateIndex("idx_users_id", "id")

	expected := `CREATE INDEX "idx_users_id" ON "users" ("id");`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableColumnAddSqlite(t *testing.T) {
	db, err := initSqliteDB("test_column_add.db")

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	defer db.Close()

	sqlTableCreate := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name:   "email",
			Type:   COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	result, err := db.Exec(sqlTableCreate)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if result == nil {
		t.Fatal("Result must not be NIL")
	}

	sqlColumnRename, err := NewBuilder(DIALECT_SQLITE).
		TableColumnAdd("users", Column{
			Name:     "name",
			Type:     COLUMN_TYPE_STRING,
			Nullable: true,
		})

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	expected := `ALTER TABLE "users" ADD COLUMN "name" TEXT;`
	if sqlColumnRename != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sqlColumnRename)
	}

	result, err = db.Exec(sqlColumnRename)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if result == nil {
		t.Fatal("Result must not be NIL")
	}

	sql := NewBuilder(DIALECT_SQLITE).Table("users").Select([]string{"id", "email", "name", "created_at", "deleted_at"})

	rows, err := db.Query(sql)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if rows == nil {
		t.Fatal("Rows must not be NIL")
	}
}

func TestBuilderTableColumnRenameSqlite(t *testing.T) {
	db, err := initSqliteDB("test_column_rename.db")

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	defer db.Close()

	sqlTableCreate := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Column(Column{
			Name:       "id",
			Type:       COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(Column{
			Name:   "email",
			Type:   COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(Column{
			Name: "created_at",
			Type: COLUMN_TYPE_DATETIME,
		}).
		Column(Column{
			Name:     "deleted_at",
			Type:     COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	result, err := db.Exec(sqlTableCreate)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if result == nil {
		t.Fatal("Result must not be NIL")
	}

	sqlColumnRename, err := NewBuilder(DIALECT_SQLITE).
		TableColumnRename("users", "email", "name")

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	expected := `ALTER TABLE "users" RENAME COLUMN "email" TO "name";`
	if sqlColumnRename != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sqlColumnRename)
	}

	result, err = db.Exec(sqlColumnRename)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if result == nil {
		t.Fatal("Result must not be NIL")
	}

	sql := NewBuilder(DIALECT_SQLITE).Table("users").Select([]string{"id", "name", "created_at", "deleted_at"})

	rows, err := db.Query(sql)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if rows == nil {
		t.Fatal("Rows must not be NIL")
	}
}

func TestBuilderTableDropMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Drop()

	expected := "DROP TABLE `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Drop()

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Drop()

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Delete()

	expected := "DELETE FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysqlExtended(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Sam",
			Type:     "OR",
		}).
		Limit(12).
		Offset(34).
		Delete()

	expected := "DELETE FROM `users` WHERE `FirstName` = \"Tom\" OR `FirstName` = \"Sam\" LIMIT 12 OFFSET 34;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Delete()

	expected := `DELETE FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqliteExtended(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Sam",
			Type:     "OR",
		}).
		Delete()

	expected := `DELETE FROM "users" WHERE "FirstName" = 'Tom' OR "FirstName" = 'Sam';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Select([]string{})

	expected := "SELECT * FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Select([]string{})

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Select([]string{})

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Where(Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(GroupBy{Column: "passport"}).
		Select([]string{"id", "first_name", "last_name"})

	expected := "SELECT `id`, `first_name`, `last_name` FROM `users` WHERE `first_name` <> \"Jane\" GROUP BY `passport` ORDER BY `first_name` ASC LIMIT 10 OFFSET 20;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Where(Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(GroupBy{Column: "passport"}).
		Select([]string{"id", "first_name", "last_name"})

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> "Jane" GROUP BY "passport" ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Where(Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(GroupBy{Column: "passport"}).
		Select([]string{"id", "first_name", "last_name"})

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> 'Jane' GROUP BY "passport" ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Limit(1).
		Insert(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	expected := "INSERT INTO `users` (`first_name`, `last_name`) VALUES (\"Tom\", \"Jones\") LIMIT 1;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertPostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Limit(1).
		Insert(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	expected := `INSERT INTO "users" ("first_name", "last_name") VALUES ("Tom", "Jones") LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Limit(1).
		Insert(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	expected := `INSERT INTO "users" ("first_name", "last_name") VALUES ('Tom', 'Jones') LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableColumnCreateSqlite(t *testing.T) {
	sql, err := NewBuilder(DIALECT_SQLITE).TableColumnAdd("table_name", Column{
		Name:     "name",
		Type:     COLUMN_TYPE_STRING,
		Length:   255,
		Nullable: true,
	})

	if err != nil {
		t.Fatal(err)
	}

	expected := `ALTER TABLE "table_name" ADD COLUMN "name" TEXT(255);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdateMysql(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Where(Where{
			Column:   "id",
			Operator: "==",
			Value:    "1",
		}).
		Limit(1).
		Update(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	expected := "UPDATE `users` SET `first_name`=\"Tom\", `last_name`=\"Jones\" WHERE `id` = \"1\" LIMIT 1;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdatePostgres(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Where(Where{
			Column:   "id",
			Operator: "==",
			Value:    "1",
		}).
		Limit(1).
		Update(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	expected := `UPDATE "users" SET "first_name"="Tom", "last_name"="Jones" WHERE "id" = "1" LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdateSqlite(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Where(Where{
			Column:   "id",
			Operator: "==",
			Value:    "1",
		}).
		Limit(1).
		Update(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	expected := `UPDATE "users" SET "first_name"='Tom', "last_name"='Jones' WHERE "id" = '1' LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectMysqlInj(t *testing.T) {
	sql := NewBuilder(DIALECT_MYSQL).
		Table("users").
		Where(Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		Select([]string{})

	expected := "SELECT * FROM `users` WHERE `id` = \"58\"\" OR 1 = 1;--\";"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgreslInj(t *testing.T) {
	sql := NewBuilder(DIALECT_POSTGRES).
		Table("users").
		Where(Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		Select([]string{})

	expected := `SELECT * FROM "users" WHERE "id" = "58"" OR 1 = 1;--";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlitelInj(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Where(Where{Column: "id", Operator: "=", Value: "58' OR 1 = 1;--"}).
		Select([]string{})

	expected := `SELECT * FROM "users" WHERE "id" = '58'' OR 1 = 1;--';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectAll(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Select([]string{"*"})

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFn(t *testing.T) {
	sql := NewBuilder(DIALECT_SQLITE).
		Table("users").
		Select([]string{"MIN(created_at)"})

	expected := `SELECT MIN(created_at) FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}
