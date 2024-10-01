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

	router.Post("/api/task", handlers.AddTaskHandler(db))
	router.Get("/api/task", handlers.GetTaskHandler(db))
	router.Put("/api/task", handlers.EditTaskHandler(db))
	router.Delete("/api/task", handlers.DeleteTaskHandler(db))
	router.Get("/api/tasks", handlers.GetTasksListHandler(db))
	router.Post("/api/task/done", handlers.TaskDoneHandler(db))

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
