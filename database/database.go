package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func CreateDB(dbFile string) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date TEXT NOT NULL DEFAULT "",
			title TEXT NOT NULL DEFAULT "",
			comment TEXT NOT NULL DEFAULT "",
			repeat VARCHAR(128) NOT NULL DEFAULT ""
		);
		CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date);
	`)
	if err != nil {
		log.Println(err)
		return
	}
}

type Task struct {
	ID      int
	Date    time.Time
	Title   string
	Comment string
	Repeat  string
}

// InsertTask вставляет новую задачу в базу данных и возвращает её идентификатор
func InsertTask(db *sql.DB, date time.Time, title, comment, repeat string) (int, error) {
	// Подготовка SQL-запроса для вставки задачи
	query := `
INSERT INTO scheduler (date, title, comment, repeat)
VALUES (?, ?, ?, ?)
	`

	// Выполнение запроса
	result, err := db.Exec(query, date.Format("20060102"), title, comment, repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка при вставке задачи: %w", err)
	}

	// Получение идентификатора созданной записи
	taskID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка при получении ID вставленной задачи: %w", err)
	}

	return int(taskID), nil
}
