package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/ds124wfegd/tech_wildberries_Go/internal/entity"
	"github.com/segmentio/kafka-go"
)

func main() {
	// Connecting to Kafka
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9094", "orders", 0)
	if err != nil {
		log.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer conn.Close()

	log.Println("Starting Kafka producer...")

	// generate and send test orders
	for i := 1; i <= 10; i++ {
		order := generateTestOrder(i)

		// make JSON
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order: %v", err)
			continue
		}

		// Write message in kafka
		_, err = conn.WriteMessages(kafka.Message{
			Topic: "orders",
			Value: orderJSON,
		})
		if err != nil {
			log.Printf("Error sending message: %v", err)
			continue
		}

		log.Printf("Sent order %d: %s", i, order.OrderUID)
		time.Sleep(2 * time.Second)
	}

	log.Println("Finished sending test orders")
}

// generate test order
func generateTestOrder(id int) *entity.Order {
	orderUID := fmt.Sprintf("test-order-%d-%d", id, time.Now().Unix())

	return &entity.Order{
		OrderUID:          orderUID,
		TrackNumber:       fmt.Sprintf("TRACK%d", id),
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        fmt.Sprintf("customer-%d", id),
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: entity.Delivery{
			Name:    fmt.Sprintf("Test User %d", id),
			Phone:   fmt.Sprintf("+972000000%d", id),
			Zip:     "2639809",
			City:    "Test City",
			Address: fmt.Sprintf("Test Address %d", id),
			Region:  "Test Region",
			Email:   fmt.Sprintf("test%d@gmail.com", id),
		},
		Payment: entity.Payment{
			Transaction:  orderUID,
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1000 + rand.Intn(1000),
			PaymentDt:    time.Now().Unix(),
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []entity.Item{
			{
				ChrtID:      9934930 + id,
				TrackNumber: fmt.Sprintf("TRACK%d", id),
				Price:       453,
				Rid:         fmt.Sprintf("rid-%d", id),
				Name:        fmt.Sprintf("Test Product %d", id),
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212 + id,
				Brand:       "Test Brand",
				Status:      202,
			},
		},
	}
}
