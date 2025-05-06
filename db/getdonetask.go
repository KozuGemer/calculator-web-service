package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/KozuGemer/calculator-web-service/models"
	_ "github.com/mattn/go-sqlite3" // Подключение драйвера для SQLite
)

// OpenDB открывает соединение с базой данных
func OpenDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./calculator.db") // Путь к вашей базе данных
	if err != nil {
		return nil, err
	}
	return db, nil
}
func GetTaskByIDAndUserID(TaskID string, UserID int) (*models.Task, error) {
	// Открытие базы данных
	db, err := sql.Open("sqlite3", "./calculator.db") // Путь к вашей базе данных
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	var task models.Task
	err = db.QueryRow("SELECT task_id, user_id, expression, result, status FROM tasks WHERE task_id = ? AND user_id = ?", TaskID, UserID).
		Scan(&task.ID, &task.UserID, &task.Expression, &task.Result, &task.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task not found")
		}
		return nil, fmt.Errorf("error fetching task: %v", err)
	}

	return &task, nil
}

// GetTasksByStatus возвращает все задачи с заданным статусом
func GetTasksByStatus(userID int, status string) ([]models.Task, error) {
	// Открытие базы данных
	db, err := sql.Open("sqlite3", "./calculator.db") // Путь к вашей базе данных
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	defer func() {
		if r := recover(); r != nil {
			log.Panicf("Recovered from panic: %s", r)
		}
	}()

	// Выполняем SQL запрос
	query := "SELECT task_id, user_id, expression, result, status FROM tasks WHERE status = ? AND user_id = ?"

	rows, err := db.Query(query, status, userID)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task

	// Чтение всех строк
	for rows.Next() {
		var task models.Task

		// Используем sql.NullFloat64 для поддержки null значений
		var result sql.NullFloat64

		// Чтение значений из строки
		if err := rows.Scan(&task.ID, &task.UserID, &task.Expression, &result, &task.Status); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Если результат не null, присваиваем значение
		if result.Valid {
			task.Result = &result.Float64
		} else {
			task.Result = nil
		}

		// Добавляем задачу в срез
		tasks = append(tasks, task)
	}

	// Проверяем на ошибку после чтения всех строк
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	// Если не найдено задач
	if len(tasks) == 0 {
		fmt.Println("No tasks found for the given user and status")
	}

	// Возвращаем задачи
	return tasks, nil
}
