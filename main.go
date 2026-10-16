package main

import (
	"fmt"
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

	data := os.Getenv("TODO_DBFILE")
	if data == "" {
		data = "scheduler.db"
	}

	if err := db.Init(data); err != nil {
		log.Fatalf("DB init error: %v", err)
	}
	defer db.DB.Close()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	host := fmt.Sprintf(":%s", port)

	if err := server.Run(host); err != nil {
		log.Fatalf("Server error: %v", err)
	}

}
