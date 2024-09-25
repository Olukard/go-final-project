package handlers

import (
	"encoding/json"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	idStr := r.URL.Query().Get("id")

	err := ValidateID(idStr)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка получения id задачи"}
		json.NewEncoder(w).Encode(response)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	task, err = db.GetTaskFromDB(idStr)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка получения данных"}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(task)
}
