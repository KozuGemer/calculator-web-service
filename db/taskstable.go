package db

import (
	"database/sql"
	"fmt"

	"github.com/KozuGemer/calculator-web-service/models"
)

func GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT id, login, password FROM users WHERE login = ?", login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}
	return &user, nil
}

// Создание задачи в базе данных
func CreateTask(TaskID int, UserID int, expression string) (*models.Task, error) {
	stmt, err := DB.Prepare("INSERT INTO tasks (task_id, user_id, expression, status) VALUES (?, ?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}

	// Статус задачи по умолчанию "pending"
	_, err = stmt.Exec(TaskID, UserID, expression, "pending")
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	// Получаем ID только что созданной задачи
	err = DB.QueryRow("SELECT task_id FROM tasks WHERE user_id = ? AND expression = ? ORDER BY task_id DESC LIMIT 1", UserID, expression).Scan(&TaskID)
	if err != nil {
		return nil, fmt.Errorf("error fetching task id: %v", err)
	}

	// Возвращаем задачу
	return &models.Task{
		ID:         TaskID,
		UserID:     UserID, // Привязываем задачу к пользователю
		Expression: expression,
		Status:     "pending",
	}, nil
}

// Обновление задачи в базе данных
func UpdateTaskStatus(taskID int, status string, result float64) error {
	stmt, err := DB.Prepare("UPDATE tasks SET status = ?, result = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}

	_, err = stmt.Exec(status, result, taskID)
	if err != nil {
		return fmt.Errorf("error executing statement: %v", err)
	}

	return nil
}

func GetNextTask() (*models.Task, error) {
	// Запрос для получения задачи со статусом "pending"
	var task models.Task
	err := DB.QueryRow(`
		SELECT id, user_id, expression, status 
		FROM tasks 
		WHERE status = 'pending' 
		LIMIT 1
	`).Scan(&task.ID, &task.UserID, &task.Expression, &task.Status)

	// Если ошибка при извлечении задачи, вернем ошибку
	if err != nil {
		if err == sql.ErrNoRows {
			// Нет задач со статусом "pending"
			return nil, fmt.Errorf("no pending tasks")
		}
		return nil, err
	}

	return &task, nil
}
