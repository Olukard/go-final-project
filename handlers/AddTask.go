package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-final-project/db"
	"go-final-project/models"
	"net/http"
	"time"
)

func AddTaskHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var task models.Task
		var buf bytes.Buffer

		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		_, err := buf.ReadFrom(r.Body)
		if err != nil {
			handleError(w, err, "Bad request", 400)
			return
		}

		if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
			handleError(w, err, "Bad request", 400)
			return
		}

		if task.Title == "" {
			err := fmt.Errorf("title cannot be empty")
			handleError(w, err, "Bad request", 400)
			return
		}

		now := time.Now()

		if task.Date == "" {
			task.Date = now.Format(TimeFormat)
		} else {
			dueDate, err := time.Parse(TimeFormat, task.Date)
			if err != nil {
				handleError(w, err, "Bad request", 500)
				return
			}

			if dueDate.Before(now) {
				if task.Repeat != "" {
					nextDueDate, err := NextDate(now, task.Date, task.Repeat)
					if err != nil {
						handleError(w, err, "Bad request", 500)
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
			handleError(w, err, "Internal server error", 400)
			return

		}
		response := models.AddTaskResponse{ID: id}
		json.NewEncoder(w).Encode(response)
	}
}
