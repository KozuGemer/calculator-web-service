// db.go
package db

import (
	"database/sql"
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
			password TEXT
		)
	`)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Database initialized successfully")
}
