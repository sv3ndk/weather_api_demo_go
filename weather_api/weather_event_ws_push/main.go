// Lambda listening to new weather events from DynamoDB stream and forwarding
// them in JSON format to all currently connected websocket clients.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamodbClient *dynamodb.Client
var dynamoTable *string
var apiGWManagementClient *apigatewaymanagementapi.Client

func init() {
	sdkConfig, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal("Could not connect to AWS API", err)
	}

	dynamoTable = aws.String(os.Getenv("DYNAMO_TABLE"))
	dynamodbClient = dynamodb.NewFromConfig(sdkConfig)

	wsClientCallbackUrl := fmt.Sprintf(
		"https://%s.execute-api.%s.amazonaws.com/%s",
		os.Getenv("API_ID"),
		os.Getenv("AWS_REGION"),
		os.Getenv("API_STAGE"),
	)
	apiGWManagementClient = apigatewaymanagementapi.NewFromConfig(
		sdkConfig,
		func(o *apigatewaymanagementapi.Options) {
			o.BaseEndpoint = &wsClientCallbackUrl
		},
	)
}

// queryActiveSessionIds retrieves the list of connection id of currently connected ws clients
func queryActiveSessionIds(ctx context.Context) ([]string, error) {
	expr, err := expression.NewBuilder().
		WithKeyCondition(
			expression.Key("PK").Equal(expression.Value("WS_SESSIONS")),
		).
		Build()

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error while building DynamoDB query: %w", err)
	}

	query := dynamodb.QueryInput{
		TableName:                 dynamoTable,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	if queryResult, err := dynamodbClient.Query(ctx, &query); err != nil {
		return nil, fmt.Errorf("error while querying DynamodDB: %w", err)
	} else {
		connectionIds := []string{}
		for _, rawEvent := range queryResult.Items {
			var connectionIdValue string
			if connectionId, ok := rawEvent["ConnectionId"]; ok {
				if err := attributevalue.Unmarshal(connectionId, &connectionIdValue); err != nil {
					log.Printf("failed to parse %v, skipping %v", connectionId, err)
				} else {
					connectionIds = append(connectionIds, connectionIdValue)
				}
			}
		}
		return connectionIds, nil
	}
}

// sendEventsToWsClients tries to forward the specified events to those websocket connection ID
// via the API gateway. Any error is just ignored
func sendEventsToWsClients(ctx context.Context, weatherEvents [][]byte, connectionIds []string) {
	var waiter sync.WaitGroup
	for _, connectionId := range connectionIds {
		log.Println("sending records to active ws connection: ", connectionId)
		for _, event := range weatherEvents {
			log.Println("sending event", string(event), " to ", connectionId)
			waiter.Add(1)
			go func() {
				defer waiter.Done()
				postInput := apigatewaymanagementapi.PostToConnectionInput{
					ConnectionId: &connectionId,
					Data:         event,
				}
				if _, err := apiGWManagementClient.PostToConnection(ctx, &postInput); err != nil {
					log.Println("failed to send event ", err)
				}

			}()
		}
	}
	waiter.Wait()
}

func handler(ctx context.Context, event events.DynamoDBEvent) {

	timeBoxedCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	connectionIds, err := queryActiveSessionIds(timeBoxedCtx)
	if err != nil {
		log.Fatal("Could not fetch active ws connections from DB", err)
	}

	if len(connectionIds) > 0 {
		weatherEvents := [][]byte{}
		for _, record := range event.Records {
			cleanEvent := map[string]any{}
			for k, v := range record.Change.NewImage {
				if k != "PK" && k != "SK" {
					if v.DataType() == events.DataTypeString {
						cleanEvent[k] = v.String()
					} else if v.DataType() == events.DataTypeNumber {
						cleanEvent[k] = v.Number()
					}
				}
			}
			if eventBytes, err := json.Marshal(cleanEvent); err != nil {
				log.Println("failed to process DynamoDB event", err)
			} else {
				weatherEvents = append(weatherEvents, eventBytes)
			}
		}

		sendEventsToWsClients(timeBoxedCtx, weatherEvents, connectionIds)

	} else {
		log.Println("no WS client connected atm")
	}
}

func main() {
	lambda.Start(handler)
}
