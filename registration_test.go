package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterAccount(t *testing.T) {
	// Регистрация нового пользователя
	reqBody := `{"login": "newuser", "password": "newpassword"}`

	// Создаем новый POST-запрос
	req, err := http.NewRequest("POST", "/api/v1/register", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаем новый ResponseRecorder для имитации ответа сервера
	recorder := httptest.NewRecorder()

	// Обработчик для регистрации (должен быть инициализирован)
	handler := http.HandlerFunc(RegisterHandler)

	// Используем ServeHTTP для обработки запроса
	handler.ServeHTTP(recorder, req)

	// Проверка статуса ответа
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	// Декодируем ответ
	var response map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	// Проверка успешной регистрации
	if response["message"] != "Registration successful" {
		t.Errorf("Expected 'Registration successful', but got %v", response["message"])
	}
}
