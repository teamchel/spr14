// pkg/api/nextdate.go
package api

import (
	"errors"
	"fmt"
	"net/http"
	"spr14/pkg/utils"
	"strconv"
	"strings"
	"time"
)

// NextDayHandler обрабатывает /api/nextdate
func NextDayHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	var err error

	if nowStr == "" {
		now = time.Now()
	} else {
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, "Invalid 'now' date", http.StatusBadRequest)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, "Missing 'date'", http.StatusBadRequest)
		return
	}

	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte(result))
}

// NextDate вычисляет следующую дату по правилу повторения
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	// Парсим начальную дату
	start, err := time.Parse(utils.DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}

	// Разбиваем правило на части
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("invalid repeat rule")
	}

	var date time.Time

	// Определяем тип правила
	switch parts[0] {
	case "d":
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'd'")
		}
		// Парсим интервал
		intervalStr := parts[1]
		interval, err := strconv.Atoi(intervalStr)
		if err != nil {
			return "", errors.New("invalid interval format")
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("interval out of range [1, 400]")
		}

		// Просто добавляем интервал к начальной дате (не ищем после now!)
		date = start.AddDate(0, 0, interval)

	case "y":
		// Добавляем 1 год к начальной дате (не ищем после now!)
		currentYear := start.Year()
		nextYear := currentYear + 1
		month := start.Month()
		day := start.Day()

		// Обработка 29 февраля для високосных годов
		if month == 2 && day == 29 {
			if !isLeapYear(nextYear) {
				day = 28 // Если следующий год не високосный, устанавливаем 28 февраля
			}
		}

		date = time.Date(nextYear, month, day, start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), start.Location())

	case "w":
		if len(parts) < 2 {
			return "", errors.New("missing days for 'w'")
		}
		weekDays, err := parseWeekDays(parts[1])
		if err != nil {
			return "", err
		}

		// Находим следующий день недели от начальной даты (не ищем после now!)
		date = findNextWeekDay(start, weekDays)

	case "m":
		if len(parts) < 2 {
			return "", errors.New("missing days for 'm'")
		}
		var monthDays []int
		var months []int

		monthParts := strings.Split(parts[1], ",")
		for _, p := range monthParts {
			d, err := parseInterval(p)
			if err != nil {
				return "", err
			}
			monthDays = append(monthDays, d)
		}

		// Опционально: вторая последовательность месяцев
		if len(parts) >= 3 {
			for _, p := range strings.Split(parts[2], ",") {
				m, err := parseInterval(p)
				if err != nil {
					return "", err
				}
				months = append(months, m)
			}
		}

		// Находим следующую дату от начальной даты (не ищем после now!)
		date = findNextMonthDay(start, monthDays, months)

	default:
		return "", errors.New("unsupported repeat rule")
	}

	return date.Format(utils.DateFormat), nil
}

// isLeapYear проверяет, является ли год високосным
func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// parseInterval преобразует строку в число
func parseInterval(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty interval")
	}
	if s[0] == '-' {
		if len(s) == 1 {
			return 0, errors.New("invalid interval")
		}
		n, err := strconv.Atoi(s)
		if err != nil || n < -31 || n > 31 {
			return 0, errors.New("invalid day in month")
		}
		return n, nil
	}
	n, err := strconv.Atoi(s)
	if err != nil || n < 1 || n > 400 {
		return 0, errors.New("invalid interval")
	}
	return n, nil
}

// parseWeekDays парсит список дней недели (1–7)
func parseWeekDays(s string) ([]int, error) {
	var days []int
	for _, part := range strings.Split(s, ",") {
		day, err := parseInterval(part)
		if err != nil || day < 1 || day > 7 {
			return nil, errors.New("invalid weekday")
		}
		days = append(days, day)
	}
	return days, nil
}

// findNextWeekDay находит ближайший день недели из списка
func findNextWeekDay(date time.Time, days []int) time.Time {
	next := date
	for {
		next = next.AddDate(0, 0, 1)
		weekday := int(next.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		for _, d := range days {
			if d == weekday {
				return next
			}
		}
	}
}

// findNextMonthDay находит ближайшую дату в месяце
func findNextMonthDay(date time.Time, days []int, months []int) time.Time {
	next := date
	for {
		next = next.AddDate(0, 0, 1)
		_, month, day := next.Date()
		// Проверяем, является ли день одним из заданных
		for _, d := range days {
			if d == day {
				// Проверяем месяц, если указаны месяцы
				if len(months) == 0 || contains(months, int(month)) {
					return next
				}
			}
		}
	}
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
