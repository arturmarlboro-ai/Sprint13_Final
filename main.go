package main

import (
	"Sprint13_Final/pkg/db"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var install bool

func main() {
	err := db.InitDB()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}

	fmt.Println("Database initialized successfully.")

}
