package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"go-final-project/db"
	"go-final-project/handlers"
)

func main() {

	db, err := db.CreateDB()
	if err != nil {
		log.Print(err)
	}

	router := chi.NewRouter()

	fileserver := http.FileServer(http.Dir("./web"))
	router.Handle("/*", fileserver)

	router.Get("/api/nextdate", handlers.NextDateHandler)

	router.Post("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.AddTaskHandler(w, r, db)
	})

	router.Get("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTaskHandler(w, r, db)
	})

	router.Put("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.EditTaskHandler(w, r, db)
	})

	router.Delete("/api/task", func(w http.ResponseWriter, r *http.Request) {
		handlers.DeleteTaskHandler(w, r, db)
	})

	router.Get("/api/tasks", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetTasksListHandler(w, r, db)
	})

	router.Post("/api/task/done", func(w http.ResponseWriter, r *http.Request) {
		handlers.TaskDoneHandler(w, r, db)
	})

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = ":7540"
	}

	log.Printf("Запускаем сервер. Порт%s\n", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		panic(err)
	}

	log.Println("Завершаем работу")
}
