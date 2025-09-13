package main

import (
	"log"
	"net/http"

	"orders/internal/app"
)

func main() {
	myApp := &app.App{}
	http.HandleFunc("/", myApp.HomeHandler)

	log.Println("Сервер запущен на http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Ошибка запуска сервера:", err)
	}
}
