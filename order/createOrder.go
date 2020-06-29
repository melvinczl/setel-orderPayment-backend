package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	sdklambda "github.com/aws/aws-sdk-go/service/lambda"
	"github.com/google/uuid"
	"github.com/melvinczl/setel-orderPayment-backend/common"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest
type OrderRequest common.OrderRequest
type OrderResponse common.OrderResponse
type PaymentRequest common.PaymentRequest
type PaymentResponse common.PaymentResponse
type Payload common.Payload

var ddb *dynamodb.DynamoDB
var lambdaClient *sdklambda.Lambda

func init() {
	region := os.Getenv("AWS_REGION")
	awsConfig := &aws.Config{
		Region: &region,
	}

	if sess, err := session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	}); err != nil {
		msg := fmt.Sprintf("Failed to connect to AWS: %s", err.Error())
		fmt.Println(msg)
	} else {
		ddb = dynamodb.New(sess)
		lambdaClient = common.GetLambdaClient(sess, awsConfig)
	}
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("Received body: %s\n", request.Body)
	var (
		req OrderRequest
		id  = uuid.New().String()
	)

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		return common.ErrorResponse(err), err
	}

	status := common.Created.String()
	if status == "Unknown" {
		err := errors.New("Invalid order status")
		return common.ErrorResponse(err), err
	}

	order := &common.Order{
		Id:          id,
		Status:      status,
		Description: req.Description,
		Amount:      req.Amount,
	}

	if err := addOrder(order); err != nil {
		return common.ErrorResponse(err), err
	}

	paymentStatus, err := makePayment(order)
	if err != nil {
		return common.ErrorResponse(err), err
	}

	if err := common.UpdateOrderStatus(lambdaClient, paymentStatus, order); err != nil {
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

func addOrder(order *common.Order) error {
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

	fmt.Println("Order created: " + order.Id)
	return nil
}

func makePayment(order *common.Order) (string, error) {
	var (
		paymentFuncName = os.Getenv("PROC_PAYMENT_FUNCTION")
		resp            Response
		paymentResp     PaymentResponse
	)

	req := &PaymentRequest{
		AuthDetails: "some auth data...",
		OrderDetails: common.Order{
			Id:          order.Id,
			Status:      order.Status,
			Description: order.Description,
			Amount:      order.Amount,
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	p := &Payload{
		Body: string(body),
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return "", err
	}

	fmt.Println("Invoking lambda: " + paymentFuncName)
	fmt.Printf("payload: %s\n", string(payload))
	result, err := lambdaClient.Invoke(&sdklambda.InvokeInput{
		FunctionName: aws.String(paymentFuncName),
		Payload:      payload,
	})
	if err != nil {
		return "", err
	}
	fmt.Printf("Result: %s\n", string(result.Payload))

	statusCode := int(*result.StatusCode)
	if statusCode != 200 {
		fmt.Println("Error in processPayment, StatusCode: " + strconv.Itoa(statusCode))
	}

	if err = json.Unmarshal(result.Payload, &resp); err != nil {
		return "", err
	}

	if err = json.Unmarshal([]byte(resp.Body), &paymentResp); err != nil {
		return "", err
	}

	fmt.Println("Payment status: " + paymentResp.Status)
	return paymentResp.Status, nil
}

func main() {
	lambda.Start(Handler)
}
