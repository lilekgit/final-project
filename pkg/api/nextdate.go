package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// DateFormat — формат даты, используемый в проекте
const DateFormat = "20060102"

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	// Валидация: repeat не должен быть пустым
	if repeat == "" {
		return "", fmt.Errorf("правило повторения не указано")
	}

	// Парсим начальную дату
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата dstart: %w", err)
	}

	// Разбиваем правило на части
	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "d":
		return nextDateByDays(date, now, parts)
	case "y":
		return nextDateByYear(date, now)
	case "w":
		return nextDateByWeekday(date, now, parts)
	case "m":
		return nextDateByMonthday(date, now, parts)
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", rule)
	}
}

// nextDateByDays обрабатывает правило "d <число>"
func nextDateByDays(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила d: отсутствует интервал")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("неверный интервал в правиле d: %w", err)
	}

	if interval <= 0 || interval > 400 {
		return "", fmt.Errorf("интервал в правиле d должен быть от 1 до 400")
	}

	// Увеличиваем дату, пока она не станет больше now
	for {
		date = date.AddDate(0, 0, interval)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(DateFormat), nil
}

// nextDateByYear обрабатывает правило "y"
func nextDateByYear(date, now time.Time) (string, error) {
	for {
		date = date.AddDate(1, 0, 0)
		if afterNow(date, now) {
			break
		}
	}

	return date.Format(DateFormat), nil
}

// nextDateByWeekday обрабатывает правило "w <дни_недели>"
func nextDateByWeekday(date, now time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("неверный формат правила w")
	}

	// Парсим дни недели
	weekdayStrs := strings.Split(parts[1], ",")
	weekdays := make(map[time.Weekday]bool)

	for _, ws := range weekdayStrs {
		ws = strings.TrimSpace(ws)
		w, err := strconv.Atoi(ws)
		if err != nil || w < 1 || w > 7 {
			return "", fmt.Errorf("недопустимый день недели: %s", ws)
		}

		// Преобразуем: 1=пн, 7=вс → в Go: 1=пн, 0=вс
		var goWeekday time.Weekday
		if w == 7 {
			goWeekday = time.Sunday
		} else {
			goWeekday = time.Weekday(w)
		}
		weekdays[goWeekday] = true
	}

	// Идём день за днём, начиная со следующего после dstart
	date = date.AddDate(0, 0, 1)
	for {
		if weekdays[date.Weekday()] && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

// nextDateByMonthday обрабатывает правило "m <дни> [месяцы]"
func nextDateByMonthday(date, now time.Time, parts []string) (string, error) {
	if len(parts) < 2 || len(parts) > 3 {
		return "", fmt.Errorf("неверный формат правила m")
	}

	// Парсим дни месяца
	dayStrs := strings.Split(parts[1], ",")
	days := make(map[int]bool)
	for _, ds := range dayStrs {
		ds = strings.TrimSpace(ds)
		d, err := strconv.Atoi(ds)
		if err != nil {
			return "", fmt.Errorf("недопустимый день месяца: %s", ds)
		}

		// Допустимые значения: 1-31, -1, -2
		if (d < 1 || d > 31) && d != -1 && d != -2 {
			return "", fmt.Errorf("недопустимый день месяца: %d", d)
		}
		days[d] = true
	}

	// Парсим месяцы (опционально)
	months := make(map[time.Month]bool)
	if len(parts) == 3 {
		monthStrs := strings.Split(parts[2], ",")
		for _, ms := range monthStrs {
			ms = strings.TrimSpace(ms)
			m, err := strconv.Atoi(ms)
			if err != nil || m < 1 || m > 12 {
				return "", fmt.Errorf("недопустимый месяц: %s", ms)
			}
			months[time.Month(m)] = true
		}
	}

	// Идём день за днём
	date = date.AddDate(0, 0, 1)
	for {
		// Проверяем месяц (если указаны конкретные месяцы)
		if len(months) > 0 && !months[date.Month()] {
			date = date.AddDate(0, 0, 1)
			continue
		}

		// Проверяем день месяца
		dayOfMonth := date.Day()
		lastDayOfMonth := daysInMonth(date.Year(), date.Month())

		matched := false
		for d := range days {
			if d > 0 && d == dayOfMonth {
				matched = true
				break
			}
			if d == -1 && dayOfMonth == lastDayOfMonth {
				matched = true
				break
			}
			if d == -2 && dayOfMonth == lastDayOfMonth-1 {
				matched = true
				break
			}
		}

		if matched && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}

		date = date.AddDate(0, 0, 1)
	}
}

// daysInMonth возвращает количество дней в месяце
func daysInMonth(year int, month time.Month) int {
	// Переходим на первый день следующего месяца и вычитаем один день
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// afterNow возвращает true, если date строго больше now (без учёта времени)
func afterNow(date, now time.Time) bool {
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	nowOnly := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return dateOnly.After(nowOnly)
}

// nextDateHandler — обработчик GET /api/nextdate
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметры
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Если now не указан, берём текущую дату
	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			http.Error(w, "Некорректный параметр now", http.StatusBadRequest)
			return
		}
	}

	// Вызываем функцию NextDate
	result, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Возвращаем результат
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}
