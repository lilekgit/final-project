package api

import (
	"net/http"

	"go_final_project/pkg/db"
)

// TasksResp — структура ответа для списка задач
type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// tasksHandler обрабатывает GET-запрос /api/tasks
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем параметр search из строки запроса (задание со звёздочкой)
	search := r.FormValue("search")

	// Получаем список задач (максимум 50)
	tasks, err := db.Tasks(search, 50)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ в JSON
	writeJson(w, TasksResp{Tasks: tasks})
}
