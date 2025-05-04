package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KozuGemer/calculator-web-service/agents"
	"github.com/KozuGemer/calculator-web-service/db"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Task описывает задачу
type Task struct {
	ID         int      `json:"id"`               // Уникальный идентификатор задачи
	UserID     int      `json:"user_id"`          // Идентификатор пользователя, которому принадлежит задача
	Expression string   `json:"expression"`       // Математическое выражение
	Result     *float64 `json:"result,omitempty"` // Результат вычислений, может быть пустым
	Status     string   `json:"status"`           // Статус задачи (например, "pending" или "done")
}

var (
	taskQueue = make(map[string]*Task)
	queueLock sync.Mutex
)

var secretKey = []byte("your-secret-key") // Должен совпадать с ключом при генерации токена

//go:embed site/*
var embeddedFiles embed.FS

// generateTaskID генерирует случайный ID для задачи
func generateTaskID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("task-%d", rand.Intn(1000000))
}
func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := embeddedFiles.ReadFile("site/register.html")
	if err != nil {
		http.Error(w, "Error loading register page", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// extractToken пытается получить токен сначала из заголовка, затем из cookie "token"
func extractToken(r *http.Request) (string, error) {
	// Сначала пытаемся получить из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], nil
		}
	}
	// Если заголовка нет, пробуем получить токен из cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", fmt.Errorf("cookie 'token' not found")
	}
	return cookie.Value, nil
}

// compareTokensHandler будет сравнивать токен из cookie с токеном из базы данных
func compareTokensHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из cookie
	tokenFromCookie, err := extractToken(r)
	if err != nil {
		http.Error(w, "Error extracting token from cookie", http.StatusBadRequest)
		return
	}

	// Получаем login пользователя из токена
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenFromCookie, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	login, ok := claims["login"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Извлекаем токен из базы данных для данного пользователя
	dbToken, err := db.GetUserToken(login)
	if err != nil {
		http.Error(w, "Error fetching token from database", http.StatusInternalServerError)
		return
	}

	// Сравниваем токены
	if tokenFromCookie == dbToken {
		// Если токены одинаковы
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"result": "Tokens are identical"})
	} else {
		// Если токены разные
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"result": "Tokens are different"})
	}
}

// loginPageHandler отдает HTML-страницу логина (без проверки токена)
func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := embeddedFiles.ReadFile("site/login.html")
	if err != nil {
		http.Error(w, "Error loading login page", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}

// indexHandler отдает основную страницу (калькулятор) для авторизованных пользователей
func indexHandler(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	data, err := embeddedFiles.ReadFile("site/index.html")
	if err != nil {
		http.Error(w, "Error loading index page", http.StatusInternalServerError)
		return
	}
	if token != "" && r.URL.Path == "/api/v1/calculate" {
		createTaskHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(data)
}
func getTokenFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("token")
	if err != nil {
		// Если cookie не найдено, возвращаем пустую строку
		return ""
	}
	return cookie.Value
}

// styleHandler отдает CSS-стили
func styleHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем токен из cookie
	token := getTokenFromCookie(r)

	// Если токен есть, отдаем стиль для страницы с калькулятором, иначе - стиль для страницы логина
	var fileName string
	if token == "" && r.URL.Path == "/style/register.css" {
		fileName = "register.css" // Если токена нет и путь /register
	} else if token == "" {
		fileName = "login.css" // Если токен отсутствует, но путь не /register
	} else {
		fileName = "style.css" // Если токен есть, загружаем основной стиль
	}
	// Открываем файл из встроенной файловой системы
	data, err := embeddedFiles.ReadFile("site/" + fileName)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	// Устанавливаем правильный MIME тип для CSS
	w.Header().Set("Content-Type", "text/css")
	w.Write(data)
}

