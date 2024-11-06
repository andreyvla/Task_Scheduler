package main

import (
	"go_final_project/database"
	"go_final_project/handlers"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := database.CreateOrGetDb()
	if err != nil {
		log.Printf("Ошибка при создании или открытии базы данных: %v", err)
		return
	}
	defer db.Close()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	mux := http.NewServeMux()
	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", handlers.TaskHandler(db))
	mux.HandleFunc("/api/tasks", handlers.GetTasks(db))
	mux.HandleFunc("/api/task/done", handlers.HandlePostTaskDone(db))

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Printf("Ошибка при запуске сервера: %v", err)
	}
}
