package main

import (
	"log"
	"net/http"
	"orders/internal/app"
	"os"

	"github.com/joho/godotenv"
	"github.com/swaggo/http-swagger"
)

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DB_CONN_STRING")
	if dbURL == "" {
		log.Fatalln("DB_CONN_STRING is not found")
	}

	driver := os.Getenv("DRIVER")
	if driver == "" {
		log.Fatalln("DRIVER is not found")
	}

	myApp, err := app.NewApp(driver, dbURL)
	if err != nil {
		log.Fatalln("Can't create db connection:", err)
	}
	defer myApp.Close()

	// Отдаем статику
	staticFileServer := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticFileServer))

	// Основные эндпоинты
	http.HandleFunc("/", myApp.HomeHandler)
	http.HandleFunc("/orders", myApp.ShowOrdersHandler)
	http.HandleFunc("/orders/{order_uid}", myApp.GetOrderByIdHandler)
	http.HandleFunc("/random/{amount}", myApp.RandomOrdersHandler)

	// Отдаем файл с документацией и рендерим его по эндпоинту /docs
	http.HandleFunc("/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yaml")
	})
	http.Handle("/docs/", httpSwagger.Handler(httpSwagger.URL("/swagger.yaml")))

	log.Println("Server is running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("Can't start the server:", err)
	}
}
