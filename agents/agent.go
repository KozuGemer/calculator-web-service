package agents

import (
	"fmt"
	"time"

	"github.com/KozuGemer/calculator-web-service/db"
	"github.com/KozuGemer/calculator-web-service/utils"
)

// Старт агента, который будет выполнять задачи из базы данных
func StartAgent() {
	for {
		// Извлекаем задачу из базы данных
		task, err := db.GetNextTask() // Получаем следующую задачу
		if err != nil {
			continue
		}

		// Если задач нет, ждем некоторое время и пробуем снова
		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		// Выполняем вычисления по полученному выражению
		result := calculate(task.Expression)
		task.Result = &result

		// Обновляем статус задачи в базе данных
		err = db.UpdateTaskStatus(task.ID, "done", result)
		if err != nil {
			fmt.Println("Error updating task:", err)
			continue
		}

		// Выводим результат в лог
		fmt.Printf("Task %d completed. Result: %f\n", task.ID, result)
	}
}

// Функция для выполнения вычислений
func calculate(expression string) float64 {
	result, err := utils.Calc(expression) // Вычисляем результат выражения
	if err != nil {
		fmt.Println("Calculation error:", err)
		return 0 // Возвращаем 0 в случае ошибки
	}
	return result
}
