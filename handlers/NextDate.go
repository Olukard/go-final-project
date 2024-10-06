package handlers

import (
	"fmt"
	"net/http"
	"time"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {

	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	_, err := ValidateDate(nowStr, TimeFormat)
	if err != nil {
		handleError(w, err, "Internal server error")
		return
	}

	_, err = ValidateDate(dateStr, TimeFormat)
	if err != nil {
		handleError(w, err, "Internal server error")
		return
	}

	err = ValdateRepeatRule(repeatStr)
	if err != nil {
		handleError(w, err, "Internal server error")
		return
	}

	nowDate, _ := time.Parse(TimeFormat, nowStr)

	repeatDate, err := NextDate(nowDate, dateStr, repeatStr)
	if err != nil {
		handleError(w, err, "Internal server error")
		return
	}

	fmt.Fprintln(w, repeatDate)
}
