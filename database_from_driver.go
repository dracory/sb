package sb

import (
	"database/sql"
	"errors"
	"time"
)

func NewDatabaseFromDriver(driverName, dataSourceName string) (DatabaseInterface, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, errors.New("failed to open DB: " + err.Error())
	}

	return &Database{
		db:             db,
		databaseType:   DatabaseDriverName(db),
		debug:          false,
		sqlLogEnabled:  false,
		sqlLog:         map[string]string{},
		sqlDurationLog: map[string]time.Duration{},
	}, nil
}
