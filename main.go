package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/KozuGemer/calculator-web-service/agents"
)

type Task struct {
	ID         string   `json:"id"`
	Expression string   `json:"expression"`
	Result     *float64 `json:"result,omitempty"`
	Status     string   `json:"status"`
}

var (
	taskQueue = make(map[string]*Task)
	queueLock sync.Mutex
)

func generateTaskID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("task-%d", rand.Intn(1000000))
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	taskID := generateTaskID()
	task := &Task{
		ID:         taskID,
		Expression: req.Expression,
		Status:     "pending",
	}

	queueLock.Lock()
	taskQueue[taskID] = task
	queueLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":         taskID,
		"expression": task.Expression,
		"status":     "201 - Accepted for Processing",
		"message":    "Task has been created and is being processed",
	})

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("site/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	queueLock.Lock()
	task, exists := taskQueue[id]
	queueLock.Unlock()

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         id,
			"expression": nil,
			"result":     nil,
			"status":     "404 - Not Found",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func getNextTaskHandler(w http.ResponseWriter, r *http.Request) {
	queueLock.Lock()
	defer queueLock.Unlock()

	for _, task := range taskQueue {
		if task.Status == "pending" {
			task.Status = "processing"
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	http.Error(w, "No pending tasks", http.StatusNotFound)
}

func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	queueLock.Lock()
	task, exists := taskQueue[id]
	queueLock.Unlock()

	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         id,
			"expression": nil,
			"result":     nil,
			"status":     "404 - Not Found",
		})
		return
	}

	var req struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "400 - Bad Request",
			"message": "Invalid JSON format",
		})
		return
	}

	queueLock.Lock()
	if task.Status == "processing" {
		task.Result = &req.Result
		task.Status = "done"
	}
	queueLock.Unlock()
	if task.Status == "done" {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         task.ID,
			"expression": task.Expression,
			"result":     task.Result,
			"status":     "200 - Task Already Completed",
		})
		return
	} else {
		w.WriteHeader(210)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":         task.ID,
			"expression": task.Expression,
			"result":     task.Result,
			"status":     "210 - OK",
		})
	}

}

// Получение всех задач (новая фича)
func getAllExpressions(w http.ResponseWriter, r *http.Request) {
	queueLock.Lock()
	defer queueLock.Unlock()

	expressions := make([]Task, 0, len(taskQueue))
	for _, task := range taskQueue {
		expressions = append(expressions, *task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"expressions": expressions,
	})
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/v1/tasks", createTaskHandler)
	http.HandleFunc("/api/v1/tasks/status", getTaskStatusHandler)
	http.HandleFunc("/api/v1/tasks/next", getNextTaskHandler)
	http.HandleFunc("/api/v1/tasks/completetask", completeTaskHandler)
	http.HandleFunc("/api/v1/expressions", getAllExpressions) // Новый маршрут
	http.Handle("/style.css", http.FileServer(http.Dir("site")))
	go agents.StartAgent("http://localhost:8080")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
