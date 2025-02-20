package agents

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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
	_, err := json.Marshal(map[string]float64{"result": *task.Result})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", serverURL+"/api/v1/tasks/completetask?id="+task.ID, nil)
	if err != nil {
		return err
	}
	req.Body = http.NoBody
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send result")
	}
	return nil
}

func calculate(expression string) float64 {
	// Здесь будет вызов функции Calc
	// Заглушка, пока нет реальной обработки:
	return 42 // TODO: заменить на реальный расчет
}

func StartAgent(serverURL string) {
	for {
		task, err := fetchTask(serverURL)
		if err != nil {
			fmt.Println("No tasks, retrying in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}

		result := calculate(task.Expression)
		task.Result = &result

		if err := sendResult(serverURL, task); err != nil {
			fmt.Println("Error sending result:", err)
		}
	}
}
