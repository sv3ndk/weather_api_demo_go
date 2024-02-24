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

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamodbClient *dynamodb.Client
var dynamoTable *string
var apiID *string
var apiStage *string
var awsRegion *string
var ctx context.Context = context.Background()

func init() {
	dynamoTable = aws.String(os.Getenv("DYNAMO_TABLE"))
	apiID = aws.String(os.Getenv("API_ID"))
	apiStage = aws.String(os.Getenv("API_STAGE"))
	awsRegion = aws.String(os.Getenv("AWS_REGION"))

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Could not connect to dynamo ", err)
	}
	dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
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
func handler(ctx context.Context, e events.DynamoDBEvent) {

	connectionIds, err := queryActiveSessionIds()
	if err != nil {
		log.Fatal("Could not fetch active ws connections from DB", err)
	}

	if len(connectionIds) > 0 {
		for _, connectionId := range connectionIds {
			log.Println("active connections: ", connectionId)
			clientCallbackUrl := fmt.Sprintf("POST https://%s.execute-api.%s.amazonaws.com/%s/@connections/%s", *apiID, *awsRegion, *apiStage, connectionId)
			log.Printf("Should now post back to %s", clientCallbackUrl)
		}
	} else {
		log.Printf("no client is connected atm")
	}

	// for _, record := range e.Records {
	// 	log.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

	// 	for name, value := range record.Change.NewImage {
	// 		if value.DataType() == events.DataTypeString {
	// 			log.Printf("Attribute name: %s, value: %s\n", name, value.String())
	// 		}
	// 	}
	// }
}

func main() {
	lambda.Start(handler)
}
