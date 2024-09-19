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
	date DATE, 
	title VARCHAR(256),
	comment VARCHAR(256),
	repeat VARCHAR(128)
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

	result, err := db.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Date, task.Repeat)
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
