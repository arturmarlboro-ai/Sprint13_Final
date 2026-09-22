// users.go (и другие сущности) — CRUD-методы вроде GetUserByID,
package db

import (
	"context"
	"database/sql"
)

type User struct {
	ID   int
	Name string
}

// Предполагается, что у тебя где-то есть глобальный или передаваемый *sql.DB
func GetUserByID(ctx context.Context, id int) (User, error) {
	var u User
	err := dbGlobal.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name)
	if err == sql.ErrNoRows {
		// можно вернуть свою ошибку или nil + err
	}
	return u, err
}
