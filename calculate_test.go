package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/KozuGemer/calculator-web-service/handlers"
)

func TestSubmitTaskAndGetStatusWithToken(t *testing.T) {
	// 1. Вход: создаем запрос для логина с правильными данными
	reqBody := `{"login": "testuser", "password": "testpassword"}`

	req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем тестовый сервер
	recorder := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler) // Обработчик для входа
	handler.ServeHTTP(recorder, req)

	// Проверка кода ответа
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	// Проверка наличия токена в cookies
	cookie := recorder.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Name != "token" {
		t.Errorf("Expected token cookie, but got %v", cookie)
	}

	// Извлекаем токен из cookies
	token := cookie[0].Value
	if token == "" {
		t.Fatal("Token should not be empty")
	}

	// 2. Отправка задачи с токеном в cookies
	reqBody = `{"expression": "2+2"}`
	req, err = http.NewRequest("POST", "/api/v1/tasks", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "token", Value: token}) // Добавляем токен в cookies

	// Создаем тестовый сервер для создания задачи
	recorder = httptest.NewRecorder()
	handler = http.HandlerFunc(handlers.CreateTaskHandler) // Обработчик для создания задачи
	handler.ServeHTTP(recorder, req)

	// Проверка кода ответа на создание задачи
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	// Получаем ID задачи из ответа
	var response map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	taskID := response["task_id"].(string)
	if taskID == "" {
		t.Fatal("Expected task_id to be returned")
	}

	// 3. Получение статуса задачи с использованием полученного taskID
	req, err = http.NewRequest("GET", "/api/v1/complete?id="+taskID, nil)
	if err != nil {
		t.Fatal(err)
	}
	req.AddCookie(&http.Cookie{Name: "token", Value: token}) // Добавляем токен в cookies

	// Создаем тестовый сервер для получения статуса задачи
	recorder = httptest.NewRecorder()
	handler = http.HandlerFunc(getTaskStatusHandler) // Обработчик для получения статуса задачи
	handler.ServeHTTP(recorder, req)

	// Проверка кода ответа
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	// Проверка, что задача возвращена с правильными данными
	var taskResponse Task
	if err := json.NewDecoder(recorder.Body).Decode(&taskResponse); err != nil {
		t.Fatal(err)
	}

	if taskResponse.Expression != "2+2" {
		t.Errorf("Expected expression '2+2', but got %s", taskResponse.Expression)
	}
	if taskResponse.Status != "pending" {
		t.Errorf("Expected status 'pending', but got %s", taskResponse.Status)
	}
}
