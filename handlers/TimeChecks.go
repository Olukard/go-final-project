package handlers

import (
	"fmt"
	"strconv"
	"time"
)

const DayLimit = 400
const TimeFormat = "20060102"

func NextDate(now time.Time, date string, repeat string) (string, error) {

	taskDate, err := time.Parse(TimeFormat, date)
	if err != nil {
		return "", fmt.Errorf("ошибка формата даты")
	}

	if repeat == "" || date == "" {
		return "", fmt.Errorf("ошибка ввода")
	}

	switch string(repeat[0]) {
	case "d":
		if len(repeat) < 3 {
			return "", fmt.Errorf("ошибка формата правила повторения")
		}
		if repeat[2:] == "1" {
			return date, nil
		}
		daysToDelay, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return "", fmt.Errorf("ошибка формата правила повторения")
		}
		if daysToDelay > DayLimit {
			return "", fmt.Errorf(`превышен максимальный лимит переноса дней. 
Запрошенный перенос: %d. Допустимый лимит: %d`, daysToDelay, DayLimit)
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
