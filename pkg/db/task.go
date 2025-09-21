package db

import (
	"database/sql"
	"fmt"
)

// Task представляет собой структуру задачи.
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет новую задачу в базу данных и возвращает её ID.
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("не удалось вставить задачу: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось получить ID вставленной задачи: %w", err)
	}

	return id, nil
}

// Tasks получает список задач из базы данных, отсортированных по дате.
func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
	rows, err := GetDB().Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		task := &Task{}
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		tasks = append(tasks, task)
	}

	// Проверяем на ошибки после итерации
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по строкам: %w", err)
	}

	return tasks, nil
}

// GetTask получает задачу из базы данных по её ID.
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := GetDB().QueryRow(query, id)

	task := &Task{}
	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача с ID %s не найдена", id)
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}

	return task, nil
}

// UpdateTask обновляет задачу в базе данных.
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := GetDB().Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества изменённых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с ID %s не найдена для обновления", task.ID)
	}

	return nil
}

// DeleteTask удаляет задачу из базы данных по её ID.
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := GetDB().Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества удалённых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с ID %s не найдена для удаления", id)
	}

	return nil
}

// UpdateTaskDate обновляет только дату задачи в базе данных.
func UpdateTaskDate(id, newDate string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := GetDB().Exec(query, newDate, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка получения количества изменённых строк: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("задача с ID %s не найдена для обновления даты", id)
	}

	return nil
}
