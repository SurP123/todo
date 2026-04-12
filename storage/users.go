package storage

import "database/sql"

type User struct {
	Id    int    `json:"ID"`
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type UsersStorage struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *UsersStorage {
	return &UsersStorage{db: db}
}

func (s *UsersStorage) AddUser(u User) (int, error) {
	res, err := s.db.Exec(
		"INSERT INTO users (login, pass) VALUES (?, ?)",
		u.Login, u.Pass,
	)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return int(id), nil
}

func (s *UsersStorage) Find(login string) (User, bool) {
	var u User
	err := s.db.QueryRow(
		"SELECT id, login, pass FROM users WHERE login = ?", login,
	).Scan(&u.Id, &u.Login, &u.Pass)

	if err == sql.ErrNoRows {
		return User{}, false
	}
	if err != nil {
		return User{}, false
	}
	return u, true
}
