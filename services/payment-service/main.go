package main

import (
	"encoding/json"
	"log"
	"time"

	"github.com/MachadoMichael/hermes/domain"
	"github.com/MachadoMichael/hermes/infra"
	"golang.org/x/exp/rand"
)

func main() {
	config := infra.MSConfig{URL: "amqp://guest:guest@rabbitmq:5672/"}
	client, err := infra.NewMSClient(config)
	if err != nil {
		log.Fatalf("Failed to create messaging system client: %v", err)
	}

	defer client.Close()

	msgs, err := client.Consume("orders")
	if err != nil {
		log.Fatalf("Failed to consume from payments queue: %v", err)
	}
	log.Println("Payment service started. Waiting for payments...")

	// Seed random for simulating payment success/failure
	rand.Seed(uint64(time.Now().UnixNano()))

	for msg := range msgs {
		var order domain.Order
		if err := json.Unmarshal(msg.Body, &order); err != nil {
			log.Printf("Failed to unmarshal order: %v", err)
			continue
		}

		log.Printf("Processing payment for order: %v", order.ID)

		paymentStatus := domain.PaymentStatus{
			OrderID:   order.ID,
			PaymentID: generatePaymentID(),
		}

		if rand.Intn(10) < 8 {
			paymentStatus.Status = domain.PAID
		} else {
			paymentStatus.Status = domain.FAILED
			paymentStatus.Reason = "Insufficient funds"
		}

		body, err := json.Marshal(paymentStatus)
		if err != nil {
			log.Printf("Failed to marshal payment status: %v", err)
			continue
		}

		if err := client.Publish("payments", body); err != nil {
			log.Printf("Failed to publish payment status: %v", err)
			continue
		}
		log.Printf("Payment status published for order: %v", order.ID)
	}

}

func generatePaymentID() string {
	return "pay-" + randomString(10)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
