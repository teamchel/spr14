package api

import (
	"net/http"
	"spr14/pkg/db"
)

// TasksResp структура для формирования ответа со списком задач.
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// TasksHandler обрабатывает GET-запросы к /api/tasks для получения списка задач.
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		return
	}

	// Получение списка задач из БД
	tasks, err := db.Tasks(50)
	if err != nil {
		// Ошибка получения списка задач
		writeJSON(w, map[string]string{"error": "Ошибка получения списка задач: " + err.Error()}, http.StatusInternalServerError) // Исправлено
		return
	}

	// Если задач нет, tasks будет пустым слайсом
	if tasks == nil {
		tasks = []*db.Task{} // Создаем пустой слайс, если nil
	}
	// Успешный ответ
	writeJSON(w, TasksResp{Tasks: tasks}, http.StatusOK) // Исправлено
}
