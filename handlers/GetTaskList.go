package handlers

import (
	"encoding/json"
	"go-final-project/db"
	"net/http"
)

func GetTasksListHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		tasks, err := db.GetListFromDB()
		if err != nil {
			handleError(w, err, "Internal server error", 400)
			return
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"tasks": tasks,
		})
	}
}
