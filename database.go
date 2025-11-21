package sb

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/dracory/uid"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/spf13/cast"
)

// == CONSTRUCTOR ============================================================

func NewDatabase(db *sql.DB, databaseType string) DatabaseInterface {
	return &Database{
		db:           db,
		databaseType: databaseType,
	}
}

// == TYPE ===================================================================

type Database struct {
	db             *sql.DB
	tx             *sql.Tx
	databaseType   string
	sqlLogEnabled  bool
	sqlLog         map[string]string
	sqlDurationLog map[string]time.Duration
	debug          bool
}

// == INTERFACE ==============================================================

var _ DatabaseInterface = (*Database)(nil)

// == PUBLIC METHODS =========================================================

func (d *Database) IsMssql() bool {
	return d.databaseType == DIALECT_MSSQL
}

func (d *Database) IsMysql() bool {
	return d.databaseType == DIALECT_MYSQL
}

func (d *Database) IsPostgres() bool {
	return d.databaseType == DIALECT_POSTGRES
}

func (d *Database) IsSqlite() bool {
	return d.databaseType == DIALECT_SQLITE
}

func (d *Database) SqlLog() []map[string]string {
	logArray := []map[string]string{}
	for key, value := range d.sqlLog {
		sqlDuration := d.sqlDurationLog[key]
		logArray = append(logArray, map[string]string{
			"sql":  value,
			"time": sqlDuration.String(),
		})
	}
	return logArray
}

func (d *Database) SqlLogEmpty() {
	d.sqlLog = map[string]string{}
	d.sqlDurationLog = map[string]time.Duration{}
}

func (d *Database) SqlLogLen() int {
	return len(d.sqlLog)
}

func (d *Database) SqlLogEnable(enable bool) {
	d.sqlLogEnabled = enable
}

func (d *Database) SqlLogShrink(leaveLast int) {
	if len(d.sqlLog) <= leaveLast {
		return
	}

	keys := []string{}

	for key := range d.sqlLog {
		keys = append(keys, key)
	}

	tempSqlLog := map[string]string{}
	tempSqlDurationLog := map[string]time.Duration{}
	lastKeys := keys[leaveLast:]
	for _, key := range lastKeys {
		tempSqlLog[key] = d.sqlLog[key]
		tempSqlDurationLog[key] = d.sqlDurationLog[key]
	}

	d.sqlLog = tempSqlLog
	d.sqlDurationLog = tempSqlDurationLog
}

func (d *Database) DebugEnable(debug bool) {
	d.debug = debug
}

func (d *Database) Type() string {
	return d.databaseType
}

func (d *Database) Open() (err error) {
	d.db, err = sql.Open("sqlite3", "test.db")
	return err
}

func (d *Database) Close() (err error) {
	return d.db.Close()
}

func (d *Database) DB() *sql.DB {
	return d.db
}

func (d *Database) BeginTransaction() (err error) {
	if d.tx != nil {
		return errors.New("transaction already in progress")
	}

	tx, err := d.db.Begin()
	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}
	d.tx = tx

	return err
}

func (d *Database) BeginTransactionWithContext(ctx context.Context, opts *sql.TxOptions) (err error) {
	if d.tx != nil {
		return errors.New("transaction already in progress")
	}

	tx, err := d.db.BeginTx(ctx, opts)

	if err != nil {
		return errors.New("failed to begin transaction: " + err.Error())
	}

	d.tx = tx

	return nil
}

func (d *Database) ExecInTransaction(ctx context.Context, fn func(d *Database) error) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	err = d.BeginTransactionWithContext(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err := d.RollbackTransaction()
			if err != nil {
				log.Println("sqldb rollback error: " + err.Error())
			}
		}
	}()

	err = fn(&Database{db: d.db, tx: d.tx})

	if err == nil {
		err = d.CommitTransaction()
	}

	return
}

func (d *Database) Exec(ctx context.Context, sqlStr string, args ...any) (sql.Result, error) {
	if d.sqlLogEnabled {
		if d.sqlLog == nil {
			d.sqlLog = map[string]string{}
			d.sqlDurationLog = map[string]time.Duration{}
		}

		sqlID := uid.HumanUid()

		d.sqlLog[sqlID] = sqlStr

		start := time.Now()
		defer func() {
			d.sqlDurationLog[sqlID] = time.Since(start)
		}()
	}

	if d.debug {
		log.Println(sqlStr)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if d.tx != nil {
		return d.tx.ExecContext(ctx, sqlStr, args...)
	}
	return d.db.ExecContext(ctx, sqlStr, args...)
}

func (d *Database) Query(ctx context.Context, sqlStr string, args ...any) (*sql.Rows, error) {
	if d.sqlLogEnabled {
		if d.sqlLog == nil {
			d.sqlLog = map[string]string{}
			d.sqlDurationLog = map[string]time.Duration{}
		}

		sqlID := uid.HumanUid()

		d.sqlLog[sqlID] = sqlStr

		start := time.Now()
		defer func() {
			d.sqlDurationLog[sqlID] = time.Since(start)
		}()
	}

	if d.debug {
		log.Println(sqlStr)
	}

	if ctx == nil {
		ctx = context.Background()
	}

	if d.tx != nil {
		return d.tx.QueryContext(ctx, sqlStr, args...)
	}
	return d.db.QueryContext(ctx, sqlStr, args...)
}

func (d *Database) CommitTransaction() (err error) {
	if d.tx == nil {
		return errors.New("no transaction in progress")
	}

	err = d.tx.Commit()

	if err != nil {
		return errors.New("failed to commit transaction: " + err.Error())
	}

	d.tx = nil // empty transaction

	return err
}

func (d *Database) RollbackTransaction() (err error) {
	if d.tx == nil {
		return errors.New("no transaction in progress")
	}

	err = d.tx.Rollback()

	if err != nil {
		return errors.New("failed to rollback transaction: " + err.Error())
	}

	d.tx = nil // empty transaction

	return err
}

func (d *Database) SelectToMapAny(ctx context.Context, sqlStr string, args ...any) ([]map[string]any, error) {
	if d.sqlLogEnabled {
		if d.sqlLog == nil {
			d.sqlLog = map[string]string{}
			d.sqlDurationLog = map[string]time.Duration{}
		}

		sqlID := uid.HumanUid()

		d.sqlLog[sqlID] = sqlStr

		start := time.Now()
		defer func() {
			d.sqlDurationLog[sqlID] = time.Since(start)
		}()
	}

	if d.debug {
		log.Println(sqlStr)
	}

	listMap := []map[string]any{}

	if ctx == nil {
		ctx = context.Background()
	}

	err := sqlscan.Select(ctx, d.db, &listMap, sqlStr, args...)
	if err != nil {
		if sqlscan.NotFound(err) {
			return []map[string]any{}, nil
		}

		return []map[string]any{}, err
	}

	return listMap, nil
}

func (d *Database) SelectToMapString(ctx context.Context, sqlStr string, args ...any) ([]map[string]string, error) {
	listMapAny, err := d.SelectToMapAny(ctx, sqlStr, args...)

	if err != nil {
		return []map[string]string{}, err
	}

	listMapString := []map[string]string{}

	for i := range listMapAny {
		mapString := cast.ToStringMapString(listMapAny[i])
		listMapString = append(listMapString, mapString)
	}

	return listMapString, nil
}

func (d *Database) Tx() *sql.Tx {
	return d.tx
}
