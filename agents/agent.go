package agents

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/KozuGemer/calculator-web-service/utils"
)

type Task struct {
	ID         string   `json:"id"`
	Expression string   `json:"expression"`
	Result     *float64 `json:"result,omitempty"`
	Status     string   `json:"status"`
}

var httpClient = &http.Client{Timeout: 1 * time.Second} // Глобальный клиент

func fetchTask(serverURL string) (*Task, error) {
	startFetch := time.Now()

	resp, err := httpClient.Get(serverURL + "/api/v1/tasks/next") // Используем глобальный клиент
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("no tasks available")
	}

	body, err := io.ReadAll(resp.Body) // Читаем тело ответа одним вызовом
	if err != nil {
		return nil, err
	}

	var task Task
	if err := json.Unmarshal(body, &task); err != nil { // Декодируем JSON отдельно
		return nil, err
	}

	fmt.Println("fetchTask() time:", time.Since(startFetch)) // Логируем только в конце

	return &task, nil
}
func sendResult(serverURL string, task *Task) error {
	data, err := json.Marshal(map[string]float64{"result": *task.Result})
	if err != nil {
		return err
	}

	url := serverURL + "/api/v1/tasks/completetask?id=" + task.ID
	req, err := http.NewRequest("POST", url, strings.NewReader(string(data))) // Оптимизированная отправка
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	startSend := time.Now()
	resp, err := httpClient.Do(req) // Используем глобальный клиент
	fmt.Println("sendResult() time:", time.Since(startSend))

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
		// startFetch := time.Now()
		task, err := fetchTask(serverURL)
		// fmt.Println("fetchTask() time:", time.Since(startFetch))
		if err != nil {
			continue
		}

		// startCalc := time.Now()
		result := calculate(task.Expression)
		// fmt.Println("calculate() time:", time.Since(startCalc))

		task.Result = &result

		startSend := time.Now()
		if err := sendResult(serverURL, task); err != nil {
			fmt.Println("Error sending result:", err)
		}
		fmt.Println("sendResult() time:", time.Since(startSend))
	}
}
