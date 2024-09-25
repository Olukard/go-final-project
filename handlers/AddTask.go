package handlers

import (
	"bytes"
	"encoding/json"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	//обозначаем структуры
	var task models.Task
	var buf bytes.Buffer

	//устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка чтения запроса"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//декодируем json
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response := models.ErrorResponse{Error: "Ошибка декодирования json"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//проверяем наличие заголовка запроса
	if task.Title == "" {
		response := models.ErrorResponse{Error: "Заголовок задачи не может быть пустым"}
		json.NewEncoder(w).Encode(response)
		return
	}

	now := time.Now()

	// Проверяем дату задачи на пустое значение
	if task.Date == "" {
		task.Date = now.Format(TimeFormat)
	} else {
		// Проверяем парсится ли дата
		dueDate, err := time.Parse(TimeFormat, task.Date)
		if err != nil {
			response := models.ErrorResponse{Error: "Неверный формат даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Проверяем, не раньше ли дата, чем сегодня
		if dueDate.Before(now) {
			// Если дата раньше, вычисляем следующую дату
			if task.Repeat != "" {
				nextDueDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := models.ErrorResponse{Error: "Неверный формат правила повторения"}
					json.NewEncoder(w).Encode(response)
					return
				}
				task.Date = nextDueDate
			} else {
				task.Date = now.Format(TimeFormat)
			}
		}
		task.Date = now.Format(TimeFormat)
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
