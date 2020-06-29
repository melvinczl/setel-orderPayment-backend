package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
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

func UpdateOrderStatus(lambdaClient *lambda.Lambda, paymentStatus string, order *Order) error {
	var (
		updateOrderFunc = os.Getenv("UPDATE_ORDER_FUNCTION")
		resp            Response
		orderResp       OrderResponse
		orderStatus     OrderStatus
	)

	switch paymentStatus {
	case PaymentConfirmed.String():
		orderStatus = Confirmed
	case PaymentDeclined.String():
		orderStatus = Cancelled
	default:
		orderStatus = Created
	}

	req := &OrderRequest{
		Status:      orderStatus,
		Description: order.Description,
		Amount:      order.Amount,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	pathparams := make(map[string]string)
	pathparams["id"] = order.Id
	payload, err := APIRequstPayload(body, pathparams)
	if err != nil {
		return err
	}

	fmt.Println("Invoking lambda: " + updateOrderFunc)
	fmt.Printf("payload: %v\n", string(payload))
	result, err := lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String(updateOrderFunc),
		Payload:      payload,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Result: %v\n", string(result.Payload))

	statusCode := int(*result.StatusCode)
	if statusCode != 200 {
		fmt.Println("Error in processPayment, StatusCode: " + strconv.Itoa(statusCode))
	}

	if err = json.Unmarshal(result.Payload, &resp); err != nil {
		return err
	}

	if err = json.Unmarshal([]byte(resp.Body), &orderResp); err != nil {
		return err
	}

	fmt.Printf("Order status: %s\n", orderResp.Status.String())
	return nil
}
