package domain

import (
	"encoding/json"
	"fmt"
)

type Order struct {
	ID    string `json:"id"`
	Items []Item `json:"items"`
}

type Item struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type StatusEnum int

const (
	PENDING StatusEnum = iota
	PAID
	FAILED
)

type PaymentStatus struct {
	OrderID   string     `json:"order_id"`
	Status    StatusEnum `json:"status"`
	Reason    string     `json:"reason"`
	PaymentID string     `json:"payment_id"`
}

// String representation of StatusEnum
func (s StatusEnum) String() string {
	return [...]string{"PENDING", "PAID", "FAILED"}[s]
}

// MarshalJSON converts StatusEnum to a JSON string
func (s StatusEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON converts a JSON string to StatusEnum
func (s *StatusEnum) UnmarshalJSON(data []byte) error {
	var statusStr string
	if err := json.Unmarshal(data, &statusStr); err != nil {
		return err
	}

	switch statusStr {
	case "PENDING":
		*s = PENDING
	case "PAID":
		*s = PAID
	case "FAILED":
		*s = FAILED
	default:
		return fmt.Errorf("invalid status: %s", statusStr)
	}

	return nil
}
