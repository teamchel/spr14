package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spr14/pkg/db"
	"time"
)

// writeJSON записывает данные в формате JSON в ResponseWriter.
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// checkDate проверяет и корректирует дату задачи.
func checkDate(task *db.Task) error {
	now := time.Now()
	nowStr := now.Format(DateFormat)

	if task.Date == "" {
		task.Date = nowStr
		return nil
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты: %w", err)
	}

	// Проверяем, не прошла ли дата
	if t.Before(now.Truncate(24 * time.Hour)) { // Сравниваем только даты
		if task.Repeat == "" {
			// Если нет повторения, ставим сегодня
			task.Date = nowStr
		} else {
			// Если есть повторение, вычисляем следующую дату
			nextDate, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				return fmt.Errorf("ошибка вычисления следующей даты: %w", err)
			}
			task.Date = nextDate
		}
	}

	return nil
}

// AddTaskHandler обрабатывает POST-запросы к /api/task для добавления задачи.
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		// Ошибка десериализации JSON (400 Bad Request)
		writeJSON(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()}, http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		// Отсутствует обязательное поле (400 Bad Request)
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"}, http.StatusBadRequest)
		return
	}

	// Проверка и коррекция даты
	err = checkDate(&task)
	if err != nil {
		// Ошибка в дате или правиле повторения (400 Bad Request)
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавление задачи в БД
	id, err := db.AddTask(&task)
	if err != nil {
		// Ошибка базы данных (500 Internal Server Error)
		writeJSON(w, map[string]string{"error": "Ошибка добавления задачи в БД: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	// Возврат ID созданной задачи (201 Created)
	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusCreated)
}
