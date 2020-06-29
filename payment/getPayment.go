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
type Payment common.Payment

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
	fmt.Println("GetPayment")

	var orderId = request.QueryStringParameters["orderId"]
	fmt.Printf("orderId: %v\n", orderId)

	payments, err := fetchPayment(orderId)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	body, err := json.Marshal(&payments)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(body),
		StatusCode: 200,
	}, nil
}

func fetchPayment(orderId string) ([]Payment, error) {
	var tableName = aws.String(os.Getenv("PAYMENT_TABLE"))
	var payments []Payment

	input := &dynamodb.ScanInput{
		TableName: tableName,
	}

	result, err := ddb.Scan(input)
	if err != nil {
		return nil, err
	}

	for _, i := range result.Items {
		payment := Payment{}

		if err := dynamodbattribute.UnmarshalMap(i, &payment); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func main() {
	lambda.Start(Handler)
}
