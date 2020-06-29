package common

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type Response events.APIGatewayProxyResponse
type Request events.APIGatewayProxyRequest

type Payload struct {
	Body           string            `json:"body"`
	PathParameters map[string]string `json:"pathParameters,omitempty"`
}

const (
	TimeLayout = "2006-01-02T15:04:05-0700"
)

// Generic Http error response
func ErrorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: 500,
	}
}

func GetLambdaClient(sess *session.Session, awsConfig *aws.Config) *lambda.Lambda {
	return lambda.New(sess, awsConfig)
}

func APIRequstPayload(payload []byte, pathParameters map[string]string) ([]byte, error) {
	p := &Payload{
		Body:           string(payload),
		PathParameters: pathParameters,
	}

	payload, err := json.Marshal(p)
	if err != nil {
		return nil, err
	}

	return payload, nil
}
