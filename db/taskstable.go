package db

import (
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
func CreateTask(userID int, expression string) (*models.Task, error) {
	stmt, err := DB.Prepare("INSERT INTO tasks (user_id, expression, status) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}

	// Статус задачи по умолчанию "pending"
	_, err = stmt.Exec(userID, expression, "pending")
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %v", err)
	}

	// Получаем ID только что созданной задачи
	var taskID int
	err = DB.QueryRow("SELECT id FROM tasks WHERE user_id = ? AND expression = ? ORDER BY id DESC LIMIT 1", userID, expression).Scan(&taskID)
	if err != nil {
		return nil, fmt.Errorf("error fetching task id: %v", err)
	}

	// Возвращаем задачу
	return &models.Task{
		ID:         taskID,
		UserID:     userID, // Привязываем задачу к пользователю
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

// Получение задач пользователя из базы данных
func GetTasksByUser(userID int) ([]models.Task, error) {
	rows, err := DB.Query("SELECT id, expression, result, status FROM tasks WHERE user_id = ?", userID)
	if err != nil {
		return nil, fmt.Errorf("error querying tasks: %v", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		if err := rows.Scan(&task.ID, &task.Expression, &task.Result, &task.Status); err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	return tasks, nil
}
