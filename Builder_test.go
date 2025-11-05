package sb_test

import (
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/dracory/sb"
	_ "modernc.org/sqlite"
)

func initSqliteDB(filepath string) (sb.DatabaseInterface, error) {
	if filepath == "" {
		return nil, errors.New("filepath is required")
	}

	err := os.Remove(filepath) // remove database

	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	sqlDB, err := sql.Open("sqlite", filepath)

	if err != nil {
		return nil, err
	}

	return sb.NewDatabase(sqlDB, sb.DIALECT_SQLITE), nil
}

func TestBuilderTableCreateMssql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "email",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := `CREATE TABLE [users] ("id" NVARCHAR(40) PRIMARY KEY NOT NULL, "email" NVARCHAR(255) NOT NULL UNIQUE, "image" VARBINARY(MAX) NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" DATETIME2 NOT NULL, "deleted_at" DATETIME2);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFull(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
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
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "email",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "short_description",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "long_description",
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `email` VARCHAR(255) NOT NULL UNIQUE, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `short_description` LONGTEXT NOT NULL, `long_description` LONGTEXT NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreatePostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "email",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "short_description",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "long_description",
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT PRIMARY KEY NOT NULL, "email" TEXT NOT NULL UNIQUE, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "short_description" TEXT NOT NULL, "long_description" TEXT NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "email",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "short_description",
			Type: sb.COLUMN_TYPE_TEXT,
		}).
		Column(sb.Column{
			Name: "long_description",
			Type: sb.COLUMN_TYPE_LONGTEXT,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		Create()

	expected := `CREATE TABLE "users"("id" TEXT(40) PRIMARY KEY NOT NULL, "email" TEXT(255) NOT NULL UNIQUE, "image" BLOB NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "short_description" TEXT NOT NULL, "long_description" TEXT NOT NULL, "created_at" DATETIME NOT NULL, "deleted_at" DATETIME);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsMysql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	expected := "CREATE TABLE IF NOT EXISTS `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsPostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	expected := `CREATE TABLE IF NOT EXISTS "users"("id" TEXT PRIMARY KEY NOT NULL, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name: "image",
			Type: sb.COLUMN_TYPE_BLOB,
		}).
		Column(sb.Column{
			Name: "price_default",
			Type: sb.COLUMN_TYPE_DECIMAL,
		}).
		Column(sb.Column{
			Name:     "price_custom",
			Type:     sb.COLUMN_TYPE_DECIMAL,
			Length:   12,
			Decimals: 10,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
			Nullable: true,
		}).
		CreateIfNotExists()

	expected := `CREATE TABLE IF NOT EXISTS "users"("id" TEXT(40) PRIMARY KEY NOT NULL, "image" BLOB NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" DATETIME NOT NULL, "deleted_at" DATETIME);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableColumnChangeMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).TableColumnChange("users", sb.Column{
		Name:   "email",
		Type:   sb.COLUMN_TYPE_STRING,
		Length: 255,
		Unique: true,
	})

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	expected := "ALTER TABLE `users` MODIFY COLUMN `email` VARCHAR(255) NOT NULL UNIQUE;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableColumnChangePostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).TableColumnChange("users", sb.Column{
		Name:   "email",
		Type:   sb.COLUMN_TYPE_STRING,
		Length: 255,
		Unique: true,
	})

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	expected := `ALTER TABLE "users" ALTER COLUMN "email" TEXT NOT NULL UNIQUE;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableColumnChangeSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).TableColumnChange("users", sb.Column{
		Name:   "email",
		Type:   sb.COLUMN_TYPE_STRING,
		Length: 255,
		Unique: true,
	})

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	expected := `ALTER TABLE "users" ALTER COLUMN "email" TEXT(255) NOT NULL UNIQUE;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateMysql(t *testing.T) {
	selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})

	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
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
	selectSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})

	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
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
	selectSQL := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})

	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
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
	selectSQL := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})

	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
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
	selectSQL := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})

	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
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
	selectSQL := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})

	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
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
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndex("idx_users_id", "id")

	expected := "CREATE INDEX `idx_users_id` ON `users` (`id`);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexPostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		CreateIndex("idx_users_id", "id")

	expected := `CREATE INDEX "idx_users_id" ON "users" ("id");`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	sqlTableCreate := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "email",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
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

	sqlColumnRename, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		TableColumnAdd("users", sb.Column{
			Name:     "name",
			Type:     sb.COLUMN_TYPE_STRING,
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

	sql := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "email", "name", "created_at", "deleted_at"})

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

	sqlTableCreate := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Column(sb.Column{
			Name:       "id",
			Type:       sb.COLUMN_TYPE_STRING,
			Length:     40,
			PrimaryKey: true,
		}).
		Column(sb.Column{
			Name:   "email",
			Type:   sb.COLUMN_TYPE_STRING,
			Length: 255,
			Unique: true,
		}).
		Column(sb.Column{
			Name: "created_at",
			Type: sb.COLUMN_TYPE_DATETIME,
		}).
		Column(sb.Column{
			Name:     "deleted_at",
			Type:     sb.COLUMN_TYPE_DATETIME,
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

	sqlColumnRename, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	sql := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "name", "created_at", "deleted_at"})

	rows, err := db.Query(sql)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if rows == nil {
		t.Fatal("Rows must not be NIL")
	}
}

func TestBuilderTableDropMysql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Drop()

	expected := "DROP TABLE `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropPostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Drop()

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Drop()

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Delete()

	expected := "DELETE FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysqlExtended(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(sb.Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(sb.Where{
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
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Delete()

	expected := `DELETE FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqliteExtended(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(sb.Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(sb.Where{
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
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Select([]string{})

	expected := "SELECT * FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Select([]string{})

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Select([]string{})

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullMysql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
		Select([]string{"id", "first_name", "last_name"})

	expected := "SELECT `id`, `first_name`, `last_name` FROM `users` WHERE `first_name` <> \"Jane\" GROUP BY `passport` ORDER BY `first_name` ASC LIMIT 10 OFFSET 20;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullPostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
		Select([]string{"id", "first_name", "last_name"})

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> "Jane" GROUP BY "passport" ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
		Select([]string{"id", "first_name", "last_name"})

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> 'Jane' GROUP BY "passport" ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertMysql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
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
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
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
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).TableColumnAdd("table_name", sb.Column{
		Name:     "name",
		Type:     sb.COLUMN_TYPE_STRING,
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
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(sb.Where{
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
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(sb.Where{
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
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(sb.Where{
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
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		Select([]string{})

	expected := "SELECT * FROM `users` WHERE `id` = \"58\"\" OR 1 = 1;--\";"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgreslInj(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		Select([]string{})

	expected := `SELECT * FROM "users" WHERE "id" = "58"" OR 1 = 1;--";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlitelInj(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(sb.Where{Column: "id", Operator: "=", Value: "58' OR 1 = 1;--"}).
		Select([]string{})

	expected := `SELECT * FROM "users" WHERE "id" = '58'' OR 1 = 1;--';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectAll(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Select([]string{"*"})

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFn(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Select([]string{"MIN(created_at)"})

	expected := `SELECT MIN(created_at) FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewDropMysql(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		View("v_users").
		Drop()

	expected := "DROP VIEW `v_users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewDropPostgres(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		View("v_users").
		Drop()

	expected := `DROP VIEW "v_users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewDropSqlite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		View("v_users").
		Drop()

	expected := `DROP VIEW "v_users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilder_TableColumnChange(t *testing.T) {
	type args struct {
		tableName string
		column    sb.Column
	}
	tests := []struct {
		name          string
		b             *sb.Builder
		args          args
		wantSqlString string
		wantErr       bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotSqlString, err := tt.b.TableColumnChange(tt.args.tableName, tt.args.column)
			if (err != nil) != tt.wantErr {
				t.Errorf("Builder.TableColumnChange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotSqlString != tt.wantSqlString {
				t.Errorf("Builder.TableColumnChange() = %v, want %v", gotSqlString, tt.wantSqlString)
			}
		})
	}
}

