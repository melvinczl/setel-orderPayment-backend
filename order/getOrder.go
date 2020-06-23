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
	"github.com/melvinczl/setel-orderPayment-backend/common"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
type Order common.Order

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
	fmt.Println("GetOrder")

	var orderId = request.PathParameters["id"]

	orders, err := fetchOrder(orderId)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	body, err := json.Marshal(&orders)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func fetchOrder(orderId string) ([]Order, error) {
	var tableName = aws.String(os.Getenv("ORDER_TABLE"))
	var orders []Order

	if orderId == "" {
		input := &dynamodb.ScanInput{
			TableName: tableName,
		}

		result, err := ddb.Scan(input)
		if err != nil {
			return nil, err
		}

		for _, i := range result.Items {
			order := Order{}

			if err := dynamodbattribute.UnmarshalMap(i, &order); err != nil {
				return nil, err
			}
			orders = append(orders, order)
		}
	} else { // Get specific order by ID
		input := &dynamodb.GetItemInput{
			TableName: tableName,
			Key: map[string]*dynamodb.AttributeValue{
				"id": {
					S: aws.String(orderId),
				},
			},
		}

		result, err := ddb.GetItem(input)
		if err != nil {
			return nil, err
		}

		order := Order{}
		if err := dynamodbattribute.UnmarshalMap(result.Item, &order); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func main() {
	lambda.Start(Handler)
}
