package common

import (
	"errors"
)

type Order struct {
	Id          string  `json:"id"`
	Status      string  `json:"status,omitempty"`
	Description string  `json:"description,omitempty"`
	Amount      float64 `json:"amount,omitempty"`
}

type OrderRequest struct {
	Status      OrderStatus `json:"status"`
	Description string      `json:"description,omitempty"`
	Amount      float64     `json:"amount,omitempty"`
}

type OrderResponse struct {
	Id          string      `json:"id,omitempty"`
	Status      OrderStatus `json:"status"`
	Description string      `json:"description,omitempty"`
	Amount      float64     `json:"amount,omitempty"`
}

type OrderStatus int

const (
	Cancelled OrderStatus = 0
	Created   OrderStatus = 1
	Confirmed OrderStatus = 2
	Delivered OrderStatus = 3
)

var (
	statuses = []string{
		"Cancelled",
		"Created",
		"Confirmed",
		"Delivered",
	}
)

// Returns status name
func (status OrderStatus) String() string {
	minStatus := Cancelled
	maxStatus := Delivered

	if status < minStatus || status > maxStatus {
		return "Unknown"
	}
	return statuses[status]
}

// Returns order status value
func GetOrderStatus(orderStatus string) (int, error) {
	for status, v := range statuses {
		if orderStatus == v {
			return status, nil
		}
	}
	return -1, errors.New("Invalid order status")
}
