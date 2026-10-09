package db

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"` // "20260924"
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"` // "d 1", "w 1,3,5", "m 1,15", "y", ""
}

func DeleteTask(ctx context.Context, id string) error {
	res, err := DB.ExecContext(ctx, `DELETE FROM scheduler WHERE id = ?`, id)
	if err != nil {
		log.Printf("DeleteTask error: %v", err)
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("DeleteTask error: %v", err)
		return err
	}
	if count == 0 {
		log.Printf("DeleteTask: no task found for id %s", id)
		return fmt.Errorf("no task found for id= %s", id)
	}
	log.Printf("deleted task with id= %s", id)
	return nil
}

func UpdateTask(ctx context.Context, task Task, id string) error {
	res, err := DB.ExecContext(ctx, `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`, task.Date, task.Title, task.Comment, task.Repeat, id)
	count, err := res.RowsAffected()
	if err != nil {
		log.Printf("UpdateTask error; %v", err)
		return err
	}
	if count == 0 {
		log.Printf("UpdateTask: no updates for id %s", id)
		return fmt.Errorf("no updates for id= %s", id)
	}
	log.Printf("updated task: %v with id= %s", task, id)
	return nil
}

func GetTaskByID(ctx context.Context, id string) (*Task, error) {
	row := DB.QueryRowContext(ctx, `SELECT CAST(id AS TEXT) as id, date, title, comment, repeat FROM scheduler WHERE id = ?`, id)

	var task Task

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		log.Printf("getTaskByID error Scanning: %v", err)
		return nil, err
	}
	log.Printf("Task found: %s, %s, %s, %s, %s", task.ID, task.Date, task.Title, task.Comment, task.Repeat)
	return &task, nil
}

func GetTasksBySearch(ctx context.Context, search string) ([]*Task, error) {

	tasks := []*Task{}

	if len(search) == 10 && search[2] == '.' && search[5] == '.' {
		t, err := time.Parse("02.01.2006", search)
		if err != nil {
			return nil, err
		}
		date := t.Format("20060102")
		rows, err := DB.QueryContext(ctx,
			`SELECT CAST(id AS TEXT) as id, date, title, comment, repeat FROM scheduler WHERE date = ? ORDER BY date`, date)
		if err != nil {
			return tasks, err
		}
		defer rows.Close()

		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return tasks, err
			}
			tasks = append(tasks, &task)
		}
		if err := rows.Err(); err != nil {
			return tasks, err
		}
		return tasks, nil

	} else {
		search2L := strings.ToLower(search)

		rows, err := DB.QueryContext(ctx,
			`SELECT CAST(id AS TEXT) as id, date, title, comment, repeat FROM scheduler ORDER BY date `)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var task Task
			err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
			if err != nil {
				return nil, err
			}
			if strings.Contains(strings.ToLower(task.Title), search2L) || strings.Contains(strings.ToLower(task.Comment), search2L) {
				tasks = append(tasks, &task)
			}
		}
		if err := rows.Err(); err != nil {
			log.Printf("rows.Err() returned: %v", err)
		}

	}
	return tasks, nil
}

func AddTask(ctx context.Context, task Task) (string, error) {
	if task.Date == "" {
		task.Date = time.Now().Format("20060102")
	}

	res, err := DB.ExecContext(ctx,
		`INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`,
		task.Date, task.Title, task.Comment, task.Repeat,
	)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	idStr := strconv.FormatInt(id, 10)
	log.Printf("INSERT INTO db: %s %s %s %s with id= %v", task.Date, task.Title, task.Comment, task.Repeat, id)
	return idStr, nil
}

func GetTasksByLimit(ctx context.Context, limit int) ([]*Task, error) {
	tasks := []*Task{}
	rows, err := DB.QueryContext(ctx, `SELECT CAST(id AS TEXT) as id, date, title, comment, repeat FROM scheduler LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		log.Printf("rows.Err() returned: %v", err)
	}

	return tasks, nil
}
