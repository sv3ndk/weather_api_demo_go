package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

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

func queryActiveSessionIds() ([]string, error) {

	expr, err := expression.NewBuilder().
		WithKeyCondition(
			expression.Key("PK").Equal(expression.Value("WS_SESSIONS")),
		).
		Build()

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error while building DynamoDB query: %w", err)
	}

	queryResult, err := dynamodbClient.Query(
		context.TODO(),
		&dynamodb.QueryInput{
			TableName:                 dynamoTable,
			KeyConditionExpression:    expr.KeyCondition(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error while querying DynamodDB: %w", err)
	}

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

// cf https://github.com/aws/aws-lambda-go/blob/main/events/README_DynamoDB.md
func handler(ctx context.Context, event events.DynamoDBEvent) {

	connectionIds, err := queryActiveSessionIds()
	if err != nil {
		log.Fatal("Could not fetch active ws connections from DB", err)
	}

	if len(connectionIds) > 0 {
		bytesEvents := [][]byte{}
		for _, record := range event.Records {
			// this is a bit lame, though I can't find a better approach
			// dynamodb basic data types seems to be quite duplicated in a bunch of incompatible packages :(
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

			eventBytes, err := json.Marshal(cleanEvent)
			if err != nil {
				log.Println("failed to process DynamoDB event", err)
			} else {
				bytesEvents = append(bytesEvents, eventBytes)
			}
		}

		for _, connectionId := range connectionIds {
			log.Println("sending records to active ws connection: ", connectionId)
			for _, event := range bytesEvents {
				log.Println("sending event", string(event), " to ", connectionId)
				outp, err := apiGWManagementClient.PostToConnection(
					ctx,
					&apigatewaymanagementapi.PostToConnectionInput{
						ConnectionId: &connectionId,
						Data:         event,
					},
				)

				if err != nil {
					log.Println("failed to send event ", err)
				} else {
					log.Println("send result ", *outp)
				}

			}
		}
	} else {
		log.Println("no WS client connected atm")
	}
}

func main() {
	lambda.Start(handler)
}