// jsHandler отдает JavaScript
func jsHandler(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromCookie(r)
	var NameJS string
	if token == "" && r.URL.Path == "/jsscripts/register.js" {
		NameJS = "register.js"
	} else {
		NameJS = "app.js"
	}

	data, err := embeddedFiles.ReadFile("site/" + NameJS)
	if err != nil {
		http.Error(w, "Error loading JS", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/javascript")
	w.Write(data)
}

// registerHandler отвечает за регистрацию пользователя
func registerHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли токен в cookie
	token := getTokenFromCookie(r)

	// Если токен существует, редиректим на главную страницу (не позволяем попасть на регистрацию)
	if token != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 12)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Регистрируем пользователя и сохраняем токен
	user, err := db.RegisterUser(creds.Login, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// loginHandler отвечает за авторизацию пользователя и генерацию JWT-токена
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем, есть ли токен в cookie
	token := getTokenFromCookie(r)

	// Если токен существует, редиректим на главную страницу (не позволяем попасть на страницу логина)
	if token != "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var creds struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	// Аутентификация пользователя
	user, err := db.AuthenticateUser(creds.Login, creds.Password)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid login or password"})
		return
	}

	// Генерируем JWT-токен с сроком действия 24 часа
	tokenStr := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": creds.Login,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := tokenStr.SignedString(secretKey)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error generating token"})
		return
	}

	// Устанавливаем токен в cookie
	cookie := &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)

	// Обновляем токен в базе данных
	user.Token = tokenString
	err = db.UpdateUserToken(user.ID, tokenString)
	if err != nil {
		http.Error(w, "Error updating token in database", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

// createTaskHandler создает новую задачу
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Expression string `json:"expression"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Получаем токен из cookie
	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "Error extracting token from cookie", http.StatusUnauthorized)
		return
	}

	// Извлекаем login из токена
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	login, ok := claims["login"].(string)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Получаем ID пользователя из базы данных
	user, err := db.GetUserByLogin(login)
	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	// Создаем задачу и сохраняем ее в базе данных
	task, err := db.CreateTask(user.ID, req.Expression)
	if err != nil {
		http.Error(w, "Error creating task", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// getTaskStatusHandler возвращает статус задачи
func getTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	queueLock.Lock()
	task, exists := taskQueue[id]
	queueLock.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// completeTaskHandler завершает задачу, устанавливая результат
func completeTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	// Блокируем очередь задач для безопасного доступа
	queueLock.Lock()
	task, exists := taskQueue[id]
	queueLock.Unlock()

	if !exists {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Получаем login пользователя из токена
	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "Error extracting token from cookie", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	login, ok := claims["login"].(int)
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Проверяем, что задача принадлежит пользователю
	if task.UserID != login {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Result float64 `json:"result"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Обновляем задачу в базе данных
	err = db.UpdateTaskStatus(task.ID, "done", req.Result)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	// Отправляем обновленную задачу
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// getAllExpressionsHandler возвращает все задачи
func getAllExpressionsHandler(w http.ResponseWriter, r *http.Request) {
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
	// Инициализация базы данных
	db.InitDB()
	http.HandleFunc("/api/v1/tasks", createTaskHandler)
	http.HandleFunc("/api/v1/tasks/status", getTaskStatusHandler)
	http.HandleFunc("/api/v1/tasks/complete", completeTaskHandler)
	http.HandleFunc("/api/v1/expressions", getAllExpressionsHandler)

	http.HandleFunc("/register", registerPageHandler) // Страница регистрации
	// Регистрация маршрутов
	http.HandleFunc("/api/v1/register", registerHandler)
	http.HandleFunc("/api/v1/login", loginHandler)
	http.HandleFunc("/login", loginPageHandler) // Страница логина
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/style/", styleHandler) // Главная страница (калькулятор)
	http.HandleFunc("/jsscripts/", jsHandler)
	// Новый обработчик для сравнения токенов
	http.HandleFunc("/api/v1/compare-tokens", compareTokensHandler)

	// Запуск агента
	go agents.StartAgent("http://localhost:8080")

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
