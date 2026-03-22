package sb_test

import (
	"context"
	stdsql "database/sql"
	"errors"
	"os"
	"strings"
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
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE TABLE [users] ("id" NVARCHAR(40) PRIMARY KEY NOT NULL, "email" NVARCHAR(255) NOT NULL UNIQUE, "image" VARBINARY(MAX) NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" DATETIME2 NOT NULL, "deleted_at" DATETIME2);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFull(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		WithInterpolatedValues().
		Select([]string{"id", "first_name", "last_name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> 'Jane' ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE TABLE `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `email` VARCHAR(255) NOT NULL UNIQUE, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `short_description` LONGTEXT NOT NULL, `long_description` LONGTEXT NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreatePostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE TABLE "users"("id" TEXT PRIMARY KEY NOT NULL, "email" TEXT NOT NULL UNIQUE, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "short_description" TEXT NOT NULL, "long_description" TEXT NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE TABLE "users"("id" TEXT(40) PRIMARY KEY NOT NULL, "email" TEXT(255) NOT NULL UNIQUE, "image" BLOB NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "short_description" TEXT NOT NULL, "long_description" TEXT NOT NULL, "created_at" DATETIME NOT NULL, "deleted_at" DATETIME);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE TABLE IF NOT EXISTS `users`(`id` VARCHAR(40) PRIMARY KEY NOT NULL, `image` LONGBLOB NOT NULL, `price_default` DECIMAL(10,2) NOT NULL, `price_custom` DECIMAL(12,10) NOT NULL, `created_at` DATETIME NOT NULL, `deleted_at` DATETIME);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\n but found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsPostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE TABLE IF NOT EXISTS "users"("id" TEXT PRIMARY KEY NOT NULL, "image" BYTEA NOT NULL, "price_default" DECIMAL(10,2) NOT NULL, "price_custom" DECIMAL(12,10) NOT NULL, "created_at" TIMESTAMP NOT NULL, "deleted_at" TIMESTAMP);`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableCreateIfNotExistsSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

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
	selectSQL, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})
	if err != nil {
		t.Fatal("Unexpected error in Select:", err)
	}

	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		Create()

	if err != nil {
		t.Fatal("Unexpected error in Create:", err)
	}

	expected := "CREATE VIEW `v_users` (`first_name`, `last_name`) AS SELECT `FirstName`, `LastName` FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreatePostgresql(t *testing.T) {
	selectSQL, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		Create()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE VIEW "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateSqlite(t *testing.T) {
	selectSQL, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		Create()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE VIEW "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateIfNotExistsMysql(t *testing.T) {
	selectSQL, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})
	if err != nil {
		t.Fatal("Unexpected error in Select:", err)
	}

	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		CreateIfNotExists()

	if err != nil {
		t.Fatal("Unexpected error in Create:", err)
	}

	expected := "CREATE OR REPLACE VIEW `v_users` (`first_name`, `last_name`) AS SELECT `FirstName`, `LastName` FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateIfNotExistsPostgresql(t *testing.T) {
	selectSQL, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		CreateIfNotExists()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE VIEW IF NOT EXISTS "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewCreateIfNotExistsSqlite(t *testing.T) {
	selectSQL, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})
	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		View("v_users").
		ViewColumns([]string{"first_name", "last_name"}).
		ViewSQL(selectSQL).
		CreateIfNotExists()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE VIEW IF NOT EXISTS "v_users" ("first_name", "last_name") AS SELECT "FirstName", "LastName" FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndex("idx_users_id", "id")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_id` ON `users` (`id`);"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexPostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		CreateIndex("idx_users_id", "id")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE INDEX "idx_users_id" ON "users" ("id");`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderCreateIndexSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		CreateIndex("idx_users_id", "id")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

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

	sqlTableCreate, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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

	if err != nil {
		t.Fatal("Error creating table:", err)
	}

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

	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "email", "name", "created_at", "deleted_at"})

	if err != nil {
		t.Fatal("Unexpected error in Select:", err)
	}

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

	sqlTableCreate, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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

	if err != nil {
		t.Fatal("Error creating table:", err)
	}

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

	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "name", "created_at", "deleted_at"})

	if err != nil {
		t.Fatal("Unexpected error in Select:", err)
	}

	rows, err := db.Query(context.Background(), sql)

	if err != nil {
		t.Fatal("Error must be NIL but got: ", err.Error())
	}

	if rows == nil {
		t.Fatal("Rows must not be NIL")
	}
}

func TestBuilderTableDropMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Drop()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP TABLE `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropPostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Drop()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDropSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Drop()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP TABLE "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysql(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Delete()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DELETE FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteMysqlExtended(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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
		WithInterpolatedValues().
		Delete()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DELETE FROM `users` WHERE `FirstName` = \"Tom\" OR `FirstName` = \"Sam\" LIMIT 12 OFFSET 34;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqlite(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Delete()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DELETE FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqliteExtended(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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
		WithInterpolatedValues().
		Delete()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DELETE FROM "users" WHERE "FirstName" = 'Tom' OR "FirstName" = 'Sam';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectMysql(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Select([]string{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT * FROM `users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgres(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Select([]string{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlite(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Select([]string{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullMysql(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
		WithInterpolatedValues().
		Select([]string{"id", "first_name", "last_name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `id`, `first_name`, `last_name` FROM `users` WHERE `first_name` <> \"Jane\" GROUP BY `passport` ORDER BY `first_name` ASC LIMIT 10 OFFSET 20;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullPostgres(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
		WithInterpolatedValues().
		Select([]string{"id", "first_name", "last_name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> 'Jane' GROUP BY "passport" ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFullSqlite(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
		WithInterpolatedValues().
		Select([]string{"id", "first_name", "last_name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "id", "first_name", "last_name" FROM "users" WHERE "first_name" <> 'Jane' GROUP BY "passport" ORDER BY "first_name" ASC LIMIT 10 OFFSET 20;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertMysql(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Limit(1).
		WithInterpolatedValues().
		Insert(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "INSERT INTO `users` (`first_name`, `last_name`) VALUES (\"Tom\", \"Jones\") LIMIT 1;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertPostgres(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Limit(1).
		WithInterpolatedValues().
		Insert(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `INSERT INTO "users" ("first_name", "last_name") VALUES ('Tom', 'Jones') LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableInsertSqlite(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Limit(1).
		WithInterpolatedValues().
		Insert(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

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
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{
			Column:   "id",
			Operator: "==",
			Value:    "1",
		}).
		Limit(1).
		WithInterpolatedValues().
		Update(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "UPDATE `users` SET `first_name`=\"Tom\", `last_name`=\"Jones\" WHERE `id` = \"1\" LIMIT 1;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdatePostgres(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{
			Column:   "id",
			Operator: "==",
			Value:    "1",
		}).
		Limit(1).
		WithInterpolatedValues().
		Update(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `UPDATE "users" SET "first_name"='Tom', "last_name"='Jones' WHERE "id" = '1' LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdateSqlite(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{
			Column:   "id",
			Operator: "==",
			Value:    "1",
		}).
		Limit(1).
		WithInterpolatedValues().
		Update(map[string]string{
			"first_name": "Tom",
			"last_name":  "Jones",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `UPDATE "users" SET "first_name"='Tom', "last_name"='Jones' WHERE "id" = '1' LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectMysqlInj(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		WithInterpolatedValues().
		Select([]string{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT * FROM `users` WHERE `id` = \"58\"\" OR 1 = 1;--\";"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectPostgreslInj(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
		WithInterpolatedValues().
		Select([]string{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT * FROM "users" WHERE "id" = '58" OR 1 = 1;--';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectSqlitelInj(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58' OR 1 = 1;--"}).
		WithInterpolatedValues().
		Select([]string{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT * FROM "users" WHERE "id" = '58'' OR 1 = 1;--';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectAll(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Select([]string{"*"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT * FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectFn(t *testing.T) {
	sql, _, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Select([]string{"MIN(created_at)"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT MIN(created_at) FROM "users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewDropMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		View("v_users").
		Drop()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP VIEW `v_users`;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewDropPostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		View("v_users").
		Drop()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP VIEW "v_users";`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderViewDropSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		View("v_users").
		Drop()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

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
	sqlQuery, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

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
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Truncate()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "TRUNCATE TABLE `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncatePostgreSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Truncate()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `TRUNCATE TABLE "users";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateSQLite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Truncate()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DELETE FROM "users";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateMSSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		Truncate()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "TRUNCATE TABLE [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsMySQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true, ResetIdentity: true})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "TRUNCATE TABLE `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsPostgreSQL(t *testing.T) {
	// Test without CASCADE
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		TruncateWithOptions(sb.TruncateOptions{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `TRUNCATE TABLE "orders"`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test with CASCADE
	sql, err = sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected = `TRUNCATE TABLE "orders" CASCADE`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsSQLite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("cache").
		TruncateWithOptions(sb.TruncateOptions{Cascade: true, ResetIdentity: true})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DELETE FROM "cache";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithOptionsMSSQL(t *testing.T) {
	// Test without ResetIdentity
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "TRUNCATE TABLE [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryCreateMethod(t *testing.T) {
	subqueryBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).Subquery()

	subquery := subqueryBuilder.
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"}).
		WithInterpolatedValues()

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, _, err := builder.WithInterpolatedValues().Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > \"1000\");"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestSubqueryValidation(t *testing.T) {
	tests := []struct {
		name        string
		operator    string
		columns     []string
		shouldError bool
		errorMsg    string
	}{
		{"IN with single column", "IN", []string{"id"}, false, ""},
		{"IN with multiple columns", "IN", []string{"id", "name"}, true, "IN/NOT IN subquery must select exactly one column"},
		{"EXISTS with multiple columns", "EXISTS", []string{"*"}, false, ""},
		{"Comparison with multiple columns", ">", []string{"id", "name"}, true, "comparison subquery must select exactly one column"},
		{"NOT IN with single column", "NOT IN", []string{"user_id"}, false, ""},
		{"NOT IN with multiple columns", "NOT IN", []string{"user_id", "total"}, true, "IN/NOT IN subquery must select exactly one column"},
		{"Equals with single column", "=", []string{"max_id"}, false, ""},
		{"Equals with multiple columns", "=", []string{"max_id", "created_at"}, true, "comparison subquery must select exactly one column"},
		{"Unknown operator", "UNKNOWN", []string{"id", "name"}, false, ""}, // Should not error for unknown operators
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a subquery with the specified columns
			subqueryInterface := sb.NewBuilder(sb.DIALECT_MYSQL).
				Table("orders")

			// Set the columns by calling Select to trigger sqlSelectColumns setting
			_, _, err := subqueryInterface.Select(tt.columns)
			if err != nil {
				t.Fatal("Unexpected error setting subquery columns:", err)
			}

			// Try to use the subquery
			if tt.operator == "IN" || tt.operator == "NOT IN" {
				_, _, err = sb.NewBuilder(sb.DIALECT_MYSQL).
					Table("users").
					Where(&sb.Where{
						Operator: tt.operator,
						Subquery: subqueryInterface.(*sb.Builder),
					}).
					Select([]string{"*"})
			} else if tt.operator == "EXISTS" || tt.operator == "NOT EXISTS" {
				_, _, err = sb.NewBuilder(sb.DIALECT_MYSQL).
					Table("users").
					Where(&sb.Where{
						Operator: tt.operator,
						Subquery: subqueryInterface.(*sb.Builder),
						IsNot:    tt.operator == "NOT EXISTS",
					}).
					Select([]string{"*"})
			} else {
				// Comparison operators
				_, _, err = sb.NewBuilder(sb.DIALECT_MYSQL).
					Table("users").
					Where(&sb.Where{
						Column:   "id",
						Operator: tt.operator,
						Subquery: subqueryInterface.(*sb.Builder),
					}).
					Select([]string{"*"})
			}

			// Check error expectation
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error but none occurred")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestSubqueryValidationWithExists(t *testing.T) {
	// EXISTS should work with any number of columns
	subqueryBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("EXISTS should not panic for multiple columns, got: %v", r)
		}
	}()

	sql, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{
			Operator: "EXISTS",
			Subquery: subqueryBuilder.(*sb.Builder),
		}).
		Select([]string{"*"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	// Should generate valid SQL
	if sql == "" {
		t.Error("Expected SQL to be generated")
	}
}

func TestSubqueryValidationWithNilSubquery(t *testing.T) {
	// This should return an error during validation
	_, _, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{
			Operator: "IN",
			Subquery: nil,
		}).
		Select([]string{"*"})

	if err == nil {
		t.Error("Expected error for nil subquery")
	} else if err.Error() != "subquery cannot be nil" {
		t.Errorf("Expected error message 'subquery cannot be nil', got '%s'", err.Error())
	}
}

// TestWhereConditionWithEmptyValue ensures that WHERE conditions with empty values
// work correctly and don't incorrectly trigger subquery validation.
// This is a regression test for a bug where regular WHERE conditions with empty
// values failed with "subquery cannot be nil" error.
func TestWhereConditionWithEmptyValue(t *testing.T) {
	dialects := []string{sb.DIALECT_MYSQL, sb.DIALECT_POSTGRES, sb.DIALECT_SQLITE, sb.DIALECT_MSSQL}

	for _, dialect := range dialects {
		t.Run(dialect, func(t *testing.T) {
			// Create a WHERE condition with empty value that previously triggered the bug
			whereCondition := sb.Where{
				Column:   "name",
				Operator: "=",
				Value:    "", // Empty value should not trigger subquery validation
				Type:     "AND",
			}

			query, args, err := sb.NewBuilder(dialect).
				Table("users").
				Where(&whereCondition).
				Select([]string{"id", "name", "email"})

			// Should NOT fail with "subquery cannot be nil" error
			if err != nil {
				t.Errorf("WHERE condition with empty value failed for dialect %s: %v", dialect, err)
				return
			}

			if query == "" {
				t.Errorf("Empty query returned for dialect %s", dialect)
				return
			}

			if args == nil {
				t.Errorf("Nil args returned for dialect %s", dialect)
				return
			}

			// Verify the query contains the expected WHERE clause
			expectedPattern := "WHERE `name` = ?"
			switch dialect {
			case sb.DIALECT_POSTGRES:
				expectedPattern = "WHERE \"name\" = $1"
			case sb.DIALECT_SQLITE:
				expectedPattern = "WHERE \"name\" = ?"
			case sb.DIALECT_MSSQL:
				expectedPattern = "WHERE [name] = @p1"
			}

			if !strings.Contains(query, expectedPattern) {
				t.Errorf("Query doesn't contain expected WHERE pattern for dialect %s. Expected: %s, Got: %s", dialect, expectedPattern, query)
			}
		})
	}
}

// Enhanced Index Support Tests

func TestBuilderCreateIndexWithOptionsBasic(t *testing.T) {
	// Test basic index creation (should work like CreateIndex)
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "email"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_email` ON `users` (`email`);"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateUniqueIndex(t *testing.T) {
	// Test UNIQUE index creation
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		CreateUniqueIndex("idx_users_email_unique", "email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE UNIQUE INDEX "idx_users_email_unique" ON "users" ("email");`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateCompositeIndex(t *testing.T) {
	// Test composite index with ordering
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		CreateCompositeIndex("idx_orders_status_date", []sb.IndexColumn{
			{Name: "status", Direction: "ASC"},
			{Name: "created_at", Direction: "DESC"},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_orders_status_date` ON `orders` (`status`, `created_at` DESC);"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreatePartialIndex(t *testing.T) {
	// Test partial index (PostgreSQL)
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		CreatePartialIndex("idx_users_active", "deleted_at IS NULL", "email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE INDEX "idx_users_active" ON "users" ("email") WHERE deleted_at IS NULL;`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateCoveringIndex(t *testing.T) {
	// Test covering index with INCLUDE clause (PostgreSQL)
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		CreateCoveringIndex("idx_orders_customer", []string{"customer_name", "total"}, "customer_id")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE INDEX "idx_orders_customer" ON "orders" ("customer_id") INCLUDE ("customer_name", "total");`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateIndexWithOptionsIfNotExists(t *testing.T) {
	// Test IF NOT EXISTS support (PostgreSQL, SQLite)
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			IfNotExists: true,
			Columns:     []sb.IndexColumn{{Name: "email"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE INDEX IF NOT EXISTS "idx_users_email" ON "users" ("email");`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateIndexWithOptionsIfNotExistsMySQL(t *testing.T) {
	// Test that MySQL silently omits IF NOT EXISTS
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			IfNotExists: true,
			Columns:     []sb.IndexColumn{{Name: "email"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_email` ON `users` (`email`);"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateMySQLFullTextIndex(t *testing.T) {
	// Test MySQL FULLTEXT index
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("articles").
		CreateIndexWithOptions("idx_articles_content", sb.IndexOptions{
			Using:   sb.INDEX_TYPE_FULLTEXT,
			Columns: []sb.IndexColumn{{Name: "title"}, {Name: "content"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE FULLTEXT INDEX `idx_articles_content` ON `articles` (`title`, `content`);"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateMySQLSpatialIndex(t *testing.T) {
	// Test MySQL SPATIAL index
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("locations").
		CreateIndexWithOptions("idx_locations_coords", sb.IndexOptions{
			Using:   sb.INDEX_TYPE_SPATIAL,
			Columns: []sb.IndexColumn{{Name: "coordinates"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE SPATIAL INDEX `idx_locations_coords` ON `locations` (`coordinates`);"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreatePostgreSQLGINIndex(t *testing.T) {
	// Test PostgreSQL GIN index
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("documents").
		CreateIndexWithOptions("idx_documents_search", sb.IndexOptions{
			Using:   sb.INDEX_TYPE_GIN,
			Columns: []sb.IndexColumn{{Name: "search_vector"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE INDEX "idx_documents_search" ON "documents" USING GIN ("search_vector");`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateMySQLPrefixLength(t *testing.T) {
	// Test MySQL prefix length on text column
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_name", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "name", Length: 50}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_name` ON `users` (`name`(50));"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateMySQLUsingBTREE(t *testing.T) {
	// Test MySQL USING BTREE inside column list
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			Using:   sb.INDEX_TYPE_BTREE,
			Columns: []sb.IndexColumn{{Name: "email"}},
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_email` ON `users` (`email` USING BTREE);"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateIndexWithStorageParams(t *testing.T) {
	// Test PostgreSQL storage parameters
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "email"}},
			Storage: "fillfactor=90",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE INDEX "idx_users_email" ON "users" ("email") WITH (fillfactor=90);`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateIndexWithComment(t *testing.T) {
	// Test MySQL index comment
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "email"}},
			Comment: "Index for user email lookups",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_email` ON `users` (`email`) COMMENT 'Index for user email lookups';"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderCreateIndexWithEscapedComment(t *testing.T) {
	// Test MySQL index comment with single quote escaping
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_name", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "name"}},
			Comment: "User's name index",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE INDEX `idx_users_name` ON `users` (`name`) COMMENT 'User''s name index';"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderDropIndexWithOptionsBasic(t *testing.T) {
	// Test basic DROP INDEX (should work like DropIndex)
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderDropIndexWithOptionsIfExists(t *testing.T) {
	// Test DROP INDEX IF EXISTS (PostgreSQL)
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
			IfExists: true,
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderDropIndexWithOptionsWithSchema(t *testing.T) {
	// Test DROP INDEX with schema (PostgreSQL)
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
			Schema: "public",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX "public"."idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderDropIndexWithOptionsIfExistsAndSchema(t *testing.T) {
	// Test DROP INDEX IF EXISTS with schema (PostgreSQL)
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
			IfExists: true,
			Schema:   "public",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX IF EXISTS "public"."idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderDropIndexWithOptionsMySQLNoIfExists(t *testing.T) {
	// Test that MySQL silently omits IF EXISTS
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexWithOptions("idx_users_email", sb.DropIndexOptions{
			IfExists: true,
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

// Error handling tests

func TestBuilderCreateIndexWithOptionsEmptyName(t *testing.T) {
	// Test error handling for empty index name
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "email"}},
		})

	if err == nil {
		t.Fatal("Expected error for empty index name")
	}
	if err.Error() != "ValidationError: index name cannot be empty" {
		t.Fatalf("Expected 'ValidationError: index name cannot be empty' but got: %v", err)
	}
}

func TestBuilderCreateIndexWithOptionsMissingTable(t *testing.T) {
	// Test error handling for missing table
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			Columns: []sb.IndexColumn{{Name: "email"}},
		})

	if err == nil {
		t.Fatal("Expected error for missing table")
	}
	if err.Error() != "ValidationError: no table specified" {
		t.Fatalf("Expected 'ValidationError: no table specified' but got: %v", err)
	}
}

func TestBuilderCreateIndexWithOptionsEmptyColumns(t *testing.T) {
	// Test error handling for empty columns
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		CreateIndexWithOptions("idx_users_email", sb.IndexOptions{
			Columns: []sb.IndexColumn{},
		})

	if err == nil {
		t.Fatal("Expected error for empty columns")
	}
	if err.Error() != "ValidationError: columns cannot be empty" {
		t.Fatalf("Expected 'ValidationError: columns cannot be empty' but got: %v", err)
	}
}

func TestBuilderDropIndexWithOptionsEmptyName(t *testing.T) {
	// Test error handling for empty index name
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexWithOptions("", sb.DropIndexOptions{})

	if err == nil {
		t.Fatal("Expected error for empty index name")
	}
	if err.Error() != "ValidationError: index name cannot be empty" {
		t.Fatalf("Expected 'ValidationError: index name cannot be empty' but got: %v", err)
	}
}

// Dialect-specific comprehensive tests

func TestBuilderEnhancedIndexSupportMySQL(t *testing.T) {
	// Test all MySQL-specific features work together
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("products").
		CreateIndexWithOptions("idx_products_search", sb.IndexOptions{
			Unique:  true,
			Using:   sb.INDEX_TYPE_FULLTEXT,
			Columns: []sb.IndexColumn{{Name: "name", Length: 100}, {Name: "description", Length: 255}},
			Comment: "Full-text search index for products",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "CREATE UNIQUE FULLTEXT INDEX `idx_products_search` ON `products` (`name`(100), `description`(255)) COMMENT 'Full-text search index for products';"
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderEnhancedIndexSupportPostgreSQL(t *testing.T) {
	// Test all PostgreSQL-specific features work together
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		CreateIndexWithOptions("idx_orders_customer_active", sb.IndexOptions{
			Unique:      true,
			IfNotExists: true,
			Using:       sb.INDEX_TYPE_GIN,
			Columns:     []sb.IndexColumn{{Name: "customer_id"}},
			Include:     []string{"order_date", "total"},
			Where:       "status = 'active'",
			Storage:     "fillfactor=70",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE UNIQUE INDEX IF NOT EXISTS "idx_orders_customer_active" ON "orders" USING GIN ("customer_id") INCLUDE ("order_date", "total") WHERE status = 'active' WITH (fillfactor=70);`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderEnhancedIndexSupportSQLite(t *testing.T) {
	// Test SQLite features (limited set)
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		CreateIndexWithOptions("idx_users_active_email", sb.IndexOptions{
			Unique:      true,
			IfNotExists: true,
			Columns:     []sb.IndexColumn{{Name: "email"}},
			Where:       "deleted_at IS NULL",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE UNIQUE INDEX IF NOT EXISTS "idx_users_active_email" ON "users" ("email") WHERE deleted_at IS NULL;`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

func TestBuilderEnhancedIndexSupportMSSQL(t *testing.T) {
	// Test MSSQL features
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("orders").
		CreateIndexWithOptions("idx_orders_customer", sb.IndexOptions{
			Unique:  true,
			Columns: []sb.IndexColumn{{Name: "customer_id", Direction: "DESC"}},
			Include: []string{"order_date", "total"},
			Where:   "status = 'active'",
		})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `CREATE UNIQUE INDEX [idx_orders_customer] ON [orders] ([customer_id] DESC) INCLUDE ([order_date], [total]) WHERE status = 'active';`
	if sql != expected {
		t.Fatalf("Expected:\n%s\nbut found:\n%s", expected, sql)
	}
}

// TDD Test: WHERE clause SQL generation fixes
func TestWhereClauseSQLOperators(t *testing.T) {
	// Test that WHERE conditions generate proper SQL with operators
	tests := []struct {
		name     string
		dialect  string
		where    *sb.Where
		expected string
	}{
		{
			name:     "MySQL email equality",
			dialect:  sb.DIALECT_MYSQL,
			where:    &sb.Where{Column: "email", Operator: "=", Value: "test@example.com"},
			expected: "WHERE `email` = ?;",
		},
		{
			name:     "PostgreSQL email equality",
			dialect:  sb.DIALECT_POSTGRES,
			where:    &sb.Where{Column: "email", Operator: "=", Value: "test@example.com"},
			expected: "WHERE \"email\" = $1;",
		},
		{
			name:     "MySQL IS NULL",
			dialect:  sb.DIALECT_MYSQL,
			where:    &sb.Where{Column: "soft_deleted_at", Operator: "IS NULL"},
			expected: "WHERE `soft_deleted_at` IS NULL;",
		},
		{
			name:     "PostgreSQL IS NULL",
			dialect:  sb.DIALECT_POSTGRES,
			where:    &sb.Where{Column: "soft_deleted_at", Operator: "IS NULL"},
			expected: "WHERE \"soft_deleted_at\" IS NULL;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := sb.NewBuilder(tt.dialect).Table("users").Where(tt.where)
			sql, _, err := builder.Select([]string{"id", "email", "name"})

			if err != nil {
				t.Fatal("Unexpected error:", err)
			}

			// Extract WHERE clause from the generated SQL
			whereIndex := strings.Index(sql, "WHERE ")
			if whereIndex == -1 {
				t.Fatalf("No WHERE clause found in SQL: %s", sql)
			}

			whereClause := sql[whereIndex:]
			if whereClause != tt.expected {
				t.Fatalf("Expected WHERE clause:\n%s\nbut found:\n%s", tt.expected, whereClause)
			}
		})
	}
}
