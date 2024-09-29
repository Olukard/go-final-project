package handlers

import (
	"bytes"
	"encoding/json"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request, db *db.DB) {

	var task models.Task
	var buf bytes.Buffer

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка чтения запроса"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response := models.ErrorResponse{Error: "Ошибка декодирования json"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" {
		response := models.ErrorResponse{Error: "Заголовок задачи не может быть пустым"}
		json.NewEncoder(w).Encode(response)
		return
	}

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(TimeFormat)
	} else {
		dueDate, err := time.Parse(TimeFormat, task.Date)
		if err != nil {
			response := models.ErrorResponse{Error: "Неверный формат даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		if dueDate.Before(now) {
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
