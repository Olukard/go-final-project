package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"go-final-project/db"
	"go-final-project/handlers"
)

var DB *sql.DB

func main() {

	fmt.Println("Проверяем наличие базы данных...")
	if !db.CheckDBexists() {
		fmt.Println("База данных не найдена, создаем...")
		db.CreateDB()
	} else {
		fmt.Println("База данных найдена.")
	}

	router := mux.NewRouter()

	// Обработчики
	router.HandleFunc("/api/nextdate", handlers.NextDateHandler).Methods("GET")
	router.HandleFunc("/api/task", handlers.AddTaskHandler).Methods("POST")
	router.HandleFunc("/api/task", handlers.GetTaskHandler).Methods("GET")
	router.HandleFunc("/api/task", handlers.EditTaskHandler).Methods("PUT")
	router.HandleFunc("/api/tasks", handlers.GetTasksListHandler).Methods("GET")

	// Обработчик для статических файлов (из директории "web")
	fileServer := http.FileServer(http.Dir("./web"))
	router.PathPrefix("/").Handler(fileServer) // Обратите внимание на router

	port := ":7540"

	fmt.Printf("Запускаем сервер. Порт%s\n", port)
	err := http.ListenAndServe(port, router) // Запускаем сервер с router
	if err != nil {
		panic(err)
	}

	fmt.Println("Завершаем работу")
}
