package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"string"`
}

type AddTaskResponse struct {
	ID    int    `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
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

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, "Ошибка декодирования json", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		http.Error(w, "Заголовок задачи не может быть пустым", http.StatusBadRequest)
		return
	}

	id, err := insertIntoDB(task)
	if err != nil {
		http.Error(w, "Ошибка добавления в базу данных", http.StatusBadRequest)
		return
	}
	response := AddTaskResponse{ID: id}
	json.NewEncoder(w).Encode(response)
}
