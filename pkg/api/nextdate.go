package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DateLayout = "20060102"

// NextDateHandler — проверяет параметры и возвращает следующую дату
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Метод не поддерживается"})
		return
	}

	nowStr := r.URL.Query().Get("now")
	dstart := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse(DateLayout, nowStr)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "неверный формат параметра now"})
			return
		}
	}

	nextDate, err := NextDate(now, dstart, repeat)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	// Здесь мы пишем чистую строку, так как тест ожидает просто текст, а не JSON
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}

// NextDate — бизнес-логика расчета даты
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	start, err := time.Parse(DateLayout, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат начальной даты: %w", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("не указан интервал в днях")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days < 1 || days > 400 {
			return "", errors.New("недопустимый интервал в днях")
		}

		next := start
		for {
			next = next.AddDate(0, 0, days)
			if next.After(now) {
				return next.Format(DateLayout), nil
			}
		}

	case "y":
		next := start
		for {
			next = next.AddDate(1, 0, 0)
			if next.After(now) {
				return next.Format(DateLayout), nil
			}
		}

	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}
