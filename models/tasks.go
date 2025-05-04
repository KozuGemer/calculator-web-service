package models

type Task struct {
	ID         int      `json:"id"`               // Уникальный идентификатор задачи
	UserID     int      `json:"user_id"`          // Идентификатор пользователя, которому принадлежит задача
	Expression string   `json:"expression"`       // Математическое выражение
	Result     *float64 `json:"result,omitempty"` // Результат вычислений, может быть пустым
	Status     string   `json:"status"`           // Статус задачи (например, "pending" или "done")
}
