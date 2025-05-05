package utils

import (
	"fmt"
	"net/http"
	"strings"
)

func ExtractToken(r *http.Request) (string, error) {
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
