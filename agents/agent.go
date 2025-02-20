package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/KozuGemer/calculator-web-service/utils"
)

type Task struct {
	ID         string   `json:"id"`
	Expression string   `json:"expression"`
	Result     *float64 `json:"result,omitempty"`
	Status     string   `json:"status"`
}

func fetchTask(serverURL string) (*Task, error) {
	resp, err := http.Get(serverURL + "/api/v1/tasks/next")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no tasks available")
	}

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}

func sendResult(serverURL string, task *Task) error {
	data, err := json.Marshal(map[string]float64{"result": *task.Result}) // Создаём JSON
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", serverURL+"/api/v1/tasks/completetask?id="+task.ID,
		bytes.NewBuffer(data)) // Передаём JSON в тело запроса
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json") // Указываем, что это JSON

	client := &http.Client{}
	resp, err := client.Do(req) // Отправляем запрос
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send result, status: %d", resp.StatusCode)
	}

	return nil
}

func calculate(expression string) float64 {
	result, err := utils.Calc(expression)
	if err != nil {
		fmt.Println("Calculation error:", err)
		return 0 // или обработать ошибку иначе
	}
	return result
}

func StartAgent(serverURL string) {
	for {
		task, err := fetchTask(serverURL)
		if err != nil {
			continue
		}

		result := calculate(task.Expression)
		task.Result = &result

		if err := sendResult(serverURL, task); err != nil {
			fmt.Println("Error sending result:", err)
		}
	}
}
