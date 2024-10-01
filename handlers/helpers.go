package handlers

import (
	"fmt"
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
		return "", fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	switch string(repeat[0]) {
	case "d":
		if len(repeat) < 3 {
			return "", fmt.Errorf("ошибка формата правила повторения: %w", err)
		}
		daysToDelay, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return "", fmt.Errorf("ошибка формата правила повторения: %w", err)
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
			return "", fmt.Errorf("ошибка формата правила повторения: %w", err)
		}

		nextDate := taskDate.AddDate(years, 0, 0)

		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(years, 0, 0)
		}

		return nextDate.Format("20060102"), nil

	case "w", "m":
		return "", fmt.Errorf("правило повторения не поддерживается: %w", err)
	default:
		return "", fmt.Errorf("неизвестное правило повторения: %s", repeat)

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
		return time.Now(), fmt.Errorf("ошибка формата даты, возвращаем сегодняшнее число")
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
		return fmt.Errorf("правило не может быть пустым")
	}

	if !validRepeatRules[string(repeat[0])] {
		return fmt.Errorf("правило не существует")
	}

	if repeat == "y" {
		return nil
	}

	if len(repeat) < 3 {
		return fmt.Errorf("ошибка формата правила повторения")
	}

	if len(repeat) > 1 {
		_, err := strconv.Atoi(repeat[2:])
		if err != nil {
			return fmt.Errorf("ошибка формата правила повторения")
		}
	}

	return nil
}
