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

	// Получение списка задач из БД (максимум 50)
	tasks, err := db.Tasks(50)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка получения списка задач: " + err.Error()})
		return
	}

	// Если задач нет, tasks будет пустым слайсом, а не nil
	// Это важно для правильной сериализации в JSON
	if tasks == nil {
		tasks = []*db.Task{} // Создаем пустой слайс, если nil
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
