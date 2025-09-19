package main

import (
	"log"
	"os"

	"spr14/pkg/db"
	"spr14/pkg/server"
)

func main() {
	port := getPort()
	log.Printf("Запуск сервера на порту %s\n", port)

	// Инициализация базы данных
	dbPath := os.Getenv("TODO_DBFILE")
	if dbPath == "" {
		dbPath = "scheduler.db"
	}

	if err := db.Init(dbPath); err != nil {
		log.Fatal("Ошибка инициализации базы данных:", err)
	}

	if err := server.Start(port); err != nil {
		log.Fatal(err)
	}
}

func getPort() string {
	if port := os.Getenv("TODO_PORT"); port != "" {
		return port
	}
	return "7540"
}
