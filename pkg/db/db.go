package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const (
	schema = `
CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT NOT NULL DEFAULT '',
	title TEXT NOT NULL DEFAULT '',
	comment TEXT DEFAULT '',
	repeat VARCHAR(128) DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);
`
)

var db *sql.DB

func Init(dbPath string) error {
	var err error

	// Если путь не задан, используем значение по умолчанию
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	// Проверяем, существует ли файл базы данных
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("Файл базы данных не найден. Создаём новую.")
	} else if err != nil {
		return fmt.Errorf("ошибка при проверке файла: %w", err)
	}

	// Открываем базу данных
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	// Проверяем соединение
	if err = db.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	// Выполняем SQL-схему (создание таблицы и индекса)
	_, err = db.Exec(schema)
	if err != nil {
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	log.Println("База данных успешно инициализирована")
	return nil
}

// AddTask добавляет задачу в базу данных
func AddTask(task *Task) (int64, error) {
	if db == nil {
		return 0, fmt.Errorf("база данных не инициализирована")
	}

	query := `
		INSERT INTO scheduler (date, title, comment, repeat) 
		VALUES (?, ?, ?, ?)
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// GetTask получает задачу по ID
func GetTask(id string) (*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		WHERE id = ?
	`
	row := db.QueryRow(query, id)
	task := &Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return nil, err
	}
	return task, nil
}

// UpdateTask обновляет задачу
func UpdateTask(task *Task) error {
	if db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	query := `
		UPDATE scheduler 
		SET date = ?, title = ?, comment = ?, repeat = ? 
		WHERE id = ?
	`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления")
	}

	return nil
}

// DeleteTask удаляет задачу
func DeleteTask(id string) error {
	if db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	query := "DELETE FROM scheduler WHERE id = ?"
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет дату задачи
func UpdateDate(next, id string) error {
	if db == nil {
		return fmt.Errorf("база данных не инициализирована")
	}

	query := "UPDATE scheduler SET date = ? WHERE id = ?"
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("некорректный идентификатор для обновления")
	}

	return nil
}

// Tasks возвращает список задач
func Tasks(limit int) ([]*Task, error) {
	if db == nil {
		return nil, fmt.Errorf("база данных не инициализирована")
	}

	query := `
		SELECT id, date, title, comment, repeat 
		FROM scheduler 
		ORDER BY date ASC 
		LIMIT ?
	`
	rows, err := db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}
