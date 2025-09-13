package app

import (
	"log"
	"net/http"
)

type App struct{}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Привет, мир !")); err != nil {
		log.Fatalln("Ошибка хэндлера: HomeHandler:", err)
	}
}
