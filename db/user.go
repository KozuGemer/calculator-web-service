package db

import (
	"fmt"
	"log"
	"time"

	"github.com/KozuGemer/calculator-web-service/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Регистрация нового пользователя
func RegisterUser(login, password string) (*models.User, error) {
	log.Println("Регистрация пользователя:", login)
	// Проверяем, существует ли логин
	var existingLogin string
	err := DB.QueryRow("SELECT login FROM users WHERE login = ?", login).Scan(&existingLogin)
	if err == nil {
		return nil, fmt.Errorf("user with this login already exists")
	}

	// Генерируем токен
	token := generateJWTToken(login)

	// Добавляем пользователя в базу данных
	stmt, err := DB.Prepare("INSERT INTO users(login, password, token) VALUES(?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %v", err)
	}

	_, err = stmt.Exec(login, string(password), token)
	if err != nil {
		return nil, fmt.Errorf("error inserting user: %v", err)
	}

	return &models.User{Login: login, Password: string(password), Token: token}, nil
}

// Авторизация пользователя
func AuthenticateUser(login, password string) (*models.User, error) {
	log.Println("Аутентификация пользователя:", login)

	var user models.User
	var hashedPassword string

	// Чётко извлекаем хэш из базы данных
	err := DB.QueryRow("SELECT id, login, password, token FROM users WHERE login = ?", login).
		Scan(&user.ID, &user.Login, &hashedPassword, &user.Token)

	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}

	// Сравнение пароля с хэшем из базы
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}
	// Генерируем новый токен для пользователя
	newToken := generateJWTToken(login)

	// Обновляем токен в БД
	_, err = DB.Exec("UPDATE users SET token = ? WHERE login = ?", newToken, login)
	if err != nil {
		return nil, fmt.Errorf("error updating token")
	}
	user.Token = newToken
	return &user, nil
}

func UpdateUserToken(userID int, newToken string) error {
	// Обновляем токен пользователя по его ID
	stmt, err := DB.Prepare("UPDATE users SET token = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("error preparing statement: %v", err)
	}
	_, err = stmt.Exec(newToken, userID)
	if err != nil {
		return fmt.Errorf("error updating token for user %d: %v", userID, err)
	}
	return nil
}

// Функция для генерации JWT токена
func generateJWTToken(login string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": login,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString([]byte("your-secret-key"))
	return tokenString
}

// Получение токена для пользователя по логину
func GetUserToken(login string) (string, error) {
	var token string
	err := DB.QueryRow("SELECT token FROM users WHERE login = ?", login).Scan(&token)
	if err != nil {
		return "", fmt.Errorf("error fetching token for user: %v", err)
	}
	return token, nil
}
