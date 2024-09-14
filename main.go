package main

import (
	"fmt"
	"net/http"
)

func main() {

	//добавить выбор порта
	port := ":7540"

	fileServer := http.FileServer(http.Dir("./web"))

	http.Handle("/", fileServer)

	fmt.Println("Проверяем наличие базы данных...")
	if !checkDBexists() {
		fmt.Println("База данных не найдена, создаем...")
		CreateDB()
	} else {
		fmt.Println("База данных найдена.")
	}

	fmt.Printf("Запускаем сервер. Порт%s\n", port)
	err := http.ListenAndServe(port, fileServer)
	if err != nil {
		panic(err)
	}

	fmt.Println("Завершаем работу")

}
