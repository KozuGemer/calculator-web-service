package db

import (
	"fmt"

	"github.com/KozuGemer/calculator-web-service/models"
)

// Регистрация нового пользователя
func RegisterUser(login, password string) (*models.User, error) {
	// Проверим, не существует ли уже такой логин
	var existingLogin string
	err := DB.QueryRow("SELECT login FROM users WHERE login = ?", login).Scan(&existingLogin)
	if err == nil {
		return nil, fmt.Errorf("user with this login already exists")
	}

	// Добавляем нового пользователя в базу данных
	stmt, err := DB.Prepare("INSERT INTO users(login, password) VALUES(?, ?)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec(login, password)
	if err != nil {
		return nil, err
	}

	// Возвращаем созданного пользователя
	var user models.User
	err = DB.QueryRow("SELECT id, login FROM users WHERE login = ?", login).Scan(&user.ID, &user.Login)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Авторизация пользователя
func AuthenticateUser(login, password string) (*models.User, error) {
	var user models.User
	err := DB.QueryRow("SELECT id, login FROM users WHERE login = ? AND password = ?", login, password).Scan(&user.ID, &user.Login)
	if err != nil {
		return nil, fmt.Errorf("invalid login or password")
	}
	return &user, nil
}
