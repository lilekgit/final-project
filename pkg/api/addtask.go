package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"go_final_project/pkg/db"
)

// addTaskHandler обрабатывает POST-запрос на добавление задачи
func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Десериализуем JSON
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	// Проверяем заголовок
	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	// Проверяем и корректируем дату
	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Добавляем задачу в БД
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем ID в виде строки
	writeJson(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}
