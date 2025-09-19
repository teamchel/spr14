// pkg/api/api.go
package api

import (
	"net/http"
)

// Init пустая функция для совместимости
func Init() {
	// Ничего не делаем
}

// TaskHandler обрабатывает все методы для /api/task
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		PutTaskHandler(w, r)
	case http.MethodGet:
		GetTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
