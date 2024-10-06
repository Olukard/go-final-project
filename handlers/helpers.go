package handlers

import (
	"encoding/json"
	"fmt"
	"go-final-project/models"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	DayLimit   = 400
	TimeFormat = "20060102"
)

// NextDate вычисляет следующую дату повторения задачи в соответсии с заданным правилом повторения. Таже валидируются форматы дат и самого правила повторения.
func NextDate(now time.Time, date string, repeat string) (string, error) {

	taskDate, err := ValidateDate(date, TimeFormat)
	if err != nil {
		return "", err
	}

	if repeat == "" || date == "" {
		return "", fmt.Errorf("error opening database: %w", err)
	}

	switch string(repeat[0]) {
	case "d":
		if len(repeat) < 3 {
			return "", fmt.Errorf("repeat rule format error, amount of days not specified: %w", err)
		}
		daysToDelay, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return "", fmt.Errorf("repeat rule format error, cannot parse amount of days: %w", err)
		}
		if daysToDelay > DayLimit {
			return "", fmt.Errorf(`maximum day delay exeeded. 
Day delay: %d. Current limit: %d`, daysToDelay, DayLimit)
		}
		nextDate := taskDate.AddDate(0, 0, daysToDelay)

		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, daysToDelay)
		}

		return nextDate.Format(TimeFormat), nil

	case "y":
		if len(repeat) < 3 {
			repeat += " 1"
		}
		years, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return "", fmt.Errorf("repeat rule format error, cannot parse amount of days: %w", err)
		}

		nextDate := taskDate.AddDate(years, 0, 0)

		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(years, 0, 0)
		}

		return nextDate.Format("20060102"), nil

	case "w", "m":
		return "", fmt.Errorf("repeat reule not yet implemented: %w", err)
	default:
		return "", fmt.Errorf("unknown repeat rule: %s", repeat)

	}
}

// ValidateID вилидирует, что поле ID является нумерическим значением
func ValidateID(id string) (err error) {
	_, err = strconv.Atoi(id)
	return err
}

// ValidateDate вилидирует, что дата является нумерическим значением, и сооветсвует заданному формату даты
func ValidateDate(date string, timeFormat string) (resultDate time.Time, err error) {
	resultDate, err = time.Parse(timeFormat, date)
	if err != nil {
		return time.Now(), fmt.Errorf("date format error, returning current date")
	}
	return resultDate, nil
}

// ValidateRepeatRule проверяет, что строка точно соотвествует заданному списку правил повторения задач. Также проверяется формат нумерического занчения в формате правила.
func ValdateRepeatRule(repeat string) (err error) {

	validRepeatRules := map[string]bool{
		"d": true, // Daily
		"w": true, // Weekly
		"m": true, // Monthly
		"y": true, // Yearly
	}

	if repeat == "" {
		return fmt.Errorf("repeat rule cannot be empty")
	}

	if !validRepeatRules[string(repeat[0])] {
		return fmt.Errorf("unknown repeat rule: %s", repeat)
	}

	if repeat == "y" {
		return nil
	}

	if len(repeat) < 3 {
		return fmt.Errorf("repeat rule format error, amount of days not specified")
	}

	if len(repeat) > 1 {
		_, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return fmt.Errorf("repeat rule format error, cannot parse amount of days")
		}
	}

	return nil
}

func handleError(w http.ResponseWriter, err error, msg string) {
	log.Printf("Error: %v", err)
	response := models.ErrorResponse{Error: msg}
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}
