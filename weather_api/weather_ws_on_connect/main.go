package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	fmt.Printf("new connection with id  %s\n", request.RequestContext.ConnectionID)
	return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 200,
		},
		nil
}

func main() {
	lambda.Start(handleRequest)
}
