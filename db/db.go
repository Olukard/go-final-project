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

	log.Println("Checking db file...")

	dbfile := os.Getenv("TODO_DBFILE")
	if dbfile == "" {
		dbfile = "scheduler.db"
	}

	_, err := os.Stat(dbfile)

	if err == nil {
		log.Println("Database found.")
		db, err := sql.Open("sqlite3", "./"+dbfile)
		return &DB{db: db}, err
	}

	log.Println("Database not found, creating database.")

	db, err := sql.Open("sqlite3", "./"+dbfile)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	_, err = db.Exec(DBinitCommand)
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	_, err = db.Exec(DBindexCommand)
	if err != nil {
		return nil, fmt.Errorf("error indexing database: %w", err)
	}

	return &DB{db: db}, nil
}

// InsertIntoDB добавляет запись задачи в базу данных
func (d *DB) InsertIntoDB(task models.Task) (int, error) {

	result, err := d.db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, fmt.Errorf("error executing db command: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error returning task: %w", err)
	}

	return int(id), nil

}

// GetTaskFromDB возвращает конкретную задачу из базы данных по её id
func (d *DB) GetTaskFromDB(id string) (task models.Task, err error) {

	err = d.db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, fmt.Errorf("error querying database: %w", err)
	}

	return task, nil
}

// UpdateTaskInDB редактирует запись задачи в базе данных
func (d *DB) UpdateTaskInDB(task models.Task) (err error) {

	_, err = d.db.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return fmt.Errorf("error updating task in db: %w", err)
	}
	return nil
}

// GetListFromDB получает список задач из базы данных
func (d *DB) GetListFromDB() (tasks []models.Task, err error) {

	rows, err := d.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 50")
	if err != nil {
		return nil, fmt.Errorf("error querying database: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("error scanning db rows: %w", err)
		}
		task.ID = fmt.Sprint(id)

		tasks = append(tasks, task)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error scanning db rows: %w", err)
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
		return fmt.Errorf("error deleting from database: %w", err)
	}
	return nil
}
