package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-final-project/db"
	"go-final-project/models"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем параметры из запроса
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeatStr := r.URL.Query().Get("repeat")

	// Проверяем, что параметры не пусты
	if nowStr == "" {
		http.Error(w, "Не задана дата", http.StatusBadRequest)
	}
	if dateStr == "" {
		http.Error(w, "Не задана дата", http.StatusBadRequest)
		return
	}
	if repeatStr == "" {
		http.Error(w, "Не задано правило повторения", http.StatusBadRequest)
		return
	}

	// Парсим дату
	nowDate, err := time.Parse(TimeFormat, nowStr)
	if err != nil {
		http.Error(w, "Неверный формат даты1", http.StatusBadRequest)
		return
	}

	repeatDate, err := NextDate(nowDate, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, repeatDate)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {

	//обозначаем структуры
	var task models.Task
	var buf bytes.Buffer

	//устанавливаем заголовок
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//читаем тело запроса
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка чтения запроса"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//декодируем json
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response := models.ErrorResponse{Error: "Ошибка декодирования json"}
		json.NewEncoder(w).Encode(response)
		return
	}

	//проверяем наличие заголовка запроса
	if task.Title == "" {
		response := models.ErrorResponse{Error: "Заголовок задачи не может быть пустым"}
		json.NewEncoder(w).Encode(response)
		return
	}

	now := time.Now()

	// Проверяем дату задачи на пустое значение
	if task.Date == "" {
		task.Date = now.Format(TimeFormat)
	} else {
		// Проверяем парсится ли дата
		dueDate, err := time.Parse(TimeFormat, task.Date)
		if err != nil {
			response := models.ErrorResponse{Error: "Неверный формат даты"}
			json.NewEncoder(w).Encode(response)
			return
		}

		// Проверяем, не раньше ли дата, чем сегодня
		if dueDate.Before(now) {
			// Если дата раньше, вычисляем следующую дату
			if task.Repeat != "" {
				nextDueDate, err := NextDate(now, task.Date, task.Repeat)
				if err != nil {
					response := models.ErrorResponse{Error: "Неверный формат правила повторения"}
					json.NewEncoder(w).Encode(response)
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
		response := models.AddTaskResponse{Error: "Ошибка добавления в базу данных"}
		json.NewEncoder(w).Encode(response)
		return

	}
	response := models.AddTaskResponse{ID: id}
	json.NewEncoder(w).Encode(response)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var task models.Task
	idStr := r.URL.Query().Get("id")

	if idStr == "" {
		response := models.ErrorResponse{Error: "Ошибка получения id задачи"}
		json.NewEncoder(w).Encode(response)
		return
	}

	task, err := db.GetTaskFromDB(idStr)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка получения данных"}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func EditTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	var buf bytes.Buffer

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка чтения запроса"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		response := models.ErrorResponse{Error: "Ошибка декодирования json"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if task.Title == "" {
		response := models.ErrorResponse{Error: "Заголовок задачи не может быть пустым"}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = ValidateDate(task.Date)
	if err != nil {
		response := models.ErrorResponse{Error: "неверное правило повторения"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = ValdateRepeatRule(task.Repeat)
	if err != nil {
		response := models.ErrorResponse{Error: "неверное правило повторения"}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = db.GetTaskFromDB(task.ID)
	if err != nil {
		response := models.ErrorResponse{Error: "Задачи не существует"}
		json.NewEncoder(w).Encode(response)
		return
	}

	_, err = time.Parse(TimeFormat, task.Date)
	if err != nil {
		response := models.ErrorResponse{Error: "Неверный формат даты"}
		json.NewEncoder(w).Encode(response)
		return
	}

	err = db.UpdateTaskInDB(task)
	if err != nil {
		response := models.ErrorResponse{Error: fmt.Sprintf("Ошибка обновления задачи: %v", err)}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}

func GetTasksListHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	tasks, err := db.GetListFromDB()
	if err != nil {
		response := models.ErrorResponse{Error: "Ошибка получения данных"}
		json.NewEncoder(w).Encode(response)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"tasks": tasks,
	})
}
