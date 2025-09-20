package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// afterNow проверяет, что дата date больше даты now.
func afterNow(date, now time.Time) bool {
	// Сравниваем только даты, игнорируя время
	y1, m1, d1 := date.Date()
	y2, m2, d2 := now.Date()
	return time.Date(y1, m1, d1, 0, 0, 0, 0, time.UTC).After(time.Date(y2, m2, d2, 0, 0, 0, 0, time.UTC))
}

// NextDate вычисляет следующую дату для задачи по правилу повторения.
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	startDate, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("неверный формат даты начала: %w", err)
	}

	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", fmt.Errorf("некорректный формат правила повторения")
	}

	switch parts[0] {
	case "d":
		if len(parts) != 2 {
			return "", fmt.Errorf("некорректный формат правила 'd': ожидается 'd <число>'")
		}
		interval, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", fmt.Errorf("некорректное значение интервала для 'd': %w", err)
		}
		if interval <= 0 || interval > 400 {
			return "", fmt.Errorf("интервал для 'd' должен быть от 1 до 400")
		}
		date := startDate
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	case "y":
		if len(parts) != 1 {
			return "", fmt.Errorf("правило 'y' не требует дополнительных параметров")
		}
		date := startDate
		for {
			date = date.AddDate(1, 0, 0) // Добавляем один год
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(DateFormat), nil

	default:
		// Базовая реализация: возвращаем ошибку для неподдерживаемых правил
		return "", fmt.Errorf("неподдерживаемое правило повторения: %s", parts[0])
	}
}

// NextDateHandler обрабатывает GET-запросы к /api/nextdate.
func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Устанавливаем правильный Content-Type
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	nowParam := r.FormValue("now")
	dateParam := r.FormValue("date")
	repeatParam := r.FormValue("repeat")

	if dateParam == "" {
		http.Error(w, "Не указана дата начала", http.StatusBadRequest)
		return
	}
	if repeatParam == "" {
		http.Error(w, "Не указано правило повторения", http.StatusBadRequest)
		return
	}

	now := time.Now()
	if nowParam != "" {
		parsedNow, err := time.Parse(DateFormat, nowParam)
		if err != nil {
			http.Error(w, fmt.Sprintf("Неверный формат даты now: %v", err), http.StatusBadRequest)
			return
		}
		// Устанавливаем время now на начало дня переданной даты для корректного сравнения
		now = time.Date(parsedNow.Year(), parsedNow.Month(), parsedNow.Day(), 0, 0, 0, 0, parsedNow.Location())
	}

	nextDate, err := NextDate(now, dateParam, repeatParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем только саму дату, без обертки в JSON
	fmt.Fprint(w, nextDate)
}
