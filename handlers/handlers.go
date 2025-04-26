package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/KozuGemer/calculator-web-service/db"
	"github.com/KozuGemer/calculator-web-service/models"
	"github.com/dgrijalva/jwt-go"
)

var secretKey = []byte("mysecretkey")

// Регистрация пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Регистрация в базе данных
	user, err := db.RegisterUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Error registering user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// Вход пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Проверка пользователя в базе данных
	user, err := db.AuthenticateUser(req.Login, req.Password)
	if err != nil {
		http.Error(w, "Invalid login or password", http.StatusUnauthorized)
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+tokenString)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
