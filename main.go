package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	// Обработчик для /api/nextdate
	router.HandleFunc("/api/nextdate", nextDateHandler).Methods("GET")

	// Обработчик для статических файлов (из директории "web")
	fileServer := http.FileServer(http.Dir("./web"))
	router.PathPrefix("/").Handler(fileServer) // Обратите внимание на router

	// Добавить выбор порта
	port := ":7540"

	fmt.Println("Проверяем наличие базы данных...")
	if !checkDBexists() {
		fmt.Println("База данных не найдена, создаем...")
		CreateDB()
	} else {
		fmt.Println("База данных найдена.")
	}

	fmt.Printf("Запускаем сервер. Порт%s\n", port)
	err := http.ListenAndServe(port, router) // Запускаем сервер с router
	if err != nil {
		panic(err)
	}

	fmt.Println("Завершаем работу")
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Проверяем, что параметры не пусты
	if nowStr == "" {
		http.Error(w, "Не задана дата", http.StatusBadRequest)
	}
	if dateStr == "" {
		http.Error(w, "Не задана дата", http.StatusBadRequest)
		return
	}
	if repeatStr == "" {
		http.Error(w, "Не задано правило повторения", http.StatusBadRequest)
		return
	}

	// Парсим дату
	nowDate, err := time.Parse("20060102", nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты1", http.StatusBadRequest)
		return
	}

	repeatDate, err := NextDate(nowDate, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, repeatDate)
}