func TestSQLiteAutoIncrementOrder(t *testing.T) {
	// Test that SQLite generates correct PRIMARY KEY AUTOINCREMENT order
	sqlQuery := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("test_table").
		Column(sb.Column{
			Name:          "id",
			Type:          sb.COLUMN_TYPE_INTEGER,
			PrimaryKey:    true,
			AutoIncrement: true,
			Nullable:      false,
		}).
		Column(sb.Column{
			Name:     "name",
			Type:     sb.COLUMN_TYPE_STRING,
			Nullable: false,
		}).
		CreateIfNotExists()

	expected := `CREATE TABLE IF NOT EXISTS "test_table"("id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL, "name" TEXT NOT NULL);`
	if sqlQuery != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sqlQuery)
	}

	// Test that the generated SQL actually works with SQLite
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(sqlQuery)
	if err != nil {
		t.Fatalf("Failed to execute generated SQL: %v", err)
	}

	// Test inserting data to verify autoincrement works
	_, err = db.Exec(`INSERT INTO test_table (name) VALUES ('test1'), ('test2')`)
	if err != nil {
		t.Fatalf("Failed to insert test data: %v", err)
	}

	// Verify the autoincrement worked
	rows, err := db.Query(`SELECT id, name FROM test_table ORDER BY id`)
	if err != nil {
		t.Fatalf("Failed to query test data: %v", err)
	}
	defer rows.Close()

	expectedData := []struct {
		id   int
		name string
	}{
		{1, "test1"},
		{2, "test2"},
	}

	i := 0
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		if i >= len(expectedData) {
			t.Fatalf("More rows than expected")
		}

		if id != expectedData[i].id || name != expectedData[i].name {
			t.Fatalf("Row %d: expected id=%d, name=%s but got id=%d, name=%s",
				i, expectedData[i].id, expectedData[i].name, id, name)
		}
		i++
	}

	if i != len(expectedData) {
		t.Fatalf("Expected %d rows but got %d", len(expectedData), i)
	}
}
