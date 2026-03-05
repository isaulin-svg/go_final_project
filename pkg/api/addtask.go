package api

import (
	"encoding/json"
	"fmt"
	"go_final_project/pkg/db"
	"net/http"
	"strings"
	"time"
)

func validateAndFixDate(dateStr, repeat string) (string, error) {
	now := time.Now()
	todayStr := now.Format(DateLayout)

	if strings.TrimSpace(dateStr) == "" {
		return todayStr, nil
	}

	t, err := time.Parse(DateLayout, dateStr)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты")
	}

	today, _ := time.Parse(DateLayout, todayStr)

	if t.Before(today) {
		if strings.TrimSpace(repeat) == "" {
			return todayStr, nil
		}
		next, err := NextDate(now, dateStr, repeat)
		if err != nil {
			return "", err
		}
		return next, nil
	}
	return dateStr, nil
}

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		task, err := db.GetTask(id)
		if err != nil {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
			return
		}
		writeJSON(w, http.StatusOK, task)

	case http.MethodPost:
		AddTaskHandler(w, r)

	case http.MethodPut:
		var task db.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка JSON"})
			return
		}

		// ВАЖНО: Проверка заголовка при редактировании (Шаг 6)
		if strings.TrimSpace(task.Title) == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
			return
		}

		fixedDate, err := validateAndFixDate(task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		task.Date = fixedDate

		if err := db.UpdateTask(task); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})

	case http.MethodDelete:
		id := r.URL.Query().Get("id")
		if id == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
			return
		}
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
	}
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка JSON"})
		return
	}

	if strings.TrimSpace(task.Title) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	fixedDate, err := validateAndFixDate(task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	task.Date = fixedDate

	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": fmt.Sprintf("%d", id)})
}

func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		err = db.DeleteTask(id)
	} else {
		nextDate, errNext := NextDate(time.Now(), task.Date, task.Repeat)
		if errNext != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": errNext.Error()})
			return
		}
		err = db.UpdateTaskDate(id, nextDate)
	}

	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
