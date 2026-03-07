package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {
	// 1. Инициализация БД
	if err := db.InitDB(); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	// Важно: закрываем БД при завершении программы
	defer db.DB.Close()

	// 2. Инициализация конфигурации
	api.AppPassword = os.Getenv("TODO_PASSWORD")
	fmt.Printf("DEBUG: Loaded password as '%s'\n", api.AppPassword)
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	// 3. Настройка маршрутов
	mux := http.NewServeMux()

	// Статика
	webDir := "./web"
	mux.Handle("/", http.FileServer(http.Dir(webDir)))

	// Публичные эндпоинты
	mux.HandleFunc("/api/signin", api.SigninHandler)
	mux.HandleFunc("/api/nextdate", api.NextDateHandler)

	// Защищенные эндпоинты
	mux.HandleFunc("/api/task", api.AuthMiddleware(api.TaskHandler))
	mux.HandleFunc("/api/tasks", api.AuthMiddleware(api.TasksHandler))
	mux.HandleFunc("/api/task/done", api.AuthMiddleware(api.DoneTaskHandler))

	// 4. Запуск
	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
