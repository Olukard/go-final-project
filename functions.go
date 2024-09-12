package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

//Функция проверки существования файла базы данных

func checkDBexists() bool {
	appPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(appPath)

	dbFile := filepath.Join(filepath.Dir(appPath), "scheduler.db")
	_, err = os.Stat(dbFile)

	fmt.Println(dbFile)

	return err == nil
}

//функция создания базы данных

func CreateDB() {
	DBinitCommand := `
	CREATE TABLE IF NOT EXISTS scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date DATE, 
	title VARCHAR(256),
	comment VARCHAR(256),
	repeat VARCHAR(128)
	);`

	db, err := sql.Open("sqlite3", "./scheduler.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(DBinitCommand)
	if err != nil {
		log.Fatal(err)
	}
}
