package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"todo/storage"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "нет токена", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := ParseJWT(tokenStr)
	if err != nil {
		http.Error(w, "невалидный токен", http.StatusUnauthorized)
		return
	}

	tasks := store.GetAllInf(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "нет токена", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := ParseJWT(tokenStr)
	if err != nil {
		http.Error(w, "невалидный токен", http.StatusUnauthorized)
		return
	}
	var newTask storage.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newTask.UserID = userID
	id := store.Add(newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "нет токена", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := ParseJWT(tokenStr)
	if err != nil {
		http.Error(w, "невалидный токен", http.StatusUnauthorized)
		return
	}

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !store.Update(id, userID) {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "нет токена", http.StatusUnauthorized)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	userID, err := ParseJWT(tokenStr)
	if err != nil {
		http.Error(w, "невалидный токен", http.StatusUnauthorized)
		return
	}

	if idStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !store.Delete(id, userID) {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
