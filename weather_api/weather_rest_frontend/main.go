// Lambda serving the REST GET requests received from the API Gateway
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var dynamoClient *dynamodb.Client
var dynamoTable *string

const iso8601Tormat = "2006-01-02T15:04:05-0700"

func init() {
	dynamoTable = aws.String(os.Getenv("DYNAMO_TABLE"))

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	dynamoClient = dynamodb.NewFromConfig(awsCfg)
}

type InputParams struct {
	DeviceId int64
	FromTime time.Time
	ToTime   time.Time
}

type WeatherEvent struct {
	DeviceId  int64
	Time      time.Time
	EventType string
	Value     float64
}

// parseParams parses a URL encoded query string
// example input: '?device_id=1&from=2024-02-17T20:13:25+0100&to=2024-02-17T20:13:55+0100'
func parseParams(params map[string]string) (InputParams, error) {
	deviceIdStr, ok1 := params["device_id"]
	fromTimeIso, ok2 := params["from"]
	toTimeIso, ok3 := params["to"]
	if !(ok1 && ok2 && ok3) {
		return InputParams{}, errors.New("missing or invalid query params")
	}

	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		return InputParams{}, errors.New("missing or invalid device_id param")
	}

	fromTime, err1 := time.Parse(iso8601Tormat, fromTimeIso)
	toTime, err2 := time.Parse(iso8601Tormat, toTimeIso)

	if err1 != nil || err2 != nil {
		log.Println(err1, err2)
		return InputParams{}, errors.New("missing or invalid query params")
	}

	return InputParams{
		DeviceId: int64(deviceId),
		FromTime: fromTime,
		ToTime:   toTime,
	}, nil
}

func queryDb(inputParams InputParams) ([]WeatherEvent, error) {
	log.Printf("will use dynamo table %s and query params %v\n", *dynamoTable, inputParams)

	expr, err := expression.NewBuilder().
		WithKeyCondition(
			expression.KeyAnd(
				expression.Key("PK").Equal(expression.Value(fmt.Sprintf("DeviceId#%d", inputParams.DeviceId))),
				expression.Key("SK").Between(
					expression.Value(fmt.Sprintf("Time#%d", inputParams.FromTime.Unix()-1)),
					expression.Value(fmt.Sprintf("Time#%d", inputParams.ToTime.Unix()+1)),
				),
			),
		).
		Build()

	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("error while building DynamoDB query: %w", err)
	}

	queryResult, err := dynamoClient.Query(
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

	events := []WeatherEvent{}
	for _, rawEvent := range queryResult.Items {
		event := WeatherEvent{}
		attributevalue.UnmarshalMap(rawEvent, &event)
		events = append(events, event)
	}

	return events, nil
}

func serverSideError() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       "failed to fetch event from db",
		StatusCode: 500,
	}
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	inputParams, err := parseParams(request.QueryStringParameters)
	if err != nil {
		log.Println(err)
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	weatherEvents, err := queryDb(inputParams)
	if err != nil {
		log.Println(err)
		return serverSideError(), nil
	}
	log.Printf("returning %d events", len(weatherEvents))

	var body string
	if jsonBytes, err := json.Marshal(weatherEvents); err != nil {
		log.Println(err)
		return serverSideError(), nil
	} else if jsonBytes == nil {
		body = "{}"
	} else {
		body = string(jsonBytes)
	}

	return events.APIGatewayProxyResponse{
		Body:       body,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
