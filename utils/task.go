package utils

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("your-secret-key")

// GenerateTaskID генерирует случайный ID для задачи
func GenerateTaskID() int {
	rand.Seed(time.Now().UnixNano())
	random := rand.Intn(1000000)
	return random
}
func GetLoginFromToken(r *http.Request) (string, error) {
	// Получаем токен из cookie
	tokenFromCookie, err := ExtractToken(r)
	if err != nil {
		return "", fmt.Errorf("Error extracting token from cookie: %v", err)
	}

	// Получаем login пользователя из токена
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenFromCookie, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("Invalid token: %v", err)
	}

	login, ok := claims["login"].(string)
	if !ok {
		return "", fmt.Errorf("Invalid token claims")
	}

	return login, nil
}
