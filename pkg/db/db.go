package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"
)

var install bool

func InitDB() error {
	dbFile := "scheduler.db"
	_, err := os.Stat(dbFile)

	if err != nil {
		install = true
	}

	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	db.SetMaxIdleConns(2)
	db.SetMaxOpenConns(5)
	db.SetConnMaxIdleTime(time.Minute * 5)
	db.SetConnMaxLifetime(time.Hour)
	t := time.Now().Format("20060102")
	fmt.Println("Database configured.")
	log.Println("Database configured" + t)

	return nil
}
