package models

import (
	"database/sql"
	"time"
)

type Todo struct {
	ID        int
	Title     string
	Created   time.Time
	Completed bool
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

	rows, err := m.DB.Query(`SELECT id, title, created, completed FROM todos WHERE completed = FALSE ORDER BY created DESC`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var todos []Todo

	for rows.Next() {
		var t Todo
		var createdUnix int64
		err := rows.Scan(&t.ID, &t.Title, &createdUnix, &t.Completed)
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

func (m *TodoModel) Complete(id int) error {
	stmt := `UPDATE todos SET completed = TRUE WHERE id = ?`

	result, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNoRecord
	}

	return nil
}
