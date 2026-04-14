package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/storage"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	tasks, err := store.GetAllInf(userID)
	if err != nil {
		http.Error(w, "ошибка получения задач", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	var newTask storage.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	newTask.UserID = userID

	id, err := store.Add(newTask)
	if err != nil {
		http.Error(w, "ошибка создания задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "параметр id обязателен", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный формат id", http.StatusBadRequest)
		return
	}

	updated, err := store.Update(id, userID)
	if err != nil {
		http.Error(w, "ошибка обновления задачи", http.StatusInternalServerError)
		return
	}
	if !updated {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "параметр id обязателен", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "неверный формат id", http.StatusBadRequest)
		return
	}

	deleted, err := store.Delete(id, userID)
	if err != nil {
		http.Error(w, "ошибка удаления задачи", http.StatusInternalServerError)
		return
	}
	if !deleted {
		http.Error(w, "задача не найдена", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
