package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
)

func EditTaskHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		_, err = ValidateDate(task.Date, TimeFormat)
		if err != nil {
			response := models.ErrorResponse{Error: "неверный формат даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		err = ValdateRepeatRule(task.Repeat)
		if err != nil {
			response := models.ErrorResponse{Error: "неверное правило повторения"}
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err = db.GetTaskFromDB(task.ID)
		if err != nil {
			response := models.ErrorResponse{Error: "Задачи не существует"}
			json.NewEncoder(w).Encode(response)
			return
		}

		_, err = ValidateDate(task.Date, TimeFormat)
		if err != nil {
			response := models.ErrorResponse{Error: "Неверный формат даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		err = db.UpdateTaskInDB(task)
		if err != nil {
			response := models.ErrorResponse{Error: fmt.Sprintf("Ошибка обновления задачи: %v", err)}
			json.NewEncoder(w).Encode(response)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
