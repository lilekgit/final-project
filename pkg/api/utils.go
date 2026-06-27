package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

// writeJson записывает данные в ответ в формате JSON
func writeJson(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}

// writeError записывает ошибку в ответ в формате JSON
func writeError(w http.ResponseWriter, message string, code int) {
	w.WriteHeader(code)
	writeJson(w, map[string]string{"error": message})
}

// checkDate проверяет и корректирует дату задачи
func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата не указана — берём сегодняшнее число
	if task.Date == "" {
		task.Date = now.Format(db.DateFormat)
	}

	// Парсим дату
	t, err := time.Parse(db.DateFormat, task.Date)
	if err != nil {
		return fmt.Errorf("некорректный формат даты")
	}

	var next string
	// Если правило повторения указано — проверяем его
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("некорректное правило повторения: %v", err)
		}
	}

	// Если дата в прошлом
	if afterNow(now, t) {
		if task.Repeat == "" {
			// Нет правила — берём сегодня
			task.Date = now.Format(db.DateFormat)
		} else {
			// Есть правило — берём вычисленную следующую дату
			task.Date = next
		}
	}

	return nil
}
