package api

import (
	"go_final_project/pkg/db"
	"net/http"
)

// DefaultTaskLimit — константа для ограничения количества задач
const DefaultTaskLimit = 50

// TasksHandler — обработчик получения списка задач
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	search := r.URL.Query().Get("search")

	tasks, err := db.Tasks(DefaultTaskLimit, search)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if tasks == nil {
		tasks = []db.Task{}
	}

	writeJSON(w, http.StatusOK, map[string]any{"tasks": tasks})
}
