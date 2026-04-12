package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/storage"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil || userID == 0 {
		http.Error(w, "не указан user_id", http.StatusBadRequest)
		return
	}

	tasks := store.GetAllInf(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask storage.Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newTask.UserID == 0 {
		http.Error(w, "не указан UserID", http.StatusBadRequest)
		return
	}

	id := store.Add(newTask)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(id)
}

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	userIDStr := r.URL.Query().Get("user_id")

	if idStr == "" || userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
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
	userIDStr := r.URL.Query().Get("user_id")

	if idStr == "" || userIDStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
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
