package db

import (
	"database/sql"
	"fmt"
	"time"
)

// DateFormat — формат даты для хранения в БД
const DateFormat = "20060102"

// SearchDateFormat — формат даты для поиска
const SearchDateFormat = "02.01.2006"

// Task описывает задачу планировщика
type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// AddTask добавляет задачу в БД и возвращает её ID
func AddTask(task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	var id int64
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

// Tasks возвращает список ближайших задач
func Tasks(search string, limit int) ([]*Task, error) {
	tasks := make([]*Task, 0)

	var rows *sql.Rows
	var err error

	switch {
	case search == "":
		// Без поиска — возвращаем все задачи
		query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`
		rows, err = DB.Query(query, limit)
	default:
		// Проверяем, является ли search датой
		if t, parseErr := time.Parse(SearchDateFormat, search); parseErr == nil {
			// Это дата — ищем по дате
			dateStr := t.Format(DateFormat)
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
			          WHERE date = ? ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, dateStr, limit)
		} else {
			// Это подстрока — ищем в заголовке и комментарии
			searchPattern := "%" + search + "%"
			query := `SELECT id, date, title, comment, repeat FROM scheduler 
			          WHERE (title LIKE ? COLLATE NOCASE OR comment LIKE ? COLLATE NOCASE) 
			          ORDER BY date LIMIT ?`
			rows, err = DB.Query(query, searchPattern, searchPattern, limit)
		}
	}

	if err != nil {
		return tasks, fmt.Errorf("ошибка выборки задач: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := new(Task)
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return tasks, fmt.Errorf("ошибка сканирования строки: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return tasks, fmt.Errorf("ошибка при обходе результатов: %w", err)
	}

	return tasks, nil
}

// GetTask возвращает задачу по её ID
func GetTask(id string) (*Task, error) {
	task := new(Task)
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка получения задачи: %w", err)
	}
	return task, nil
}

// UpdateTask обновляет задачу в БД
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества затронутых строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// DeleteTask удаляет задачу из БД по ID
func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DB.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества удалённых строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}

// UpdateDate обновляет только дату задачи
func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DB.Exec(query, next, id)
	if err != nil {
		return fmt.Errorf("ошибка обновления даты: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("ошибка проверки количества затронутых строк: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}

	return nil
}
