// db.go
package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite драйвер
)

var DB *sql.DB

// Инициализация базы данных
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./calculator.db") // Создаем или открываем базу данных
	if err != nil {
		log.Fatal(err)
	}

	// Создаем таблицу пользователей, если она еще не существует
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			login TEXT UNIQUE,
			password TEXT,
			token TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}
	_, err = DB.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			task_id INTEGER,
			user_id INTEGER,
			expression TEXT,
			result TEXT,
			status TEXT
		)
	`)
	if err != nil {
		fmt.Println(err)
	}
	_, err = DB.Exec(`
		PRAGMA foreign_keys=OFF;

		-- Проверяем, есть ли уже столбец token в таблице users
		ALTER TABLE users ADD COLUMN token TEXT;
		
		PRAGMA foreign_keys=ON;
	`)
	if err != nil {
		// Если столбец уже существует, эта команда вызовет ошибку, которую можно игнорировать
		if err.Error() != "table users has no column named token" {
			return
		}
	}
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized successfully")
}
