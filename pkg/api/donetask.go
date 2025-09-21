package api

import (
	"fmt"
	"net/http"
	"spr14/pkg/db"
	"strconv"
	"time"
)

// DoneTaskHandler обрабатывает POST-запросы к /api/task/done для отметки задачи как выполненной.
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"}, http.StatusBadRequest)
		return
	}

	// Проверка, что ID является числом
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Некорректный идентификатор"}, http.StatusBadRequest)
		return
	}

	// Получение задачи из БД
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()}, http.StatusNotFound)
		return
	}

	// Если правило повторения не указано, удаляем задачу
	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			// Ошибка при удалении
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
	} else {
		// Если есть правило повторения, вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, map[string]string{"error": "Ошибка вычисления следующей даты: " + err.Error()}, http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		err = db.UpdateTaskDate(id, nextDate)
		if err != nil {
			// Ошибка при обновлении даты
			writeJSON(w, map[string]string{"error": err.Error()}, http.StatusInternalServerError)
			return
		}
	}

	// Возврат пустого JSON в случае успеха
	fmt.Fprint(w, "{}")
}
