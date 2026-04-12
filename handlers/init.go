package handlers

import "todo/storage"

var store *storage.Storage
var users *storage.UsersStorage

func Init(s *storage.Storage, u *storage.UsersStorage) {
	store = s
	users = u
}
