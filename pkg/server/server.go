package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
)

func Run() {
	port := "7540"
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		port = envPort
	}

	webDir := "./web"

	// Регистрируем API-обработчики
	api.Init()

	// Раздаём статические файлы
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Сервер запущен на порту %s", port)
	log.Printf("Раздаём файлы из директории: %s", webDir)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
