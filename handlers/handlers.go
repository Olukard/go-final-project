package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-final-project/db"
	"go-final-project/models"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
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

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	//обозначаем структуры
	var task models.Task
	var buf bytes.Buffer

	//устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response := models.AddTaskResponse{Error: "Ошибка чтения запроса"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//декодируем json
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response := models.AddTaskResponse{Error: "Ошибка декодирования json"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//проверяем наличие заголовка запроса
	if task.Title == "" {
		response := models.AddTaskResponse{Error: "Заголовок задачи не может быть пустым"}
		json.NewEncoder(w).Encode(response)
		return
	}

	now := time.Now()

	// Проверяем дату задачи на пустое значение
	if task.Date == "" {
		task.Date = now.Format("20060102")
	} else {
		// Проверяем парсится ли дата
		dueDate, err := time.Parse("20060102", task.Date)
		if err != nil {
			response := models.AddTaskResponse{Error: "Неверный формат даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Проверяем, не раньше ли дата, чем сегодня
		if dueDate.Before(now) {
			// Если дата раньше, вычисляем следующую дату
			if task.Repeat != "" {
				nextDueDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := models.AddTaskResponse{Error: "Неверный формат правила повторения"}
					json.NewEncoder(w).Encode(response)
					return
				}
				task.Date = nextDueDate
			} else {
				task.Date = now.Format("20060102")
			}
		}
	}

	id, err := db.InsertIntoDB(task)
	if err != nil {
		response := models.AddTaskResponse{Error: "Ошибка добавления в базу данных"}
		json.NewEncoder(w).Encode(response)
		return

	}
	response := models.AddTaskResponse{ID: id}
	json.NewEncoder(w).Encode(response)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks, err := db.GetFromDB()
	if err != nil {
		response := models.AddTaskResponse{Error: "Ошибка получения данных"}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
	})
}
