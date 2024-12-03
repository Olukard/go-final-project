package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"go-final-project/db"
)

func TaskDoneHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		idStr := r.URL.Query().Get("id")

		err := ValidateID(idStr)
		if err != nil {
			handleError(w, err, "Bad request", 500)
			return
		}

		task, err := db.GetTaskFromDB(idStr)
		if err != nil {
			handleError(w, err, "Internal server error", 400)
			return
		}

		if task.Repeat == "" {
			db.DeleteTaskFromDB(idStr)
			json.NewEncoder(w).Encode(map[string]interface{}{})
			return
		}

		task.Date, err = NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			handleError(w, err, "Internal server error", 400)
			return
		}

		if err := db.UpdateTaskInDB(task); err != nil {
			handleError(w, err, "Internal server error", 400)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
