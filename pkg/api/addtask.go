// pkg/api/addtask.go
package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spr14/pkg/db"
	"spr14/pkg/utils"
	"time"
)

// writeJSON отправляет данные в формате JSON
func writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// AddTaskHandler обрабатывает POST /api/task
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"})
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка парсинга JSON"})
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Проверка даты
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(utils.DateFormat)
	} else {
		t, err := time.Parse(utils.DateFormat, task.Date)
		if err != nil {
			writeJSON(w, map[string]string{"error": "некорректный формат даты"})
			return
		}
		// Если дата меньше сегодняшней - устанавливаем сегодня
		if t.Before(truncateToDay(now)) {
			task.Date = now.Format(utils.DateFormat)
		}
	}

	// Если есть правило повторения - вычисляем следующую дату
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "неподдерживаемое правило повторения"})
			return
		}
		task.Date = next
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка сохранения задачи"})
		return
	}

	writeJSON(w, map[string]string{"id": fmt.Sprintf("%d", id)})
}

// GetTaskHandler обрабатывает GET /api/task
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

// PutTaskHandler обрабатывает PUT /api/task
func PutTaskHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка чтения запроса"})
		return
	}

	var task db.Task
	err = json.Unmarshal(body, &task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "ошибка парсинга JSON"})
		return
	}

	// Проверка обязательных полей
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "не указан заголовок задачи"})
		return
	}

	// Проверка даты
	now := time.Now()
	if task.Date == "" {
		task.Date = now.Format(utils.DateFormat)
	} else {
		t, err := time.Parse(utils.DateFormat, task.Date)
		if err != nil {
			writeJSON(w, map[string]string{"error": "некорректный формат даты"})
			return
		}
		// Если дата меньше сегодняшней - устанавливаем сегодня
		if t.Before(truncateToDay(now)) {
			task.Date = now.Format(utils.DateFormat)
		}
	}

	// Если есть правило повторения - вычисляем следующую дату
	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "неподдерживаемое правило повторения"})
			return
		}
		task.Date = next
	}

	// Обновляем задачу
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешное обновление — пустой ответ
	writeJSON(w, map[string]string{})
}

// DoneTaskHandler обрабатывает POST /api/task/done
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	now := time.Now()

	if task.Repeat == "" {
		// Одноразовая задача — удаляем
		err = db.DeleteTask(id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	} else {
		// Периодическая задача — вычисляем следующую дату
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}

		err = db.UpdateDate(next, id)
		if err != nil {
			writeJSON(w, map[string]string{"error": err.Error()})
			return
		}
	}

	// Успешно — пустой ответ
	writeJSON(w, map[string]string{})
}

// DeleteTaskHandler обрабатывает DELETE /api/task/delete
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Успешное удаление — пустой ответ
	writeJSON(w, map[string]string{})
}

// truncateToDay обрезает время до начала дня
func truncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
