package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KozuGemer/calculator-web-service/db"
	"github.com/KozuGemer/calculator-web-service/models"
	"github.com/KozuGemer/calculator-web-service/utils"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	fmt.Println(req.Expression)
	login, err := utils.GetLoginFromToken(r)
	if err != nil {
		fmt.Print(err)
	}
	// Получаем пользователя из базы данных
	user, err := db.GetUserByLogin(login)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// Генерация уникального ID задачи
	taskID := utils.GenerateTaskID()

	// Создаем новую задачу
	task := models.Task{
		ID:         taskID,
		UserID:     user.ID,
		Expression: req.Expression,
		Status:     "pending",
	}

	// Декодируем тело запроса в структуру
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		// В случае ошибки возвращаем клиенту ошибку
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Сохраняем задачу в базе данных
	_, err = db.CreateTask(taskID, user.ID, req.Expression)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	// Отправляем задачу в ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         task.ID,
		"expression": task.Expression,
		"status":     task.Status,
	})
}
