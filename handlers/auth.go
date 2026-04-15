package handlers

import (
	"encoding/json"
	"net/http"
	"todo/storage"
	"unicode"

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

	if len(req.Password) < 8 {
		http.Error(w, "пароль должен содержать 8 символов", http.StatusBadRequest)
		return
	}
	var UpperLetter bool = false
	for i := 0; i < len(req.Password); i++ {
		if unicode.IsUpper(rune(req.Password[i])) {
			UpperLetter = true
			break
		}
	}
	if !UpperLetter {
		http.Error(w, "Пароль должен содержать хотя бы одну заглавную букву", http.StatusBadRequest)
		return
	}
	_, exists, err := users.Find(req.Username)
	if err != nil {
		http.Error(w, "ошибка проверки пользователя", http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, "логин уже занят", http.StatusConflict)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "ошибка хеширования пароля", http.StatusInternalServerError)
		return
	}

	id, err := users.AddUser(storage.User{Login: req.Username, Pass: string(hash)})
	if err != nil {
		http.Error(w, "ошибка создания пользователя", http.StatusInternalServerError)
		return
	}

	token, err := CreateJWT(id)
	if err != nil {
		http.Error(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req authRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	user, exists, err := users.Find(req.Username)
	if err != nil {
		http.Error(w, "ошибка поиска пользователя", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "пользователь не найден", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(req.Password))
	if err != nil {
		http.Error(w, "неверный пароль", http.StatusUnauthorized)
		return
	}

	token, err := CreateJWT(user.Id)
	if err != nil {
		http.Error(w, "ошибка создания токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
