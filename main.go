package main

import (
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

	mux := http.NewServeMux()

	mux.Handle("GET /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.GetTaskHandler)))
	mux.Handle("POST /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.CreateTaskHandler)))
	mux.Handle("PUT /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.UpdateTaskHandler)))
	mux.Handle("DELETE /api/tasks", handlers.AuthMiddleware(http.HandlerFunc(handlers.DeleteTaskHandler)))
	mux.HandleFunc("POST /auth/login", handlers.LoginHandler)
	mux.HandleFunc("POST /auth/register", handlers.RegisterHandler)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/index.html")
	})
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates/tasks.html")
	})

	log.Println("Server is listening...")
	log.Fatal(http.ListenAndServe(":8181", handlers.LoggerMiddleware(mux)))
}
