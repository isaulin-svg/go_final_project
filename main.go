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

	// 2. Раздача фронтенда
	webDir := "./web"
	fs := http.FileServer(http.Dir(webDir))
	http.Handle("/", fs)

	// 3. Публичные эндпоинты (БЕЗ AuthMiddleware)
	// Эти ручки должны быть открыты для прохождения базовых тестов
	http.HandleFunc("/api/signin", api.SigninHandler)
	http.HandleFunc("/api/nextdate", api.NextDateHandler)

	// 4. Защищенные эндпоинты (С AuthMiddleware)
	// Используем http.HandleFunc и оборачиваем саму функцию обработчика
	http.HandleFunc("/api/task", api.AuthMiddleware(api.TaskHandler))
	http.HandleFunc("/api/tasks", api.AuthMiddleware(api.TasksHandler))
	http.HandleFunc("/api/task/done", api.AuthMiddleware(api.DoneTaskHandler))

	// 5. Определение порта и запуск
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	fmt.Printf("Сервер запущен на http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
