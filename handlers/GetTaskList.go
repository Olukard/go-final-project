package handlers

import (
	"encoding/json"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
)

func GetTasksListHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		tasks, err := db.GetListFromDB()
		if err != nil {
			response := models.ErrorResponse{Error: "Ошибка получения данных"}
			json.NewEncoder(w).Encode(response)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"tasks": tasks,
		})
	}
}
