package main

import (
	"Sprint13_Final/pkg/db"
	"Sprint13_Final/pkg/server"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

var install bool

func main() {
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logfile.Close()

	log.SetOutput(logfile)

	err = db.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
		os.Exit(1)
	}
	log.Println("Database initialized successfully.")

	server.Run("localhost:7540")

}
