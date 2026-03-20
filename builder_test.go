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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
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
	selectSQL, err := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})
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
	selectSQL, err := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})
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
	selectSQL, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})
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
	selectSQL, err := sb.NewBuilder(sb.DIALECT_MYSQL).Table("users").Select([]string{"FirstName", "LastName"})
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
	selectSQL, err := sb.NewBuilder(sb.DIALECT_POSTGRES).Table("users").Select([]string{"FirstName", "LastName"})
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
	selectSQL, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"FirstName", "LastName"})
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

	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "email", "name", "created_at", "deleted_at"})

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

	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).Table("users").Select([]string{"id", "name", "created_at", "deleted_at"})

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
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DELETE FROM `users` WHERE `FirstName` = \"Tom\" OR `FirstName` = \"Sam\" LIMIT 12 OFFSET 34;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableDeleteSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DELETE FROM "users" WHERE "FirstName" = 'Tom' OR "FirstName" = 'Sam';`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectMysql(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
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
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "first_name", Operator: "!=", Value: "Jane"}).
		OrderBy("first_name", "asc").
		Limit(10).
		Offset(20).
		GroupBy(sb.GroupBy{Column: "passport"}).
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
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Limit(1).
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
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Limit(1).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Limit(1).
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
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "UPDATE `users` SET `first_name`=\"Tom\", `last_name`=\"Jones\" WHERE `id` = \"1\" LIMIT 1;"
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdatePostgres(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `UPDATE "users" SET "first_name"='Tom', "last_name"='Jones' WHERE "id" = '1' LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableUpdateSqlite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `UPDATE "users" SET "first_name"='Tom', "last_name"='Jones' WHERE "id" = '1' LIMIT 1;`
	if sql != expected {
		t.Fatal("Expected:\n", expected, "\nbut found:\n", sql)
	}
}

func TestBuilderTableSelectMysqlInj(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
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
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58\" OR 1 = 1;--"}).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		Where(&sb.Where{Column: "id", Operator: "=", Value: "58' OR 1 = 1;--"}).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
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

	// Test with ResetIdentity
	sql, err = sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		TruncateWithOptions(sb.TruncateOptions{ResetIdentity: true})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected = "TRUNCATE TABLE [users]; DBCC CHECKIDENT ('users', RESEED, 0)"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderTruncateWithComplexTableName(t *testing.T) {
	// Test with schema.table format
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("public.users").
		Truncate()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `TRUNCATE TABLE "public"."users";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test with database.schema.table format
	sql2, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("dbo.users").
		Truncate()

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected = "TRUNCATE TABLE [dbo].[users];"
	if sql2 != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql2)
	}
}

