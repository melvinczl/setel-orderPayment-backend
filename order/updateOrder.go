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

func Handler(ctx context.Context, request Request) (Response, error) {
	fmt.Println("UpdateOrder")

	var (
		id  = request.PathParameters["id"]
		req OrderRequest
	)

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return errorResponse(err), err
	}

	if err := updateOrder(id, req); err != nil {
		return errorResponse(err), err
	}

	return Response{
		Body:       request.Body,
		StatusCode: 200,
	}, nil
}

func updateOrder(orderId string, req OrderRequest) error {
	var (
		tableName = aws.String(os.Getenv("ORDER_TABLE"))
		status    = "status"
	)

	orderStatus, err := req.Status.String()
	if err != nil {
		return err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(orderId),
			},
		},
		ExpressionAttributeNames: map[string]*string{
			"#s": &status,
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":s": {
				S: aws.String(orderStatus),
			},
		},
		UpdateExpression: aws.String("set #s = :s"),
		ReturnValues:     aws.String("UPDATED_NEW"),
	}

	_, err = ddb.UpdateItem(input)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
