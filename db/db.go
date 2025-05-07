package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "embed"

	_ "github.com/mattn/go-sqlite3" // SQLite драйвер
)

var DB *sql.DB

func InitDB() {
	var err error
	// Логирование начала инициализации
	log.Println("Starting database initialization...")

	// Открываем временную базу данных
	DB, err = sql.Open("sqlite3", "./calculator.db")
	if err != nil {
		log.Fatal("Failed to open the embedded database:", err)
	}
	log.Println("Database opened successfully.")

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
		log.Fatal("Failed to create users table:", err)
	}

	// Создаем таблицу задач
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
		log.Fatal("Failed to create tasks table:", err)
	}

	// Применяем миграции для добавления столбца token
	_, err = DB.Exec(`
		PRAGMA foreign_keys=OFF;

		-- Проверяем, есть ли уже столбец token в таблице users
		ALTER TABLE users ADD COLUMN token TEXT;

		PRAGMA foreign_keys=ON;
	`)
	if err != nil {
		// Если столбец уже существует, эта команда вызовет ошибку, которую можно игнорировать
		if err.Error() != "table users has no column named token" {
			fmt.Println("Failed to add token column:", err)
		}
	}
	log.Println("Database initialized successfully")
}
