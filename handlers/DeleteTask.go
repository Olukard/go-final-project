package handlers

import (
	"encoding/json"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	idStr := r.URL.Query().Get("id")

	err := ValidateID(idStr)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка получения id задачи"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = db.DeleteTaskFromDB(idStr)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка удаления задачи"}
		json.NewEncoder(w).Encode(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})

}