func TestBuilderTruncateErrorHandling(t *testing.T) {
	// Test error when no table is specified
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).Truncate()

	if err == nil {
		t.Fatalf("Expected error when no table specified")
	}

	expectedMsg := "ValidationError: no table specified"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderTruncateWithOptionsErrorHandling(t *testing.T) {
	// Test error when no table is specified
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).TruncateWithOptions(sb.TruncateOptions{})

	if err == nil {
		t.Fatalf("Expected error when no table specified")
	}

	expectedMsg := "ValidationError: no table specified"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderTruncateUnsupportedDialect(t *testing.T) {
	// Test error for unsupported dialect
	builder := sb.NewBuilder("unknown")
	_, err := builder.Table("users").Truncate()

	if err == nil {
		t.Fatalf("Expected error for unsupported dialect")
	}

	expectedMsg := "ValidationError: unsupported dialect: unknown"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

// Test Builder DropIndex functionality

func TestBuilderDropIndexMySQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndex("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexPostgreSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndex("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexSQLite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		DropIndex("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexMSSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		DropIndex("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX [idx_users_email] ON [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsMySQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexIfExists("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsPostgreSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexIfExists("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsSQLite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		DropIndexIfExists("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexIfExistsMSSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		DropIndexIfExists("idx_users_email")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX IF EXISTS [idx_users_email] ON [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexWithSchemaPostgreSQL(t *testing.T) {
	// Test with schema
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `DROP INDEX IF EXISTS "public"."idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test without schema (should use regular behavior)
	sql, err = sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithSchema("idx_users_email", "")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected = `DROP INDEX IF EXISTS "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// Test with empty schema and index name
	_, err = sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		DropIndexWithSchema("", "")

	if err == nil {
		t.Fatal("Expected error for empty index name")
	}

	expectedMsg := "ValidationError: index name cannot be empty"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderDropIndexWithSchemaOtherDialects(t *testing.T) {
	// MySQL should fall back to regular DropIndex
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "DROP INDEX `idx_users_email` ON `users`;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// SQLite should fall back to regular DropIndex
	sql, err = sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected = `DROP INDEX "idx_users_email";`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}

	// MSSQL should fall back to regular DropIndex
	sql, err = sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		DropIndexWithSchema("idx_users_email", "public")

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected = "DROP INDEX [idx_users_email] ON [users];"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderDropIndexErrorHandling(t *testing.T) {
	// Test error when no table is specified
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).DropIndex("idx_users_email")

	if err == nil {
		t.Fatalf("Expected error when no table specified")
	}

	expectedMsg := "ValidationError: no table specified"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderDropIndexEmptyIndexName(t *testing.T) {
	// Test error when index name is empty
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		DropIndex("")

	if err == nil {
		t.Fatalf("Expected error when index name is empty")
	}

	expectedMsg := "ValidationError: index name cannot be empty"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderDropIndexIfExistsErrorHandling(t *testing.T) {
	// Test error when no table is specified
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).DropIndexIfExists("idx_users_email")

	if err == nil {
		t.Fatalf("Expected error when no table specified")
	}

	expectedMsg := "ValidationError: no table specified"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderDropIndexWithSchemaErrorHandling(t *testing.T) {
	// Test error when no table is specified
	_, err := sb.NewBuilder(sb.DIALECT_POSTGRES).DropIndexWithSchema("idx_users_email", "public")

	if err == nil {
		t.Fatalf("Expected error when no table specified")
	}

	expectedMsg := "ValidationError: no table specified"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderDropIndexUnsupportedDialect(t *testing.T) {
	// Test error for unsupported dialect
	builder := sb.NewBuilder("unknown")
	_, err := builder.Table("users").DropIndex("idx_users_email")

	if err == nil {
		t.Fatalf("Expected error for unsupported dialect")
	}

	expectedMsg := "ValidationError: unsupported dialect: unknown"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

// Test Builder JOIN functionality

func TestBuilderJoinInnerMySQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `orders`.*, `users`.`name` FROM `orders`  INNER JOIN `users` ON orders.user_id = users.id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinLeftPostgreSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		LeftJoin("profiles", "orders.user_id = profiles.user_id").
		Select([]string{"orders.*", "profiles.avatar"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "orders".*, "profiles"."avatar" FROM "orders"  LEFT JOIN "profiles" ON orders.user_id = profiles.user_id;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinRightMSSQL(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("orders").
		RightJoin("users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT [orders].*, [users].[name] FROM [orders]  RIGHT JOIN [users] ON orders.user_id = users.id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinWithAliasSQLite(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("orders").
		JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "orders.user_id = p.user_id").
		Select([]string{"orders.*", "p.avatar"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "orders".*, "p"."avatar" FROM "orders"  LEFT JOIN "profiles" AS "p" ON orders.user_id = p.user_id;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinMultiple(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		LeftJoin("profiles", "users.id = profiles.user_id").
		Select([]string{"orders.total", "users.name", "profiles.avatar"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `orders`.`total`, `users`.`name`, `profiles`.`avatar` FROM `orders`  INNER JOIN `users` ON orders.user_id = users.id  LEFT JOIN `profiles` ON users.id = profiles.user_id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinFullOuter(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Join(sb.JOIN_FULL, "users", "orders.user_id = users.id").
		Select([]string{"orders.*", "users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "orders".*, "users"."name" FROM "orders"  FULL JOIN "users" ON orders.user_id = users.id;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinCross(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("orders").
		Join(sb.JOIN_CROSS, "users", "1=1").
		Select([]string{"orders.*", "users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "orders".*, "users"."name" FROM "orders"  CROSS JOIN "users" ON 1=1;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinWithWhere(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "active"}).
		Select([]string{"orders.*", "users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `orders`.*, `users`.`name` FROM `orders`  INNER JOIN `users` ON orders.user_id = users.id WHERE `orders`.`status` = \"active\";"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinWithOrderBy(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		LeftJoin("users", "orders.user_id = users.id").
		OrderBy("orders.created_at", "DESC").
		Select([]string{"orders.*", "users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := `SELECT "orders".*, "users"."name" FROM "orders"  LEFT JOIN "users" ON orders.user_id = users.id ORDER BY "orders"."created_at" DESC;`
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}
func TestBuilderJoinErrorHandlingEmptyCondition(t *testing.T) {
	// Test error when ON condition is empty
	builder := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders")

	builder.InnerJoin("users", "")

	// Error should be returned when trying to generate SQL
	_, err := builder.Select([]string{"*"})

	if err == nil {
		t.Fatalf("Expected error when ON condition is empty")
	}

	expectedMsg := "ValidationError: ON condition cannot be empty"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderJoinWithAliasErrorHandlingEmptyCondition(t *testing.T) {
	// Test error when ON condition is empty
	builder := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders")

	builder.JoinWithAlias(sb.JOIN_LEFT, "profiles", "p", "")

	// Error should be returned when trying to generate SQL
	_, err := builder.Select([]string{"*"})

	if err == nil {
		t.Fatalf("Expected error when ON condition is empty")
	}

	expectedMsg := "ValidationError: ON condition cannot be empty"
	if err.Error() != expectedMsg {
		t.Fatalf("Expected error message: %s but got: %v", expectedMsg, err)
	}
}

func TestBuilderJoinComplexTableName(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("public.users", "orders.user_id = public.users.id").
		Select([]string{"orders.*", "public.users.name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `orders`.*, `public`.`users`.`name` FROM `orders`  INNER JOIN `public`.`users` ON orders.user_id = public.users.id;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinUpdate(t *testing.T) {
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		InnerJoin("users", "orders.user_id = users.id").
		Where(&sb.Where{Column: "orders.status", Operator: "=", Value: "pending"}).
		Update(map[string]string{"orders.status": "processed"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "UPDATE `orders` SET `orders`.`status`=\"processed\"  INNER JOIN `users` ON orders.user_id = users.id WHERE `orders`.`status` = \"pending\";"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderJoinNoJoins(t *testing.T) {
	// Test that queries work normally without joins
	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		Select([]string{"*"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT * FROM `orders` WHERE `status` = \"active\";"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

// Test Builder Subquery functionality

func TestBuilderSubqueryInMySQL(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > \"1000\");"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryInPostgreSQL(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT \"name\" FROM \"users\" WHERE \"id\" IN (SELECT * FROM \"orders\" WHERE \"total\" > '1000');"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryInSQLite(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_SQLITE).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT \"name\" FROM \"users\" WHERE \"id\" IN (SELECT * FROM \"orders\" WHERE \"total\" > '1000');"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryInMSSQL(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MSSQL).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT [name] FROM [users] WHERE [id] IN (SELECT * FROM [orders] WHERE [total] > 1000);"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryNotInMySQL(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "inactive"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		NotInSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `id` NOT IN (SELECT * FROM `orders` WHERE `status` = \"inactive\");"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryExistsMySQL(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"}).
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Exists(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name", "email"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name`, `email` FROM `users` WHERE EXISTS (SELECT * FROM `orders` WHERE `user_id` = \"users.id\" AND `status` = \"active\");"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryNotExistsPostgreSQL(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("orders").
		Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"}).
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"})

	builder, err := sb.NewBuilder(sb.DIALECT_POSTGRES).
		Table("users").
		NotExists(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name", "email"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT \"name\", \"email\" FROM \"users\" WHERE NOT EXISTS (SELECT * FROM \"orders\" WHERE \"user_id\" = 'users.id' AND \"status\" = 'active');"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryWithWhereClause(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `status` = \"active\" AND `id` IN (SELECT * FROM `orders` WHERE `total` > \"1000\");"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryWithOrderBy(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	builder.OrderBy("name", "ASC")

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > \"1000\") ORDER BY `name` ASC;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryWithLimit(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	builder.Limit(10)

	sql, err := builder.Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `id` IN (SELECT * FROM `orders` WHERE `total` > \"1000\") LIMIT 10;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryComparison(t *testing.T) {
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "user_id", Operator: "=", Value: "users.id"})

	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Where(&sb.Where{
			Column:   "id",
			Operator: ">",
			Subquery: subquery.(*sb.Builder),
		}).
		Select([]string{"name"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT `name` FROM `users` WHERE `id` > (SELECT * FROM `orders` WHERE `user_id` = \"users.id\");"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryErrorHandlingNilSubquery(t *testing.T) {
	// Test error when subquery is nil
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		Exists(nil)

	if err == nil {
		t.Fatalf("Expected error but none occurred")
	}
	if err.Error() != "ArgumentError: subquery cannot be nil" {
		t.Fatalf("Expected 'ArgumentError: subquery cannot be nil' but got: %v", err)
	}
}

func TestBuilderSubqueryErrorHandlingNilSubqueryNotExists(t *testing.T) {
	// Test error when subquery is nil for NotExists
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		NotExists(nil)

	if err == nil {
		t.Fatalf("Expected error but none occurred")
	}
	if err.Error() != "ArgumentError: subquery cannot be nil" {
		t.Fatalf("Expected 'ArgumentError: subquery cannot be nil' but got: %v", err)
	}
}

func TestBuilderSubqueryErrorHandlingNilSubqueryIn(t *testing.T) {
	// Test error when subquery is nil for InSubquery
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(nil)

	if err == nil {
		t.Fatalf("Expected error but none occurred")
	}
	if err.Error() != "ArgumentError: subquery cannot be nil" {
		t.Fatalf("Expected 'ArgumentError: subquery cannot be nil' but got: %v", err)
	}
}

func TestBuilderSubqueryErrorHandlingNilSubqueryNotIn(t *testing.T) {
	// Test error when subquery is nil for NotInSubquery
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		NotInSubquery(nil)

	if err == nil {
		t.Fatalf("Expected error but none occurred")
	}
	if err.Error() != "ArgumentError: subquery cannot be nil" {
		t.Fatalf("Expected 'ArgumentError: subquery cannot be nil' but got: %v", err)
	}
}
func TestBuilderSubqueryComplex(t *testing.T) {
	// Complex subquery with multiple conditions
	subquery := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("order_items").
		Where(&sb.Where{Column: "quantity", Operator: ">", Value: "5"}).
		Where(&sb.Where{Column: "price", Operator: ">", Value: "100"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("orders").
		Where(&sb.Where{Column: "status", Operator: "=", Value: "active"}).
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	builder.OrderBy("created_at", "DESC").
		Limit(20)

	sql, err := builder.Select([]string{"*"})

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	expected := "SELECT * FROM `orders` WHERE `status` = \"active\" AND `id` IN (SELECT * FROM `order_items` WHERE `quantity` > \"5\" AND `price` > \"100\") ORDER BY `created_at` DESC LIMIT 20;"
	if sql != expected {
		t.Fatalf("Expected: %s but found: %s", expected, sql)
	}
}

func TestBuilderSubqueryCreateMethod(t *testing.T) {
	// Test the Subquery() method
	subqueryBuilder := sb.NewBuilder(sb.DIALECT_MYSQL).Subquery()

	// Test that the subquery builder works normally
	subquery := subqueryBuilder.
		Table("orders").
		Where(&sb.Where{Column: "total", Operator: ">", Value: "1000"})

	builder, err := sb.NewBuilder(sb.DIALECT_MYSQL).
		Table("users").
		InSubquery(subquery)

	if err != nil {
		t.Fatal("Unexpected error:", err)
	}

	sql, err := builder.Select([]string{"name"})

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
			_, err := subqueryInterface.Select(tt.columns)
			if err != nil {
				t.Fatal("Unexpected error setting subquery columns:", err)
			}

			// Try to use the subquery
			if tt.operator == "IN" || tt.operator == "NOT IN" {
				_, err = sb.NewBuilder(sb.DIALECT_MYSQL).
					Table("users").
					Where(&sb.Where{
						Operator: tt.operator,
						Subquery: subqueryInterface.(*sb.Builder),
					}).
					Select([]string{"*"})
			} else if tt.operator == "EXISTS" || tt.operator == "NOT EXISTS" {
				_, err = sb.NewBuilder(sb.DIALECT_MYSQL).
					Table("users").
					Where(&sb.Where{
						Operator: tt.operator,
						Subquery: subqueryInterface.(*sb.Builder),
						IsNot:    tt.operator == "NOT EXISTS",
					}).
					Select([]string{"*"})
			} else {
				// Comparison operators
				_, err = sb.NewBuilder(sb.DIALECT_MYSQL).
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

	sql, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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
	_, err := sb.NewBuilder(sb.DIALECT_MYSQL).
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
