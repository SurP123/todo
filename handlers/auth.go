package handlers

import (
	"encoding/json"
	"net/http"
	"todo/storage"

	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "заполни все поля", http.StatusBadRequest)
		return
	}

	_, exists := users.Find(req.Username)
	if exists {
		http.Error(w, "логин уже занят", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	id, err := users.AddUser(storage.User{Login: req.Username, Pass: string(hash)})
	if err != nil {
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	token, err := CreateJWT(id)
	if err != nil {
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, exists := users.Find(req.Username)
	if !exists {
		http.Error(w, "пользователь не найден", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password))
	if err != nil {
		http.Error(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token, err := CreateJWT(user.Id)
	if err != nil {
		http.Error(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
