package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-final-project/db"
	"go-final-project/models"
)

func TaskDoneHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		idStr := r.URL.Query().Get("id")

		err := ValidateID(idStr)
		if err != nil {
			response := models.ErrorResponse{Error: "Ошибка получения id задачи"}
			json.NewEncoder(w).Encode(response)
			return
		}

		task, err := db.GetTaskFromDB(idStr)
		if err != nil {
			response := models.ErrorResponse{Error: "Ошибка получения данных"}
			json.NewEncoder(w).Encode(response)
			return
		}

		if task.Repeat == "" {
			db.DeleteTaskFromDB(idStr)
			json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}

		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			response := models.ErrorResponse{Error: "Ошибка установки даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := db.UpdateTaskInDB(task); err != nil {
			response := models.ErrorResponse{Error: "Ошибка обновления задачи"}
			json.NewEncoder(w).Encode(response)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
