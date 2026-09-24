// users.go (и другие сущности) — CRUD-методы вроде GetUserByID,
// pkg/db/repository.go
package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      int64
	Date    time.Time
	Title   string
	Comment string
	Repeat  bool
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task *Task) (int64, error) {

	dateStr := task.Date.Format(time.RFC3339)

	res, err := r.db.ExecContext(ctx,
		"INSERT INTO tasks (date, title, comment, repeat) VALUES (?, ?, ?, ?)",
		dateStr,
		task.Title,
		task.Comment,
		task.Repeat,
	)
	if err != nil {
		return 0, fmt.Errorf("insert task: %w", err)
	}

	// Шаг 2: получаем ID через last_insert_rowid()
	var id int64
	row := r.db.QueryRowContext(ctx, "SELECT last_insert_rowid()")
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("get last insert id: %w", err)
	}
	fmt.Println(res)
	return id, nil
}

// сюда добавь List, Get, Update и т.д. по аналогии
