package main

import (
	"log"
	"net/http"
	"orders/internal/app"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("GOOSE_DBSTRING")
	if dbURL == "" {
		log.Fatalln("GOOSE_DBSTRING is not found")
	}

	myApp, err := app.NewApp("postgres", dbURL)
	if err != nil {
		log.Fatalln("Can't create db connection:", err)
	}
	defer myApp.DB.Close()

	staticFileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	http.HandleFunc("/", myApp.HomeHandler)
	http.HandleFunc("/orders", myApp.ShowOrdersHandler)
	http.HandleFunc("/orders/{order_uid}", myApp.GetOrderByIdHandler)
	http.HandleFunc("/random/{amount}", myApp.RandomOrdersHandler)

	log.Println("Server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Can't start the server:", err)
	}
}
