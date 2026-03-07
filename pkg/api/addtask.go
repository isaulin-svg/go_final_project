package api

import (
	"encoding/json"
	"errors"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

// TaskHandler — диспетчер для работы с задачами
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
	}
}

// DoneTaskHandler — обработка завершения задачи
func DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	} else {
		newDate, err := NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		task.Date = newDate
		if err := db.UpdateTask(task); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

// Вспомогательные хендлеры
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка JSON"})
		return
	}
	if err := validateTask(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	id, err := db.AddTask(task)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Ошибка JSON"})
		return
	}
	if err := validateTask(&task); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := db.UpdateTask(task); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Задача не найдена"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{})
}

// validateTask — валидация данных
func validateTask(task *db.Task) error {
	if task.Title == "" {
		return errors.New("заголовок не может быть пустым")
	}
	now := time.Now().Truncate(24 * time.Hour)
	if task.Date != "" {
		t, err := time.Parse("20060102", task.Date)
		if err != nil {
			return errors.New("неверный формат даты")
		}
		if t.Before(now) {
			task.Date = time.Now().Format("20060102")
		}
	} else {
		task.Date = time.Now().Format("20060102")
	}
	if task.Repeat != "" {
		if _, err := NextDate(now, task.Date, task.Repeat); err != nil {
			return errors.New("неверный формат правила повторения")
		}
	}
	return nil
}
