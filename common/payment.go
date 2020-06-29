package common

import (
// "errors"
)

type Payment struct {
	Id          string  `json:"id"`
	Amount      float64 `json:"amount,omitempty"`
	Description string  `json:"description,omitempty"`
	OrderId     string  `json:"orderId,omitempty"`
	Status      string  `json:"status,omitempty"`
	BillingInfo string  `json:"billingInfo,omitempty"` //mock billing info object...
	CreatedAt   string  `json:"createdAt,omitempty"`
}

type PaymentRequest struct {
	AuthDetails  string `json:"authDetails,omitempty"` //mock auth info...
	OrderDetails Order  `json:"orderDetails,omitempty"`
}

type PaymentResponse struct {
	RefId  string  `json:"refId,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	Status string  `json:"status,omitempty"`
}

type PaymentStatus int

const (
	PaymentDeclined  PaymentStatus = 0
	PaymentConfirmed PaymentStatus = 1
)

var (
	orderStatuses = []string{
		"Declined",
		"Confirmed",
	}
)

// Returns payment status name
func (status PaymentStatus) String() string {
	minStatus := PaymentDeclined
	maxStatus := PaymentConfirmed

	if status < minStatus || status > maxStatus {
		return "Unknown"
	}
	return orderStatuses[status]
}
