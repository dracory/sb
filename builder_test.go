package sb_test

import (
	"context"
	stdsql "database/sql"
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

	sqlDB, err := stdsql.Open("sqlite", filepath)

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
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
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

	defer func() {
		_ = db.Close()
	}()

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

	result, err := db.Exec(context.Background(), sqlTableCreate)

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

	result, err = db.Exec(context.Background(), sqlColumnRename)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if result == nil {
		t.Fatal("Result must not be NIL")
	}

	sql := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "email", "name", "created_at", "deleted_at"})

	rows, err := db.Query(context.Background(), sql)

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

	defer func() {
		_ = db.Close()
	}()

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

	result, err := db.Exec(context.Background(), sqlTableCreate)

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

	result, err = db.Exec(context.Background(), sqlColumnRename)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if result == nil {
		t.Fatal("Result must not be NIL")
	}

	sql := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "name", "created_at", "deleted_at"})

	rows, err := db.Query(context.Background(), sql)

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
		Where(&sb.Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(&sb.Where{
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
		Where(&sb.Where{
			Column:   "FirstName",
			Operator: "==",
			Value:    "Tom",
		}).
		Where(&sb.Where{
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
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
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
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
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
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
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
		Where(&sb.Where{
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
		Where(&sb.Where{
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
		Where(&sb.Where{
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
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		Select([]string{})

	expected := "SELECT * FROM `users` WHERE `id` = \"58\"\" OR 1 = 1;--\";"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgreslInj(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		Select([]string{})

	expected := `SELECT * FROM "users" WHERE "id" = "58"" OR 1 = 1;--";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlitelInj(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58' OR 1 = 1;--"}).
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
	db, err := stdsql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open SQLite database: %v", err)
	}
	defer func() {
		_ = db.Close()
	}()

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
	defer func() {
		_ = rows.Close()
	}()

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

// Test Builder Truncate functionality

func TestBuilderTruncateMySQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Truncate()

	expected := "TRUNCATE TABLE `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncatePostgreSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Truncate()

	expected := `TRUNCATE TABLE "users";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateSQLite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Truncate()

	expected := `DELETE FROM "users";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateMSSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		Truncate()

	expected := "TRUNCATE TABLE [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsMySQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true, ResetIdentity: true})

	expected := "TRUNCATE TABLE `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsPostgreSQL(t *testing.T) {
	// Test without CASCADE
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		TruncateWithOptions(sb.TruncateOptions{})

	expected := `TRUNCATE TABLE "orders";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test with CASCADE
	sql = sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true})

	expected = `TRUNCATE TABLE "orders" CASCADE;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsSQLite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("cache").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true, ResetIdentity: true})

	expected := `DELETE FROM "cache";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsMSSQL(t *testing.T) {
	// Test without ResetIdentity
	sql := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{})

	expected := "TRUNCATE TABLE [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test with ResetIdentity
	sql = sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{ResetIdentity: true})

	expected = "TRUNCATE TABLE [users]; DBCC CHECKIDENT ('users', RESEED, 0)"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithComplexTableName(t *testing.T) {
	// Test with schema.table format
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("public.users").
		Truncate()

	expected := `TRUNCATE TABLE "public"."users";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test with database.schema.table format
	var sql2 string
	sql2 = sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("dbo.users").
		Truncate()

	expected = "TRUNCATE TABLE [dbo].[users];"
	if sql2 != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql2)
	}
}

func TestBuilderTruncateErrorHandling(t *testing.T) {
	// Test panic when no table is specified
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when no table specified")
		}

		expectedMsg := "In method Truncate() no table specified to truncate!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).Truncate()
}

func TestBuilderTruncateWithOptionsErrorHandling(t *testing.T) {
	// Test panic when no table is specified
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when no table specified")
		}

		expectedMsg := "In method TruncateWithOptions() no table specified to truncate!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).TruncateWithOptions(sb.TruncateOptions{})
}

func TestBuilderTruncateUnsupportedDialect(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic for unsupported dialect")
		}

		expectedMsg := "unsupported dialect: unknown"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder("unknown").
		Table("users").
		Truncate()
}

// Test Builder DropIndex functionality

func TestBuilderDropIndexMySQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndex("idx_users_email")

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexPostgreSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndex("idx_users_email")

	expected := `DROP INDEX "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexSQLite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		DropIndex("idx_users_email")

	expected := `DROP INDEX "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexMSSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		DropIndex("idx_users_email")

	expected := "DROP INDEX [idx_users_email] ON [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsMySQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexIfExists("idx_users_email")

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsPostgreSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexIfExists("idx_users_email")

	expected := `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsSQLite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		DropIndexIfExists("idx_users_email")

	expected := `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsMSSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		DropIndexIfExists("idx_users_email")

	expected := "DROP INDEX IF EXISTS [idx_users_email] ON [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexWithSchemaPostgreSQL(t *testing.T) {
	// Test with schema
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	expected := `DROP INDEX IF EXISTS "public"."idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test without schema (should use regular behavior)
	sql = sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithSchema("idx_users_email", "")

	expected = `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexWithSchemaOtherDialects(t *testing.T) {
	// MySQL should fall back to regular DropIndex
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// SQLite should fall back to regular DropIndex
	sql = sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	expected = `DROP INDEX "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// MSSQL should fall back to regular DropIndex
	sql = sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	expected = "DROP INDEX [idx_users_email] ON [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexErrorHandling(t *testing.T) {
	// Test panic when no table is specified
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when no table specified")
		}

		expectedMsg := "In method DropIndex() no table specified to drop index from!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).DropIndex("idx_users_email")
}

func TestBuilderDropIndexEmptyIndexName(t *testing.T) {
	// Test panic when index name is empty
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when index name is empty")
		}

		expectedMsg := "In method DropIndex() index name cannot be empty!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndex("")
}

func TestBuilderDropIndexIfExistsErrorHandling(t *testing.T) {
	// Test panic when no table is specified
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when no table specified")
		}

		expectedMsg := "In method DropIndexIfExists() no table specified to drop index from!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).DropIndexIfExists("idx_users_email")
}

func TestBuilderDropIndexWithSchemaErrorHandling(t *testing.T) {
	// Test panic when no table is specified
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when no table specified")
		}

		expectedMsg := "In method DropIndexWithSchema() no table specified to drop index from!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_POSTGRES).DropIndexWithSchema("idx_users_email", "public")
}

func TestBuilderDropIndexUnsupportedDialect(t *testing.T) {
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic for unsupported dialect")
		}

		expectedMsg := "unsupported dialect: unknown"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder("unknown").
		Table("users").
		DropIndex("idx_users_email")
}

// Test Builder JOIN functionality

func TestBuilderJoinInnerMySQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})

	expected := "SELECT `orders`.*, `users`.`name` FROM `orders`  INNER JOIN `users` ON orders.user_id = users.id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinLeftPostgreSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		LeftJoin("profiles", "orders.user_id = profiles.user_id").
		Select([]string{"orders.*", "profiles.avatar"})

	expected := `SELECT "orders".*, "profiles"."avatar" FROM "orders"  LEFT JOIN "profiles" ON orders.user_id = profiles.user_id;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinRightMSSQL(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("orders").
		RightJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})

	expected := "SELECT [orders].*, [users].[name] FROM [orders]  RIGHT JOIN [users] ON orders.user_id = users.id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinWithAliasSQLite(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("orders").
		JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
		Select([]string{"orders.*", "p.avatar"})

	expected := `SELECT "orders".*, "p"."avatar" FROM "orders"  LEFT JOIN "profiles" AS "p" ON orders.user_id = p.user_id;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinMultiple(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		LeftJoin("profiles", "users.id = profiles.user_id").
		Select([]string{"orders.total", "users.name", "profiles.avatar"})

	expected := "SELECT `orders`.`total`, `users`.`name`, `profiles`.`avatar` FROM `orders`  INNER JOIN `users` ON orders.user_id = users.id  LEFT JOIN `profiles` ON users.id = profiles.user_id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinFullOuter(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Join(sb.JOIN_FULL, "users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})

	expected := `SELECT "orders".*, "users"."name" FROM "orders"  FULL JOIN "users" ON orders.user_id = users.id;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinCross(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("orders").
		Join(sb.JOIN_CROSS, "users", "1=1").
		Select([]string{"orders.*", "users.name"})

	expected := `SELECT "orders".*, "users"."name" FROM "orders"  CROSS JOIN "users" ON 1=1;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinWithWhere(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
		Select([]string{"orders.*", "users.name"})

	expected := "SELECT `orders`.*, `users`.`name` FROM `orders`  INNER JOIN `users` ON orders.user_id = users.id WHERE `orders`.`status` = \"active\";"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinWithOrderBy(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		LeftJoin("users", "orders.user_id = users.id").
		OrderBy("orders.created_at", "DESC").
		Select([]string{"orders.*", "users.name"})

	expected := `SELECT "orders".*, "users"."name" FROM "orders"  LEFT JOIN "users" ON orders.user_id = users.id ORDER BY "orders"."created_at" DESC;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinErrorHandlingEmptyCondition(t *testing.T) {
	// Test panic when ON condition is empty
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when ON condition is empty")
		}

		expectedMsg := "In method Join() ON condition cannot be empty!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "")
}

func TestBuilderJoinWithAliasErrorHandlingEmptyCondition(t *testing.T) {
	// Test panic when ON condition is empty
	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("Expected panic when ON condition is empty")
		}

		expectedMsg := "In method JoinWithAlias() ON condition cannot be empty!"
		if recovered != expectedMsg {
			t.Fatalf("Expected panic message: %s but got: %v", expectedMsg, recovered)
		}
	}()

	sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "")
}

func TestBuilderJoinComplexTableName(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("public.users", "orders.user_id = public.users.id").
		Select([]string{"orders.*", "public.users.name"})

	expected := "SELECT `orders`.*, `public`.`users`.`name` FROM `orders`  INNER JOIN `public`.`users` ON orders.user_id = public.users.id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinUpdate(t *testing.T) {
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Where(&sb.Where{Column: "users.id", Operator: "=", Value: "1"}).
		Update(map[string]string{"orders.status": "processed"})

	expected := "UPDATE `orders` SET `orders`.`status`=\"processed\"  INNER JOIN `users` ON orders.user_id = users.id WHERE `users`.`id` = \"1\";"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinNoJoins(t *testing.T) {
	// Test that queries work normally without joins
	sql := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		Select([]string{"*"})

	expected := "SELECT * FROM `orders` WHERE `status` = \"active\";"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}
