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
			handleError(w, err, "Bad request", 500)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			handleError(w, err, "Bad request", 500)
			return
		}

		if task.Title == "" {
			err := fmt.Errorf("title cannot be empty")
			handleError(w, err, "Bad request", 500)
			return
		}

		_, err = ValidateDate(task.Date, TimeFormat)
		if err != nil {
			handleError(w, err, "Bad request", 500)
			return
		}

		err = ValdateRepeatRule(task.Repeat)
		if err != nil {
			handleError(w, err, "Bad request", 500)
			return
		}

		_, err = db.GetTaskFromDB(task.ID)
		if err != nil {
			handleError(w, err, "Bad request", 500)
			return
		}

		_, err = ValidateDate(task.Date, TimeFormat)
		if err != nil {
			handleError(w, err, "Bad request", 500)
			return
		}

		err = db.UpdateTaskInDB(task)
		if err != nil {
			handleError(w, err, "Internal server error", 400)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})
	}
}
