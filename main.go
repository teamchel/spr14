package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"spr14/pkg/api"
	"spr14/pkg/db"
)

func main() {
	// Шаг 1: Запуск веб-сервера
	port := "7540" // Порт по умолчанию
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	// Шаг 2: Инициализация БД
	dbFile := "scheduler.db"
	if envDBFile := os.Getenv("TODO_DBFILE"); envDBFile != "" {
		dbFile = envDBFile
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Закрытие БД при завершении работы сервера
	defer db.Close()

	// Шаг 3: Инициализация API обработчиков
	api.Init()

	// Обработчик для файлов фронтенда
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
