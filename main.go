package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Проверяем наличие базы данных...")
	if !checkDBexists() {
		fmt.Println("База данных не найдена, создаем...")
		CreateDB()
	} else {
		fmt.Println("База данных найдена.")
	}

	router := mux.NewRouter()

	// Обработчики
	router.HandleFunc("/api/nextdate", nextDateHandler).Methods("GET")
	router.HandleFunc("/api/task", addTaskHandler).Methods("POST")
	router.HandleFunc("/api/tasks", getTaskHandler).Methods("GET")

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
