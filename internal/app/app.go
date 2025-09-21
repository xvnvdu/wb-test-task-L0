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

func (a *App) ShowOrdersHandler(w http.ResponseWriter, r *http.Request) {
	queries := db.New(a.DB)

	orders, err := queries.GetOrders(r.Context())
	if err != nil {
		log.Fatalln("Error getting orders:", err)
	}

	deliveries, err := queries.GetDelivery(r.Context())
	if err != nil {
		log.Fatalln("Error getting deliveries:", err)
	}

	payments, err := queries.GetPayment(r.Context())
	if err != nil {
		log.Fatalln("Error getting payments:", err)
	}

	items, err := queries.GetItems(r.Context())
	if err != nil {
		log.Fatalln("Error getting items:", err)
	}

	deliveriesMap := make(map[string]generator.Delivery)
	paymentsMap := make(map[string]generator.Payment)
	itemsMap := make(map[string][]generator.Item)

	for _, delivery := range deliveries {
		deliveriesMap[delivery.OrderUid] = generator.Delivery{
			Name:    delivery.Name,
			Phone:   delivery.Phone,
			Zip:     delivery.Zip,
			City:    delivery.City,
			Address: delivery.Address,
			Region:  delivery.Region,
			Email:   delivery.Email,
		}
	}

	for _, payment := range payments {
		paymentsMap[payment.OrderUid] = generator.Payment{
			Transaction:  payment.Transaction,
			RequestID:    payment.RequestID.String,
			Currency:     payment.Currency,
			Provider:     payment.Provider,
			Amount:       int(payment.Amount),
			PaymentDT:    int(payment.PaymentDt),
			Bank:         payment.Bank,
			DeliveryCost: int(payment.DeliveryCost),
			GoodsTotal:   int(payment.GoodsTotal),
			CustomFee:    int(payment.CustomFee),
		}
	}

	for _, item := range items {
		if itemsMap[item.OrderUid] == nil {
			itemsMap[item.OrderUid] = []generator.Item{}
		}
		itemsMap[item.OrderUid] = append(itemsMap[item.OrderUid], generator.Item{
			ChrtID:      int(item.ChrtID),
			TrackNumber: item.TrackNumber,
			Price:       int(item.Price),
			Rid:         item.Rid,
			Name:        item.Name,
			Sale:        int(item.Sale),
			Size:        item.Size,
			TotalPrice:  int(item.TotalPrice),
			NmID:        int(item.NmID),
			Brand:       item.Brand,
			Status:      int(item.Status),
		})
	}

	var ordersList []*generator.Order
	for _, order := range orders {
		ordersList = append(ordersList, &generator.Order{
			OrderUID:          order.OrderUid,
			TrackNumber:       order.TrackNumber,
			Entry:             order.Entry,
			Delivery:          deliveriesMap[order.OrderUid],
			Payment:           paymentsMap[order.OrderUid],
			Items:             itemsMap[order.OrderUid],
			Locale:            order.Locale,
			InternalSignature: order.InternalSignature.String,
			CustomerID:        order.CustomerID,
			DeliveryService:   order.DeliveryService,
			Shardkey:          order.Shardkey,
			SmID:              int(order.SmID),
			DateCreated:       order.DateCreated,
			OofShard:          order.OofShard,
		})
	}

	orderJSON, err := json.MarshalIndent(ordersList, "", "    ")
	if err != nil {
		log.Println("Error marshalling JSON:", err)
	}

	if _, err := w.Write([]byte(orderJSON)); err != nil {
		log.Fatalln("Handler error: ShowOrdersHandler:", err)
	}
}

func (a *App) RandomOrdersHandler(w http.ResponseWriter, r *http.Request) {
	value := r.PathValue("amount")
	amount, err := strconv.Atoi(value)

	if err != nil {
		log.Println("Error in internal/app/app.go: line 130:", err)
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
