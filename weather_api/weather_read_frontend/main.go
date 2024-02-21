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

func init() {
	dynamoTable = aws.String(os.Getenv("DYNAMO_TABLE"))

	awsCfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}

	dynamoClient = dynamodb.NewFromConfig(awsCfg)
}

// prints any stuct. This is convenient since it dereference all pointers
func logAny(v any) {
	j, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(j))
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

// expects a query string like ?device_id=1&from=2024-02-17T20:13:25%2B0100&to=2024-02-17T20:13:55%2B0100'
// ('%2B' is the URL encoding of '+')
func parseParams(params map[string]string) (InputParams, error) {
	deviceIdStr, ok1 := params["device_id"]
	fromTimeIso, ok2 := params["from"]
	toTimeIso, ok3 := params["to"]
	if !(ok1 && ok2 && ok3) {
		return InputParams{}, errors.New("missing or invalid query params")
	}

	deviceId, err := strconv.Atoi(deviceIdStr)
	if err != nil {
		return InputParams{}, errors.New("missing or invalid query params")
	}

	iso8601Tormat := "2006-01-02T15:04:05-0700"

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
	log.Printf("will use dynamo table %s and query params %s\n", *dynamoTable, inputParams)

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
		return nil, err
	}

	log.Printf("key condition:")
	logAny(expr.KeyCondition())
	logAny(expr.Names())
	logAny(expr.Values())

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
		fmt.Println(err)
		return nil, err
	}

	events := []WeatherEvent{}
	for _, rawEvent := range queryResult.Items {
		event := WeatherEvent{}
		attributevalue.UnmarshalMap(rawEvent, &event)
		events = append(events, event)
	}

	return events, nil
}

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	inputParams, err := parseParams(request.QueryStringParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}

	weatherEvents, err := queryDb(inputParams)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}
	log.Println("returning events", weatherEvents)

	jsonBytes, err := json.Marshal(weatherEvents)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, nil
	}

	var body string
	if jsonBytes == nil {
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
