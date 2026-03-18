package models

import (
	"database/sql"
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

	rows, err := m.DB.Query(`SELECT id, title, created FROM todos WHERE completed IS NULL ORDER BY created DESC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo
		var createdUnix int64
		err := rows.Scan(&t.ID, &t.Title, &createdUnix)
		if err != nil {
			return nil, err
		}

		t.Created = time.Unix(createdUnix, 0).UTC()
		todos = append(todos, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
