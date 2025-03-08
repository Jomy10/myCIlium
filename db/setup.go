package db

import (
	"context"
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	sqldblogger "github.com/simukti/sqldb-logger"
)

var db *sql.DB

type DatabaseLogger struct{}

func (logger *DatabaseLogger) Log(ctx context.Context, level sqldblogger.Level, msg string, data map[string]interface{}) {
	log.Printf("[%s:%s] %s", ctx, level, msg)
}

func SetupDatabase() error {
	dsn := "./data.db"
	_db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}
	db = _db

	// Enable logging
	loggerAdapter := DatabaseLogger{}
	db = sqldblogger.OpenDriver(dsn, db.Driver(), &loggerAdapter,
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
