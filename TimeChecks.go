package main

import (
	"fmt"
	"strconv"
	"time"
)

const DayLimit = 400

func NextDate(now time.Time, date string, repeat string) (string, error) {

	taskDate, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("ошибка формата даты")
	}

	if repeat == "" || date == "" {
		return "", fmt.Errorf("ошибка ввода")
	}

	switch string(repeat[0]) {
	case "d":
		days, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return "", fmt.Errorf("ошибка формата правила повторения")
		}
		if days > DayLimit {
			return "", fmt.Errorf(`превышен максимальный лимит переноса дней. 
Запрошенный перенос: %d. Допустимый лимит: %d`, days, DayLimit)
		}
		nextDate := taskDate.AddDate(0, 0, days)

		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, days)
		}

		return nextDate.Format("20060102"), nil

	case "y":
		if len(repeat) < 3 {
			repeat += " 1"
		}
		years, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return "", fmt.Errorf("ошибка формата правила повторения")
		}

		nextDate := taskDate.AddDate(years, 0, 0)

		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(years, 0, 0)
		}

		return nextDate.Format("20060102"), nil

	case "w", "m":
		return "", fmt.Errorf("правило повторения не поддерживается")
	default:
		return "", fmt.Errorf("неизвестное правило повторения: %s", repeat)

	}
}
