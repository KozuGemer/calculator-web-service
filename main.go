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
	json.NewEncoder(w).Encode(map[string]string{"id": taskID})
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
	fmt.Println("Запрос статуса для задачи:", id)
	fmt.Println("Текущее состояние задач:", taskQueue)
	task, exists := taskQueue[id]
	queueLock.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
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
	var req struct {
		Result float64 `json:"result"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	queueLock.Lock()
	task, exists := taskQueue[id]
	if exists && task.Status == "processing" {
		task.Result = &req.Result
		task.Status = "done"
	}
	queueLock.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/v1/tasks", createTaskHandler)
	http.HandleFunc("/api/v1/tasks/status", getTaskStatusHandler)
	http.HandleFunc("/api/v1/tasks/next", getNextTaskHandler)
	http.HandleFunc("/api/v1/tasks/completetask", completeTaskHandler)
	http.Handle("/style.css", http.FileServer(http.Dir("site")))
	go agents.StartAgent("http://localhost:8080")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
