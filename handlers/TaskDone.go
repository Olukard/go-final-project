package handlers

import (
	"go-final-project/models"
	"net/http"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	idStr := r.URL.Query().Get("id")

}
