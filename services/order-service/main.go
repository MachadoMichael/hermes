package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MachadoMichael/hermes/domain"
	"github.com/MachadoMichael/hermes/infra"
)

func main() {
	config := infra.MSConfig{URL: "amqp://guest:guest@rabbitmq:5672/"}
	client, err := infra.NewMSClient(config)
	if err != nil {
		log.Fatalf("Failed to connect to messaging system: %v", err)
	}
	defer client.Close()

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		var order domain.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("Received order: %s", order.ID)

		body, err := json.Marshal(order)
		if err != nil {
			http.Error(w, "Failed to serialize order", http.StatusInternalServerError)
			return
		}

		if err := client.Publish("orders", body); err != nil {
			http.Error(w, "Failed to publish order", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Fatal(w.Write([]byte("Order created and published")))
	})

	log.Println("Order service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
