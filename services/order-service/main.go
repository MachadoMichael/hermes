package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MachadoMichael/hermes/domain"
	"github.com/MachadoMichael/hermes/infra"
)

func main() {
	config := infra.MSConfig{URL: "amqp://guess:guess@localhost:5672/"}
	client, err := infra.NewMSClient(config)
	if err != nil {
		log.Fatalf("Failed to create messaging system client: %v", err)
	}

	defer client.Close()

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		var order domain.Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

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
		w.Write([]byte("Order created and published successfully"))
	})

	log.Println("Order service started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
