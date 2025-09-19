// pkg/api/tasks.go
package api

import (
	"net/http"
	"spr14/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// TasksHandler обрабатывает GET /api/tasks
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	const limit = 50 // максимальное количество задач

	tasks, err := db.Tasks(limit)
	if err != nil {
		// Используем правильную функцию writeJSON
		writeJSON(w, map[string]interface{}{"error": "ошибка получения задач"})
		return
	}

	// Отправляем ответ с задачами
	resp := TasksResp{Tasks: tasks}
	writeJSON(w, resp)
}
