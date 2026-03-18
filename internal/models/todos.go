package models

import (
	"database/sql"
	"errors"
	"time"
)

type Todo struct {
	ID        int
	Title     string
	Created   time.Time
	Completed time.Time
}

type TodoModel struct {
	DB *sql.DB
}

func (m *TodoModel) Insert(title string) error {
	stmt := `INSERT INTO todos (title, created) VALUES(?, ?)`

	_, err := m.DB.Exec(stmt, title, time.Now().UTC().Unix())
	if err != nil {
		return err
	}
	return nil
}

func (m *TodoModel) InProgress() ([]Todo, error) {
	stmt := `SELECT id, title, created FROM todos WHERE completed IS NULL ORDER BY created DESC`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
	}

	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo

		err := rows.Scan(&t.ID, &t.Title, &t.Created)
		if err != nil {
			return nil, err
		}

		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
