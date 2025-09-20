package api

import (
	"net/http"
)

// Init регистрирует все API обработчики.
func Init() {
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			AddTaskHandler(w, r)
		case http.MethodGet:
			GetTaskHandler(w, r)
		case http.MethodPut:
			UpdateTaskHandler(w, r)
		case http.MethodDelete:
			DeleteTaskHandler(w, r)
		default:
			http.Error(w, `{"error":"Метод не поддерживается"}`, http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/tasks", TasksHandler)
	http.HandleFunc("/api/task/done", DoneTaskHandler)
}
