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
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
	}

	_, err = ValidateDate(dateStr, TimeFormat)
	if err != nil {
		http.Error(w, "Неверный формат даты", http.StatusBadRequest)
	}

	err = ValdateRepeatRule(repeatStr)
	if err != nil {
		http.Error(w, "Неверный формат правила повторения", http.StatusBadRequest)
	}

	nowDate, _ := time.Parse(TimeFormat, nowStr)

	repeatDate, err := NextDate(nowDate, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, repeatDate)
}
