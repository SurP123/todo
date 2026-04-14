package storage

import (
	"context"
)

type Task struct {
	ID     int    `json:"ID"`
	UserID int    `json:"UserID"`
	Text   string `json:"Text"`
	Done   bool   `json:"Done"`
}

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) GetAllInf(userID int) ([]Task, error) {
	rows, err := Pool.Query(context.Background(),
		"SELECT id, user_id, text, done FROM tasks WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		err := rows.Scan(&t.ID, &t.UserID, &t.Text, &t.Done)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (s *Storage) Add(t Task) (int, error) {
	var id int
	err := Pool.QueryRow(context.Background(),
		"INSERT INTO tasks (user_id, text, done) VALUES ($1, $2, $3) RETURNING id",
		t.UserID, t.Text, t.Done).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *Storage) Update(id int, userID int) (bool, error) {
	tag, err := Pool.Exec(context.Background(),
		"UPDATE tasks SET done = NOT done WHERE id = $1 AND user_id = $2",
		id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}

func (s *Storage) Delete(id int, userID int) (bool, error) {
	tag, err := Pool.Exec(context.Background(),
		"DELETE FROM tasks WHERE id = $1 AND user_id = $2",
		id, userID)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
