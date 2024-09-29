package db

import (
	"database/sql"
	"fmt"
	"go-final-project/models"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

const DBinitCommand = `
	CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date VARCHAR(256) NOT NULL DEFAULT "", 
	title VARCHAR(256) NOT NULL DEFAULT "",
	comment VARCHAR(256) NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
	);`

const DBindexCommand = `
	CREATE INDEX id_indx ON scheduler (date)
	`

// CreateDB создает файл базы данных с индексакцией в соотвествии с заданными константами DBinitCommand и DBindexCommand
func CreateDB() (*DB, error) {

	log.Println("Проверяем наличие базы данных...")

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "scheduler.db"
	}

	_, err := os.Stat(dbfile)

	if err == nil {
		log.Println("База данных найдена.")
		db, err := sql.Open("sqlite3", "./"+dbfile)
		return &DB{db: db}, fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	log.Println("База данных не найдена, создаем.")

	db, err := sql.Open("sqlite3", "./"+dbfile)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	_, err = db.Exec(DBinitCommand)
	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации базы данных: %w", err)
	}

	_, err = db.Exec(DBindexCommand)
	if err != nil {
		return nil, fmt.Errorf("ошибка индексации базы данных: %w", err)
	}

	return &DB{db: db}, nil
}

// InsertIntoDB добавляет запись задачи в базу данных
func (d *DB) InsertIntoDB(task models.Task) (int, error) {

	result, err := d.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("ошибка добавления данных: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка возврата задачи: %w", err)
	}

	return int(id), nil

}

// GetTaskFromDB возвращает конкретную задачу из базы данных по её id
func (d *DB) GetTaskFromDB(id string) (task models.Task, err error) {

	err = d.db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, fmt.Errorf("ошибка получения данных из базы данных: %w", err)
	}

	return task, nil
}

// UpdateTaskInDB редактирует запись задачи в базе данных
func (d *DB) UpdateTaskInDB(task models.Task) (err error) {

	_, err = d.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("ошибка изменения задачи в базе данных: %w", err)
	}
	return nil
}

// GetListFromDB получает список задач из базы данных
func (d *DB) GetListFromDB() (tasks []models.Task, err error) {

	rows, err := d.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, fmt.Errorf("ошибка получения списка зада из базы данных: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка считывания данных: %w", err)
		}
		task.ID = fmt.Sprint(id)

		tasks = append(tasks, task)

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("ошибка считывания данных: %w", err)
		}

	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}

// DeleteTaskFromDB удаляет задачу из базы данных
func (d *DB) DeleteTaskFromDB(id string) (err error) {

	_, err = d.db.Exec("DELETE FROM scheduler WHERE id = ?",
		id)
	if err != nil {
		return fmt.Errorf("ошибка удаления из базы данных: %w", err)
	}
	return nil
}
