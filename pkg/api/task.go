package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"spr14/pkg/db"
	"strconv"
)

// GetTaskHandler обрабатывает GET-запросы к /api/task для получения задачи по ID.
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Проверка, что ID является числом
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Некорректный идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, task)
}

// UpdateTaskHandler обрабатывает PUT-запросы к /api/task для обновления задачи.
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	var task db.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка десериализации JSON: " + err.Error()})
		return
	}

	// Проверка обязательных полей
	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор задачи"})
		return
	}
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверка и коррекция даты (аналогично добавлению)
	err = checkDate(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Обновление задачи в БД
	err = db.UpdateTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возврат пустого JSON в случае успеха
	fmt.Fprint(w, "{}")
}

// DeleteTaskHandler обрабатывает DELETE-запросы к /api/task для удаления задачи.
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	// Проверка, что ID является числом
	_, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Некорректный идентификатор"})
		return
	}

	// Удаление задачи из БД
	err = db.DeleteTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Возврат пустого JSON в случае успеха
	fmt.Fprint(w, "{}")
}
