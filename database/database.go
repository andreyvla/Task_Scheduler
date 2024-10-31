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

func GetTaskByID(db *sql.DB, id string) (*Task, error) {
	var task Task
	var dateString string
	query := "SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?"
	err := db.QueryRow(query, id).Scan(&task.ID, &dateString, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Задача не найдена")
		}
		log.Println("Ошибка при выполнении запроса:", err)
		return nil, fmt.Errorf("Ошибка при получении задачи")
	}

	// Преобразуем строку даты в тип time.Time
	task.Date, err = time.Parse("20060102", dateString)
	if err != nil {
		return nil, fmt.Errorf("Ошибка при преобразовании даты")
	}

	return &task, nil
}

func UpdateTask(db *sql.DB, id int, date time.Time, title, comment, repeat string) error {
	query := `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ? 
		WHERE id = ?`

	result, err := db.Exec(query, date.Format("20060102"), title, comment, repeat, id)
	if err != nil {
		return fmt.Errorf("ошибка при обновлении задачи: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка при получении количества обновленных строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("Задача не найдена")
	}

	return nil
}

func UpdateTaskDate(db *sql.DB, taskID int, newDate time.Time) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	_, err := db.Exec(query, newDate.Format("20060102"), taskID)
	return err
}
func DeleteTask(db *sql.DB, taskID int) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	_, err := db.Exec(query, taskID)
	return err
}
