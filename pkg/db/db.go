package db

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT,
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

// Init инициализирует подключение к базе данных и при необходимости создает таблицу.
func Init(dbFile string) error {
	// Проверяем, существует ли файл базы данных
	_, err := os.Stat(dbFile)
	install := err != nil

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть базу данных: %w", err)
	}

	// Если файла не было, создаем таблицу и индекс
	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("не удалось создать таблицу scheduler: %w", err)
		}
	}

	return nil
}

// GetDB возвращает экземпляр подключения к базе данных.
// Используется другими функциями пакета db.
func GetDB() *sql.DB {
	return db
}
