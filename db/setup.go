package db

import (
	"database/sql"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/mattn/go-sqlite3"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/logrusadapter"
)

var db *sql.DB

func SetupDatabase() error {
	dsn := "./data.db"
	_db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	db = _db

	// Enable logging
	logrus := log.New()
	logrus.Level = log.TraceLevel
	logrus.Formatter = &log.TextFormatter{}
	db = sqldblogger.OpenDriver(dsn, db.Driver(), logrusadapter.New(logrus),
		sqldblogger.WithSQLQueryAsMessage(true),
		sqldblogger.WithMinimumLevel(sqldblogger.LevelTrace),
		sqldblogger.WithPreparerLevel(sqldblogger.LevelTrace),
		sqldblogger.WithQueryerLevel(sqldblogger.LevelTrace),
		sqldblogger.WithExecerLevel(sqldblogger.LevelTrace),
	)
	err = db.Ping()
	if err != nil {
		CloseDatabase()
		return err
	}

	sql := `
	CREATE TABLE IF NOT EXISTS Meta (
		dbVersion integer NOT NULL
	);

	CREATE TABLE IF NOT EXISTS BuildRequest (
		id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
		repo string NOT NULL,
		revision string NULL,
		platform string NOT NULL,
		status integer NOT NULL,
		requested string NOT NULL DEFAULT CURRENT_TIMESTAMP,
		requestedBy string NOT NULL,
		statusDate string NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS par_Status (
		id int NOT NULL,
		name string NOT NULL
	);
	`

	_, err = db.Exec(sql)
	if err != nil {
		CloseDatabase()
		return err
	}

	rows, err := db.Query("SELECT dbVersion FROM Meta")
	if err != nil {
		CloseDatabase()
		return err
	}
	defer rows.Close()
	var dbVersion int = 0
	if rows.Next() {
		err = rows.Scan(&dbVersion)
		if err != nil {
			CloseDatabase()
			return err
		}
	}

	if dbVersion <= 0 {
		sql = `
		INSERT INTO Meta (dbVersion) VALUES (1);

		INSERT INTO par_Status (id, name)
		VALUES
			(1, 'Requested'),
			(2, 'Started'),
			(3, 'Finished');
		`

		_, err = db.Exec(sql)
		if err != nil {
			CloseDatabase()
			return err
		}
	}

	return nil
}

func CloseDatabase() {
	db.Close()
	db = nil
}

func ResetDatabase() {
	log.Fatal(os.Remove("data.db"))
}
