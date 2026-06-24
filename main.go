package main

import (
	"log"
	"os"

	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {
	// Определяем путь к файлу БД
	dbFile := "scheduler.db"
	if envDB := os.Getenv("TODO_DBFILE"); envDB != "" {
		dbFile = envDB
	}

	// Инициализируем БД
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Запускаем сервер
	server.Run()
}
