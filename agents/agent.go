package agents

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
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

	return &task, nil
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func sendResult(serverURL string, task *Task) error {
	// Используем пул буферов
	buf := bufferPool.Get().(*bytes.Buffer)
	defer bufferPool.Put(buf)
	buf.Reset()

	// Кодируем данные в буфер
	if err := json.NewEncoder(buf).Encode(map[string]float64{"result": *task.Result}); err != nil {
		return err
	}

	// Создаем URL с помощью strings.Builder
	var urlBuilder strings.Builder
	urlBuilder.WriteString(serverURL)
	urlBuilder.WriteString("/api/v1/tasks/completetask?id=")
	urlBuilder.WriteString(task.ID)
	url := urlBuilder.String()

	// Создаем запрос
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)

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

		if err := sendResult(serverURL, task); err != nil {
			fmt.Println("Error sending result:", err)
		}
	}
}
