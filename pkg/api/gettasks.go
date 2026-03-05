package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

// TasksHandler возвращает список задач с поддержкой поиска
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")

	// Устанавливаем лимит 50 задач согласно ТЗ
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	// Фронтенд ожидает объект с ключом "tasks"
	if tasks == nil {
		tasks = []db.Task{}
	}

	writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}
