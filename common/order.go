package common

import (
	"errors"
)

type OrderStatus int

const (
	Cancelled OrderStatus = 0
	Created   OrderStatus = 1
	Confirmed OrderStatus = 2
	Delivered OrderStatus = 3
)

type Order struct {
	Id          string  `json:"id"`
	Status      string  `json:"status,omitempty"`
	Description string  `json:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
}

type OrderRequest struct {
	Status      OrderStatus `json:"status,omitempty"`
	Description string      `json:"description,omitempty"`
	Amount      float64     `json:"amount,omitempty"`
}

// Returns status name
func (status OrderStatus) String() (string, error) {
	statuses := []string{
		"Cancelled",
		"Created",
		"Confirmed",
		"Delivered",
	}
	minStatus := Cancelled
	maxStatus := Delivered

	if status < minStatus || status > maxStatus {
		return "Unknown", errors.New("Invalid order status")
	}
	return statuses[status], nil
}
