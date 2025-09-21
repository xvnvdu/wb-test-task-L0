package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"orders/cmd/generator"
	"strconv"

	_ "github.com/lib/pq"
)

type App struct {
	DB *sql.DB
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello there !")); err != nil {
		log.Fatalln("Handler error: HomeHandler:", err)
	}
}

func (a *App) RandomOrdersHandler(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("amount")
	amount, err := strconv.Atoi(value)

	if err != nil {
		log.Println("Error in internal/app/app.go: line 32:", err)
		if _, err := w.Write([]byte("Ooops! Something went wrong!\n" +
			"To generate random order(s), please consider using an INTEGER value.")); err != nil {
			log.Fatalln("Handler error: RandomOrdersHandler:", err)
		}
	} else {
		order := generator.MakeRandomOrder(amount)

		orderJSON, err := json.MarshalIndent(order, "", "    ")
		if err != nil {
			log.Println("Error marshalling JSON:", err)
		}

		if _, err := w.Write([]byte(orderJSON)); err != nil {
			log.Fatalln("Handler error: RandomOrdersHandler:", err)
		}
	}
}

func NewApp(driverName, dataSourceName string) (*App, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		mainErr := err

		if closeErr := db.Close(); closeErr != nil {
			log.Println("Database connection can't be closed:", closeErr)
		}
		return nil, mainErr
	}

	app := &App{DB: db}
	return app, nil
}
