package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

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
type OrderRequest common.OrderRequest

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

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
	fmt.Println("Received body: ", request.Body)
	var (
		req OrderRequest
		id  = uuid.New().String()
	)

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return common.ErrorResponse(err), err
	}

	created, err := common.Created.String()
	if err != nil {
		return common.ErrorResponse(err), err
	}

	order := &Order{
		Id:          id,
		Status:      created,
		Description: req.Description,
		Amount:      req.Amount,
	}

	if err := addOrder(order); err != nil {
		return common.ErrorResponse(err), err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func addOrder(order *Order) error {
	var tableName = aws.String(os.Getenv("ORDER_TABLE"))

	item, err := dynamodbattribute.MarshalMap(order)
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
