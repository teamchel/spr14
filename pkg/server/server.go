// pkg/server/server.go
package server

import (
	"log"
	"net/http"

	"spr14/pkg/api"
)

const webDir = "./web"

func Start(port string) error {
	// Регистрация обработчиков API
	http.HandleFunc("/api/nextdate", api.NextDayHandler)
	http.HandleFunc("/api/task", api.TaskHandler)
	http.HandleFunc("/api/tasks", api.TasksHandler)
	http.HandleFunc("/api/task/done", api.DoneTaskHandler)
	http.HandleFunc("/api/task/delete", api.DeleteTaskHandler)
	http.HandleFunc("/api/signin", api.SigninHandler)

	// Файловый сервер для статики
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	log.Printf("Сервер запущен на порту %s\n", port)
	return http.ListenAndServe(":"+port, nil)
}
