package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

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
func Handler(ctx context.Context, request Request) (Response, error) {
	fmt.Println("Received body: ", request.Body)
	var (
		req OrderRequest
		id  = uuid.New().String()
	)

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(err), err
	}

	created, err := Created.String()
	if err != nil {
		return errorResponse(err), err
	}

	order := &Order{
		Id:          id,
		Status:      created,
		Description: req.Description,
		Amount:      req.Amount,
	}

	if err := addOrder(order); err != nil {
		return errorResponse(err), err
	}

	body, err := json.Marshal(order)
	if err != nil {
		return errorResponse(err), err
	}

	return Response{
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
