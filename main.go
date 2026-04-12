package main

import (
	"fmt"
	"log"
	"net/http"
	"todo/handlers"
	"todo/storage"
)

func main() {
	db, err := storage.NewDB("todo.db")
	if err != nil {
		log.Fatal("ошибка подключения к базе:", err)
	}
	defer db.Close()

	handlers.Init(storage.NewStorage(db), storage.NewUsers(db))

	http.HandleFunc("GET /api/tasks", handlers.GetTaskHandler)
	http.HandleFunc("POST /api/tasks", handlers.CreateTaskHandler)
	http.HandleFunc("PUT /api/tasks", handlers.UpdateTaskHandler)
	http.HandleFunc("DELETE /api/tasks", handlers.DeleteTaskHandler)
	http.HandleFunc("POST /auth/login", handlers.LoginHandler)
	http.HandleFunc("POST /auth/register", handlers.RegisterHandler)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	http.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/tasks.html")
	})

	fmt.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8181", nil))
}
