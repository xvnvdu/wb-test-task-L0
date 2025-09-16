package app

import (
	"database/sql"
	"log"
	"net/http"

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
