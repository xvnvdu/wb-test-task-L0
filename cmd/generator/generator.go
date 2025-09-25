package generator

import (
	"time"

	gf "github.com/brianvoe/gofakeit/v7"
)

func MakeRandomOrder(amount int) []*Order {
	var orders []*Order

	for range amount {
		orderUID := gf.UUID()
		trackNumber := gf.Regex("[A-Z0-9]{14}")

		itemsCount := gf.Number(1, 5)
		items := make([]Item, itemsCount)

		var totalGoodsPrice int
		for i := 0; i < itemsCount; i++ {
			price := gf.Number(100, 10000)
			sale := gf.Number(0, 90)
			totalPrice := price * (100 - sale) / 100

			items[i] = Item{
				OrderUID:    orderUID,
				ChrtID:      gf.Number(1000000, 9999999),
				TrackNumber: trackNumber,
				Price:       price,
				Rid:         gf.UUID(),
				Name:        gf.ProductName(),
				Sale:        sale,
				Size:        gf.RandomString([]string{"S", "M", "L", "XL"}),
				TotalPrice:  totalPrice,
				NmID:        gf.Number(1000000, 9999999),
				Brand:       gf.Company(),
				Status:      gf.Number(100, 500),
			}

			totalGoodsPrice += totalPrice
		}

		deliveryCost := gf.Number(0, 1500)
		customFee := gf.Number(0, 100)

		orders = append(orders, &Order{
			OrderUID:    orderUID,
			TrackNumber: trackNumber,
			Entry:       gf.RandomString([]string{"WBIL", "WBIL_1", "WBIL_2", "WBIL_3"}),
			Delivery: Delivery{
				Name:    gf.Name(),
				Phone:   gf.Phone(),
				Zip:     gf.Zip(),
				City:    gf.City(),
				Address: gf.Street(),
				Region:  gf.State(),
				Email:   gf.Email(),
			},
			Payment: Payment{
				// Transaction:  gf.UUID(),
				Transaction:  gf.Regex("[A-Z0-9]{14}"),
				RequestID:    gf.RandomString([]string{"", "1", "2", "3", "4", "5", "6", "7"}),
				Currency:     gf.CurrencyShort(),
				Provider:     gf.RandomString([]string{"wbpay", "tpay", "sberpay", "applepay"}),
				Amount:       deliveryCost + totalGoodsPrice + customFee,
				PaymentDT:    gf.Number(1000000000, 9999999999),
				Bank:         gf.RandomString([]string{"TBank", "Alpha", "Sber", "Ozon"}),
				DeliveryCost: deliveryCost,
				GoodsTotal:   totalGoodsPrice,
				CustomFee:    customFee,
			},
			Items:             items,
			Locale:            gf.RandomString([]string{"ru", "en", "cz", "es", "uk", "dk"}),
			InternalSignature: gf.UUID(),
			CustomerID:        gf.UUID(),
			DeliveryService:   gf.RandomString([]string{"SDEK", "Pochta Rossii"}),
			Shardkey:          gf.RandomString([]string{"4", "5", "6", "7", "8", "9"}),
			SmID:              gf.Number(10, 199),
			DateCreated:       time.Now(),
			OofShard:          gf.RandomString([]string{"1", "2", "3", "4", "5", "6", "7"}),
		})
	}

	return orders
}
