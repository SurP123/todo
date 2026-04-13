package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo/storage"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
	tasks := store.GetAllInf(userID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int)
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
	userID := r.Context().Value(userIDKey).(int)

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
	userID := r.Context().Value(userIDKey).(int)

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
