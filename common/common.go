package common

import (
	"github.com/aws/aws-lambda-go/events"
)

// type Response events.APIGatewayProxyResponse
// type Request events.APIGatewayProxyRequest

// Generic Http error response
func ErrorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: 500,
	}
}
