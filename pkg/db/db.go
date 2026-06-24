package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

// schema — SQL-команды для создания таблицы и индекса
const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title VARCHAR(256) NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

// DB — глобальная переменная для хранения соединения с БД
var DB *sql.DB

// Init открывает базу данных и при необходимости создаёт таблицу с индексом
func Init(dbFile string) error {
	// Проверяем, существует ли файл БД
	_, err := os.Stat(dbFile)
	install := false
	if err != nil {
		if os.IsNotExist(err) {
			install = true
		} else {
			return fmt.Errorf("ошибка проверки файла БД: %w", err)
		}
	}

	// Открываем БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Проверяем подключение
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Если файла не было — создаём таблицу и индекс
	if install {
		fmt.Println("📦 Файл БД не найден, создаём таблицу scheduler...")
		_, err = DB.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка создания таблицы: %w", err)
		}
		fmt.Println("✅ Таблица scheduler успешно создана")
	} else {
		fmt.Println("✅ БД найдена:", dbFile)
	}

	return nil
}
