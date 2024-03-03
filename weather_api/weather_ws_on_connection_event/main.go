// Lambda keeping track of the connection id of the currently connected websocket clients.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamodbClient *dynamodb.Client
var dynamoTable *string

const SESSION_PK string = "WS_SESSIONS"

func init() {
	dynamoTable = aws.String(os.Getenv("DYNAMO_TABLE"))

	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("Could not connect to dynamo ", err)
	}
	dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
}

// handleRequest is triggered by the API Gateway any time a ws client connects or disconnects
func handleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	if request.RequestContext.RouteKey == "$connect" {
		log.Printf("new connection with id %s\n", request.RequestContext.ConnectionID)
		if err := storeConnectionID(ctx, request.RequestContext.ConnectionID); err != nil {
			log.Println(err)
			return serverError("could not persist connection id"), err
		}
	} else if request.RequestContext.RouteKey == "$disconnect" {
		log.Printf("connection id %s is now stopped \n", request.RequestContext.ConnectionID)
		if err := removeConnectionID(ctx, request.RequestContext.ConnectionID); err != nil {
			log.Println(err)
			return serverError("could not clean up connection id"), err
		}
	} else {
		log.Println("unexpected route key", request.RequestContext.RouteKey)
		return serverError(""), nil
	}

	return events.APIGatewayProxyResponse{
			Body:       "",
			StatusCode: 200,
		},
		nil
}

func storeConnectionID(ctx context.Context, connectionId string) error {
	putItem := dynamodb.PutItemInput{
		TableName: dynamoTable,
		Item: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: SESSION_PK,
			},
			"SK": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("Id#%s", connectionId),
			},
			"ConnectionId": &types.AttributeValueMemberS{
				Value: connectionId,
			},
		},
	}

	if _, err := dynamodbClient.PutItem(ctx, &putItem); err != nil {
		return fmt.Errorf("error while inserting connection id %s in DyanmoDB: %w", connectionId, err)
	}
	return nil
}

func removeConnectionID(ctx context.Context, connectionId string) error {
	deletedItem := dynamodb.DeleteItemInput{
		TableName: dynamoTable,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{
				Value: SESSION_PK,
			},
			"SK": &types.AttributeValueMemberS{
				Value: fmt.Sprintf("Id#%s", connectionId),
			},
		},
	}

	if _, err := dynamodbClient.DeleteItem(ctx, &deletedItem); err != nil {
		return fmt.Errorf("error while removing sessiongId %s: %w", connectionId, err)
	}
	return nil
}

func serverError(msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       msg,
		StatusCode: 500,
	}
}

func main() {
	lambda.Start(handleRequest)
}
