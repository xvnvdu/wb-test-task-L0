package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"orders/cmd/generator"
	db "orders/internal/database"
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
		log.Println("Error in internal/app/app.go: line 27:", err)
		if _, err := w.Write([]byte("Ooops! Something went wrong!\n" +
			"To generate random order(s), please consider using an INTEGER value.")); err != nil {
			log.Fatalln("Handler error: RandomOrdersHandler:", err)
		}
	} else {
		orders := generator.MakeRandomOrder(amount)
		queries := db.New(a.DB)

		for _, order := range orders {
			err := queries.CreateOrder(r.Context(), db.CreateOrderParams{
				OrderUid:    order.OrderUID,
				TrackNumber: order.TrackNumber,
				Entry:       order.Entry,
				Locale:      order.Locale,
				InternalSignature: sql.NullString{
					String: order.InternalSignature,
					Valid:  order.InternalSignature != "",
				},
				CustomerID:      order.CustomerID,
				DeliveryService: order.DeliveryService,
				Shardkey:        order.Shardkey,
				SmID:            int32(order.SmID),
				DateCreated:     order.DateCreated,
				OofShard:        order.OofShard,
			})
			if err != nil {
				log.Fatalln("Error inserting order:", err)
			}

			err = queries.CreateDelivery(r.Context(), db.CreateDeliveryParams{
				OrderUid: order.OrderUID,
				Name:     order.Delivery.Name,
				Phone:    order.Delivery.Phone,
				Zip:      order.Delivery.Zip,
				City:     order.Delivery.City,
				Address:  order.Delivery.Address,
				Region:   order.Delivery.Region,
				Email:    order.Delivery.Email,
			})
			if err != nil {
				log.Fatalln("Error inserting delivery:", err)
			}

			err = queries.CreatePayment(r.Context(), db.CreatePaymentParams{
				OrderUid:    order.OrderUID,
				Transaction: order.Payment.Transaction,
				RequestID: sql.NullString{
					String: order.Payment.RequestID,
					Valid:  order.Payment.RequestID != "",
				},
				Currency:     order.Payment.Currency,
				Provider:     order.Payment.Provider,
				Amount:       int32(order.Payment.Amount),
				PaymentDt:    int64(order.Payment.PaymentDT),
				Bank:         order.Payment.Bank,
				DeliveryCost: int32(order.Payment.DeliveryCost),
				GoodsTotal:   int32(order.Payment.GoodsTotal),
				CustomFee:    int32(order.Payment.CustomFee),
			})
			if err != nil {
				log.Fatalln("Error inserting payment:", err)
			}

			for _, item := range order.Items {
				err = queries.CreateItem(r.Context(), db.CreateItemParams{
					OrderUid:    order.OrderUID,
					ChrtID:      int32(item.ChrtID),
					TrackNumber: item.TrackNumber,
					Price:       int32(item.Price),
					Rid:         item.Rid,
					Name:        item.Name,
					Sale:        int32(item.Sale),
					Size:        item.Size,
					TotalPrice:  int32(item.TotalPrice),
					NmID:        int32(item.NmID),
					Brand:       item.Brand,
					Status:      int32(item.Status),
				})
				if err != nil {
					log.Fatalln("Error inserting item:", err)
				}
			}
		}

		orderJSON, err := json.MarshalIndent(orders, "", "    ")
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
