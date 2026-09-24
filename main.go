package main

import (
	"log"
	"os"

	"Sprint13_Final/pkg/db"
	"Sprint13_Final/pkg/server"
)

func main() {
	logfile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer logfile.Close()
	log.SetOutput(logfile)

	if err := db.Init("scheduler.db"); err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer db.DB.Close()

	if err := server.Run("localhost:7540"); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
