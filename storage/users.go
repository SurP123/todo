package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

type User struct {
	Id    int    `json:"ID"`
	Login string `json:"login"`
	Pass  string `json:"pass"`
}

type UsersStorage struct{}

func NewUsers() *UsersStorage {
	return &UsersStorage{}
}

func (s *UsersStorage) AddUser(u User) (int, error) {
	var id int
	err := Pool.QueryRow(context.Background(),
		"INSERT INTO users (login, pass) VALUES ($1, $2) RETURNING id",
		u.Login, u.Pass).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UsersStorage) Find(login string) (User, bool, error) {
	var u User
	err := Pool.QueryRow(context.Background(),
		"SELECT id, login, pass FROM users WHERE login = $1", login,
	).Scan(&u.Id, &u.Login, &u.Pass)

	if errors.Is(err, pgx.ErrNoRows) {
		return User{}, false, nil
	}
	if err != nil {
		return User{}, false, err
	}
	return u, true, nil
}
