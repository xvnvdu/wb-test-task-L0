package app

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"orders/cmd/generator"

	k "orders/internal/kafka"
	repo "orders/internal/repository"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
)

type App struct {
	kafkaConsumer *kafka.Reader
	kafkaProducer *kafka.Writer
	repo          *repo.Repository
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		notFound, err := os.ReadFile("web/templates/404.html")
		if err != nil {
			log.Println("Error reading file:", err)
			http.Error(w, "Nothing's here...", http.StatusNotFound)
			return
		}
		if _, err := w.Write([]byte(notFound)); err != nil {
			log.Fatalln("Handler error: HomeHandler:", err)
		}
		return
	}

	html, err := os.ReadFile("web/templates/index.html")
	if err != nil {
		log.Println("Error reading file:", err)
		http.Error(w, "Nothing's here...", http.StatusNotFound)
		return
	}

	if _, err := w.Write([]byte(html)); err != nil {
		log.Fatalln("Handler error: HomeHandler:", err)
	}
}

func (a *App) GetOrderByIdHandler(w http.ResponseWriter, r *http.Request) {
	order_uid := r.PathValue("order_uid")
	ctx := context.Background()

	orderData, err := a.repo.GetOrderById(order_uid, ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Order not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}

	orderJSON, err := json.MarshalIndent(orderData, "", "    ")
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(orderJSON)); err != nil {
		log.Fatalln("Handler error: GetOrderByIdHandler:", err)
	}
}

func (a *App) ShowOrdersHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	ordersList, err := a.repo.GetAllOrders(ctx)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	ordersJSON, err := json.MarshalIndent(ordersList, "", "    ")
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := w.Write([]byte(ordersJSON)); err != nil {
		log.Fatalln("Handler error: ShowOrdersHandler:", err)
	}
}

func (a *App) RandomOrdersHandler(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("amount")
	amount, err := strconv.Atoi(value)

	badRequest := func() {
		result, err := os.ReadFile("web/templates/400.html")
		if err != nil {
			http.Error(w, "Nothing's here...", http.StatusNotFound)
			log.Fatalln("Error reading file:", err)
		}
		if _, err := w.Write([]byte(result)); err != nil {
			log.Fatalln("Handler error: RandomOrdersHandler:", err)
		}
	}

	if err != nil {
		log.Println("Error in internal/app/app.go: amount, err := strconv.Atoi(value):", err)
		badRequest()
	} else {
		if amount <= 0 {
			log.Println("Error creating orders: Value is equal or less than zero")
			badRequest()
			return
		}

		ctx := context.Background()
		orders := generator.MakeRandomOrder(amount)

		orderJSON, err := json.MarshalIndent(orders, "", "    ")
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		err = k.WriteMessage(a.kafkaProducer, ctx, orderJSON)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
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

	repo := &repo.Repository{DB: db}

	k.CreateTopic()
	reader := k.CreateReader()
	writer := k.CreateWriter()

	go k.StartConsuming(reader, repo)

	app := &App{kafkaConsumer: reader, kafkaProducer: writer, repo: repo}
	return app, nil
}

func (a App) Close() {
	err := a.repo.DB.Close()
	if err != nil {
		log.Fatalln("Database connection can't be closed:", err)
	}

	err = a.kafkaConsumer.Close()
	if err != nil {
		log.Fatalln("Kafka stream can't be closed:", err)
	}

	err = a.kafkaProducer.Close()
	if err != nil {
		log.Fatalln("Kafka producer can't be closed:", err)
	}
}
