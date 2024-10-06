package handlers

import (
	"encoding/json"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
)

func GetTaskHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		var task models.Task
		idStr := r.URL.Query().Get("id")

		err := ValidateID(idStr)
		if err != nil {
			handleError(w, err, "Internal server error")
			return
		}

		task, err = db.GetTaskFromDB(idStr)
		if err != nil {
			handleError(w, err, "Internal server error")
			return
		}

		json.NewEncoder(w).Encode(task)
	}
}
