package db

import (
	"context"
	"log"
	"time"
)

type Task struct {
	ID      int64  `json:"id"`
	Date    string `json:"date"` // "20260924"
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"` // "d 1", "w 1,3,5", "m 1,15", "y", ""
}

func AddTask(ctx context.Context, task Task) (int64, error) {
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	res, err := DB.ExecContext(ctx,
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	log.Printf("INSERT INTO db: %s %s %s %s with id= %v", task.Date, task.Title, task.Comment, task.Repeat, id)
	return id, nil
}
