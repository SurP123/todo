package storage

import "database/sql"

type Task struct {
	ID     int    `json:"ID"`
	UserID int    `json:"UserID"`
	Text   string `json:"Text"`
	Done   bool   `json:"Done"`
}

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

func (s *Storage) GetAllInf(userID int) []Task {
	rows, err := s.db.Query(
		"SELECT id, user_id, text, done FROM tasks WHERE user_id = ?", userID,
	)
	if err != nil {
		return []Task{}
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var t Task
		rows.Scan(&t.ID, &t.UserID, &t.Text, &t.Done)
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Storage) Add(t Task) int {
	res, err := s.db.Exec(
		"INSERT INTO tasks (user_id, text, done) VALUES (?, ?, ?)",
		t.UserID, t.Text, t.Done,
	)
	if err != nil {
		return 0
	}
	id, _ := res.LastInsertId()
	return int(id)
}

func (s *Storage) Update(id int, userID int) bool {
	res, err := s.db.Exec(
		"UPDATE tasks SET done = NOT done WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return false
	}
	rows, _ := res.RowsAffected()
	return rows > 0
}

func (s *Storage) Delete(id int, userID int) bool {
	res, err := s.db.Exec(
		"DELETE FROM tasks WHERE id = ? AND user_id = ?",
		id, userID,
	)
	if err != nil {
		return false
	}
	rows, _ := res.RowsAffected()
	return rows > 0
}
