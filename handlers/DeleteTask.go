package handlers

import (
	"encoding/json"
	"go-final-project/db"
	"net/http"
)

func DeleteTaskHandler(db *db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")

		idStr := r.URL.Query().Get("id")

		err := ValidateID(idStr)
		if err != nil {
			handleError(w, err, "Internal server error")
			return
		}

		err = db.DeleteTaskFromDB(idStr)
		if err != nil {
			handleError(w, err, "Internal server error")
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{})

	}
}
