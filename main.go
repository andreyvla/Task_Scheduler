package main

import (
	"database/sql"
	"go_final_project/database"
	"go_final_project/handlers"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	dbFile := filepath.Join(".", "scheduler.db")
	_, err := os.Stat(dbFile)

	if os.IsNotExist(err) {
		log.Println("Database file does not exist. Creating a new one.")
		database.CreateDB(dbFile)
	} else if err != nil {
		log.Fatalf("Ошибка при проверке файла базы данных: %v", err)
	} else {
		log.Println("Database file exists.")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	mux := http.NewServeMux()
	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))
	mux.HandleFunc("/api/nextdate", handlers.NextDateHandler)
	mux.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.TaskHandler(w, r, db)
	})
	mux.HandleFunc("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			handlers.GetTasks(w, r, db)
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandlePostTaskDone(w, r, db)
	})
	err = http.ListenAndServe(":7540", mux)
	if err != nil {
		panic(err)
	}
}
