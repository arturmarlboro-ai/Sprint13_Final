package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init(dbFile string) error {
	log.Println("Init db")
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	_, err = DB.Exec(`
			CREATE TABLE IF NOT EXISTS scheduler (
				id      INTEGER PRIMARY KEY AUTOINCREMENT,
				date    TEXT NOT NULL,
				title   TEXT NOT NULL,
				comment TEXT,
				repeat  TEXT
			);
			CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
		`)
	return err
}
