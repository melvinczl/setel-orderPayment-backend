package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"github.com/melvinczl/setel-orderPayment-backend/common"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
type Order common.Order
type Payment common.Payment
type PaymentRequest common.PaymentRequest
type PaymentResponse common.PaymentResponse

var ddb *dynamodb.DynamoDB

func init() {
	region := os.Getenv("AWS_REGION")

	if session, err := session.NewSession(&aws.Config{
		Region: &region,
	}); err != nil {
		msg := fmt.Sprintf("Failed to connect to AWS: %s", err.Error())
		fmt.Println(msg)
	} else {
		ddb = dynamodb.New(session)
	}
}

func Handler(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Received body: %s\n", request.Body)
	var (
		req PaymentRequest
		id  = uuid.New().String()
	)

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return common.ErrorResponse(err), err
	}

	status := common.PaymentConfirmed.String()
	if processPayment(req) == false {
		status = common.PaymentDeclined.String()
	}
	if status == "Unknown" {
		err := errors.New("Invalid payment status")
		return common.ErrorResponse(err), err
	}

	createdAt := time.Now().Format(common.TimeLayout)
	_, err := time.Parse(common.TimeLayout, createdAt)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	payment := &Payment{
		Id:          id,
		Amount:      req.OrderDetails.Amount,
		Description: req.OrderDetails.Description,
		OrderId:     req.OrderDetails.Id,
		Status:      status,
		BillingInfo: "",
		CreatedAt:   createdAt,
	}

	if err := addPayment(payment); err != nil {
		return common.ErrorResponse(err), err
	}

	paymentResp := &PaymentResponse{
		RefId:  payment.Id,
		Amount: payment.Amount,
		Status: payment.Status,
	}

	body, err := json.Marshal(paymentResp)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func processPayment(req PaymentRequest) bool {
	//some mock logic...
	rSrc := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(rSrc)
	result := rnd.Intn(4)

	if result > 1 {
		return true
	}
	return false
}

func addPayment(payment *Payment) error {
	var tableName = aws.String(os.Getenv("PAYMENT_TABLE"))

	item, err := dynamodbattribute.MarshalMap(payment)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: tableName,
	}

	if _, err := ddb.PutItem(input); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
