package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

// getTaskHandler обрабатывает GET-запрос /api/task?id=<id>
func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, task)
}

// updateTaskHandler обрабатывает PUT-запрос /api/task
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeError(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeError(w, "не указан идентификатор задачи", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, map[string]any{})
}

// deleteTaskHandler обрабатывает DELETE-запрос /api/task?id=<id>
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJson(w, map[string]any{})
}

// doneTaskHandler обрабатывает POST-запрос /api/task/done?id=<id>
func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		writeError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	// Получаем задачу
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	// Очищаем repeat от пробелов
	task.Repeat = strings.TrimSpace(task.Repeat)

	// Если нет правила повторения — удаляем задачу
	if task.Repeat == "" {
		if err := db.DeleteTask(id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		// Есть правило повторения — вычисляем следующую дату
		now := time.Now()
		nextDate, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, "ошибка вычисления следующей даты: "+err.Error(), http.StatusBadRequest)
			return
		}

		// Обновляем дату задачи
		if err := db.UpdateDate(nextDate, id); err != nil {
			writeError(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJson(w, map[string]any{})
}
