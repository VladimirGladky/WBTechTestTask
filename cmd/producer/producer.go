package main

import (
	kafka "WBTechTestTask/internal/kafka/producer"
	"WBTechTestTask/internal/models"
	"context"

	"github.com/google/uuid"
)

func main() {
	ctx := context.Background()
	p := kafka.NewProducer([]string{"localhost:9092"}, "orders")

	order := models.Order{
		OrderUid: uuid.New().String(),
		Delivery: &models.Delivery{
			Name:    "John Doe",
			Phone:   "1234567890",
			Zip:     "123456",
			City:    "Moscow",
			Address: "123 Main St",
			Region:  "Moscow",
			Email:   "john.doe@example.com",
		},
		Payment: &models.Payment{
			Transaction:  "123",
			RequestId:    "123",
			Currency:     "RUB",
			Provider:     "provider",
			Amount:       123,
			PaymentDt:    123,
			Bank:         "bank",
			DeliveryCost: 123,
			GoodsTotal:   123,
			CustomFee:    123,
		},
		Items: []models.Item{
			{
				ChrtId:      123,
				Price:       123,
				TrackNumber: "123",
				Size:        "123",
				Rid:         "123",
				Name:        "123",
				Sale:        123,
				TotalPrice:  123,
				NmId:        123,
				Brand:       "123",
				Status:      123,
			},
		},
		Locale:            "en",
		TrackNumber:       "123",
		Entry:             "123",
		InternalSignature: "123",
		CustomerId:        "123",
		DeliveryService:   "123",
		ShardKey:          "123",
		SmId:              123,
		DateCreated:       "2025-09-15T04:00:56+03:00",
		OofShard:          "123",
	}
	err := p.SendOrderMessage(ctx, order)
	if err != nil {
		return
	}
}
