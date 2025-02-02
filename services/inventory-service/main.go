package main

import (
	"encoding/json"
	"log"

	"github.com/MachadoMichael/hermes/domain"
	"github.com/MachadoMichael/hermes/infra"
)

type InventoryUpdate struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func main() {
	config := infra.MSConfig{URL: "amqp://guest:guest@rabbitmq:5672/"}
	client, err := infra.NewMSClient(config)
	if err != nil {
		log.Fatalf("Failed to create messaging system client: %v", err)
	}

	defer client.Close()

	msgs, err := client.Consume("payments")
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}
	log.Println("Invertory service started. Waiting for payments...")

	for msg := range msgs {
		var paymentStatus domain.PaymentStatus
		if err := json.Unmarshal(msg.Body, &paymentStatus); err != nil {
			log.Printf("Failed to unmarshal payment status: %v", err)
			continue
		}

		if paymentStatus.Status != domain.PAID {
			log.Printf("Skipping inventory update for order: %s (Payment status: %v)", paymentStatus.OrderID, paymentStatus.Status)
			continue
		}

		log.Printf("Processing inventory update for order: %s", paymentStatus.OrderID)

		inventoryUpdate := InventoryUpdate{
			OrderID: paymentStatus.OrderID,
			Status:  "Updated",
		}

		body, err := json.Marshal(inventoryUpdate)
		if err != nil {
			log.Printf("Failed to marshal inventory update: %v", err)
			continue
		}

		if err := client.Publish("inventory", body); err != nil {
			log.Printf("Failed to publish inventory update: %v", err)
			continue
		}

		log.Printf("Inventory updated for order: %v", paymentStatus.OrderID)
	}

}
