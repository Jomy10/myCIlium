package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func SetupDatabase() error {
	_db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		return err
	}
	db = _db

	sql := `
	CREATE TABLE IF NOT EXISTS Meta (
		dbVersion int
	);

	CREATE TABLE IF NOT EXISTS BuildRequest (
		id int,
		repo string,
		revision string,
		platform string,
		status int,
		requested string DEFAULT CURRENT_TIMESTAMP,
		requestedBy string,
		statusDate string DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS par_Status (
		id int,
		name string
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
