package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KozuGemer/calculator-web-service/agents"
	"github.com/KozuGemer/calculator-web-service/db"
	"github.com/KozuGemer/calculator-web-service/handlers"
	"github.com/KozuGemer/calculator-web-service/models"
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

// loginPageHandler отдает HTML-страницу логина (без проверки токена)
func allTasksPageHandler(w http.ResponseWriter, r *http.Request) {
	data, err := embeddedFiles.ReadFile("site/alltasks.html")
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
		handlers.CreateTaskHandler(w, r)
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
	} else if token != "" && r.URL.Path == "/style/alltasks.css" {
		fileName = "alltasks.css"
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
	} else if token != "" && r.URL.Path == "/jsscripts/alltasks.js" {
		NameJS = "alltasks.js"
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

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now().Add(-24 * time.Hour), // Устанавливаем прошедшую дату для удаления cookie
	})
	w.Write([]byte("Logout successful"))
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
	// Получаем login пользователя из токена
	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "Error extracting token from cookie", http.StatusUnauthorized)
		return
	}

	// Разбираем JWT и получаем claims
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Получаем login из токена
	login, ok := claims["login"].(string) // предполагается, что login - строка
	if !ok {
		http.Error(w, "Invalid token claims", http.StatusUnauthorized)
		return
	}

	// Получаем user_id пользователя по login
	userID, err := getUserIDByLogin(login)
	if err != nil {
		http.Error(w, "Error fetching user ID", http.StatusInternalServerError)
		return
	}

	// Извлекаем задачу по task_id и user_id из базы данных
	task, err := db.GetTaskByIDAndUserID(id, userID)
	if err != nil {
		http.Error(w, "Error fetching task", http.StatusInternalServerError)
		return
	}

	// Если задача не найдена или она не принадлежит пользователю
	if task == nil {
		http.Error(w, "Task not found or does not belong to the user", http.StatusNotFound)
		return
	}

	// Если задача уже завершена, возвращаем ее результат
	if task.Status == "done" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(task) // Отправляем результат задачи
		return
	}

	// Если задача не завершена, продолжаем процесс обновления
	var req struct {
		Result float64 `json:"result"`
	}

	// Читаем результат из тела запроса
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Обновляем статус задачи в базе данных и сохраняем результат
	err = db.UpdateTaskStatus(task.ID, "done", req.Result)
	if err != nil {
		http.Error(w, "Error updating task", http.StatusInternalServerError)
		return
	}

	// Обновляем задачу в памяти (если необходимо)
	task.Status = "done"
	task.Result = &req.Result

	// Отправляем обновленную задачу
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Мы сохраняем задачи для каждого пользователя, чтобы не делать запрос несколько раз
var cachedTasks = make(map[int][]models.Task)

// Получить user_id по login
func getUserIDByLogin(login string) (int, error) {
	// Открытие базы данных
	db, err := sql.Open("sqlite3", "./calculator.db")
	if err != nil {
		return 0, fmt.Errorf("error opening database: %v", err)
	}
	defer db.Close()

	// Выполняем SQL запрос для поиска user_id по login
	var userID int
	err = db.QueryRow("SELECT id FROM users WHERE login = ?", login).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Если пользователя с таким login нет
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("error querying user_id: %v", err)
	}

	return userID, nil
}

func getAllExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	token, err := extractToken(r)
	if err != nil {
		http.Error(w, "Error extracting token from cookie", http.StatusUnauthorized)
		return
	}

	// Получаем login пользователя из токена
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

	// Получаем user_id для этого login
	userID, err := getUserIDByLogin(login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем все задачи для пользователя с состоянием "done"
	tasks, err := db.GetTasksByStatus(userID, "done")
	if err != nil {
		http.Error(w, "Error fetching tasks", http.StatusInternalServerError)
		return
	}

	// Возвращаем задачи пользователю в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func main() {
	// Инициализация базы данных
	db.InitDB()
	http.HandleFunc("/api/v1/tasks", handlers.CreateTaskHandler)
	http.HandleFunc("/api/v1/tasks/status", getTaskStatusHandler)
	http.HandleFunc("/api/v1/tasks/complete", completeTaskHandler)
	http.HandleFunc("/api/v1/expressions", getAllExpressionsHandler)
	http.HandleFunc("/alltasks", allTasksPageHandler)
	http.HandleFunc("/logout", logout)
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
	go agents.StartAgent()

	fmt.Println("Server is running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Server error:", err)
	}
}
