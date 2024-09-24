package db

//привести ошибки к единому формату

import (
	"database/sql"
	"fmt"
	"go-final-project/models"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const DBfile = "scheduler.db"

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

//checkDBexists проверяет существование файла базы данных в директории проекта

func CheckDBexists() bool {

	_, err := os.Stat(DBfile)

	return err == nil
}

//CreateDB создает файл базы данных с индексакцией в соотвествии с заданными константами DBinitCommand и DBindexCommand

func CreateDB() {

	db, err := sql.Open("sqlite3", "./"+DBfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(DBinitCommand)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("База данных создана, проводим индексацию...")

	_, err = db.Exec(DBindexCommand)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Индексация завершена.")

}

func InsertIntoDB(task models.Task) (int, error) {

	db, err := sql.Open("sqlite3", "./"+DBfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	return int(id), nil

}

func GetTaskFromDB(id string) (task models.Task, err error) {
	db, err := sql.Open("sqlite3", "./"+DBfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.QueryRow("SELECT * FROM scheduler WHERE id = ?", id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		return task, err
	}

	return task, nil

}

func GetListFromDB() (tasks []models.Task, err error) {
	db, err := sql.Open("sqlite3", "./"+DBfile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task

		var id int64

		err := rows.Scan(&id, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		task.ID = fmt.Sprint(id)

		tasks = append(tasks, task)

		if err := rows.Err(); err != nil {
			log.Println(err)
			return nil, err
		}

	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	return tasks, nil
}
