package main

//привести ошибки к единому формату

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

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

func checkDBexists() bool {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	return err == nil
}

//CreateDB создает файл базы данных с индексакцией в соотвествии с заданными константами DBinitCommand и DBindexCommand

func CreateDB() {

	db, err := sql.Open("sqlite3", "./scheduler.db")
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

func insertIntoDB(task Task) (int, error) {
	db, err := sql.Open("sqlite3", "./scheduler.db")
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

func getFromDB() ([]Task, error) {
	var tasks []Task

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, date, title, comment, repeat FROM scheduler")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var task Task

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
		tasks = []Task{}
	}
	return tasks, nil
}
