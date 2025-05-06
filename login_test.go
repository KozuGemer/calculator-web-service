package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginAndTokenInCookies(t *testing.T) {
	recorder := httptest.NewRecorder()
	// Данные для входа
	reqBody := `{"login": "ddenisora", "password": "hered#gid821"}`

	// Создаем запрос для входа
	req, err := http.NewRequest("POST", "/api/v1/login", strings.NewReader(reqBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	handler := http.HandlerFunc(LoginHandler) // Обработчик для входа
	handler.ServeHTTP(recorder, req)

	// Проверка, что код ответа правильный
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %v", status)
	}

	// Проверка, что токен добавлен в cookie
	cookie := recorder.Result().Cookies()
	if len(cookie) == 0 || cookie[0].Name != "token" {
		t.Errorf("Expected token cookie, but got %v", cookie)
	}

	// Извлекаем токен для последующего использования
	token := cookie[0].Value
	t.Logf("Token in cookies: %s", token)

	// Для следующего теста мы используем этот токен
	if token == "" {
		t.Fatal("Token should not be empty")
	}
}
