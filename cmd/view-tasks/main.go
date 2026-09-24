package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func main() {
	db, err := sql.Open("sqlite", "scheduler.db")
	if err != nil {
		log.Fatalf("can't open db: %v", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
	if err != nil {
		// Если таблица ещё не создана — покажем понятную ошибку
		if err.Error() == "no such table: scheduler" {
			fmt.Println("Таблица scheduler ещё не создана. Запусти сервер один раз, чтобы она появилась.")
			return
		}
		log.Fatalf("query error: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-5s %-10s %-30s %-20s %s\n", "ID", "Date", "Title", "Comment", "Repeat")
	fmt.Println("------------------------------------------------------------------------------")

	count := 0
	for rows.Next() {
		var id int64
		var date, title, comment, repeat string
		if err := rows.Scan(&id, &date, &title, &comment, &repeat); err != nil {
			log.Fatalf("scan error: %v", err)
		}
		fmt.Printf("%-5d %-10s %-30s %-20s %s\n", id, date, title, comment, repeat)
		count++
	}

	if count == 0 {
		fmt.Println("Задач пока нет.")
	}
}
