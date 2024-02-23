package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var dynamodbClient *dynamodb.Client
var dynamoTable *string
var ctx context.Context = context.Background()

func init() {
	dynamoTable = aws.String(os.Getenv("DYNAMO_TABLE"))

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Could not connect to dynamo ", err)
	}
	dynamodbClient = dynamodb.NewFromConfig(sdkConfig)
}

type WeatherEvent struct {
	DeviceId  int64
	Time      time.Time
	EventType string
	Value     float64
}

func randomPressureEvent(deviceId int64) WeatherEvent {
	return WeatherEvent{
		DeviceId:  deviceId,
		Time:      time.Now(),
		EventType: "Pressure",
		Value:     float64(rand.Int31n(100) + 950),
	}
}

func randomTemperatureEvent(deviceId int64) WeatherEvent {
	return WeatherEvent{
		DeviceId:  deviceId,
		Time:      time.Now(),
		EventType: "Temperature",
		Value:     rand.Float64()*40 - 10,
	}
}

func randomHumidityEvent(deviceId int64) WeatherEvent {
	return WeatherEvent{
		DeviceId:  deviceId,
		Time:      time.Now(),
		EventType: "Humidity",
		Value:     rand.Float64() * 100,
	}
}

func randomWindSpeedEvent(deviceId int64) WeatherEvent {
	return WeatherEvent{
		DeviceId:  deviceId,
		Time:      time.Now(),
		EventType: "WindSpeed",
		Value:     float64(rand.Int31n(50)),
	}
}

func randomWindDirectionEvent(deviceId int64) WeatherEvent {
	return WeatherEvent{
		DeviceId:  deviceId,
		Time:      time.Now(),
		EventType: "WindDirection",
		Value:     rand.Float64() * 360,
	}
}

func randomEvents(deviceId int64) []WeatherEvent {
	return []WeatherEvent{
		randomPressureEvent(deviceId),
		randomTemperatureEvent(deviceId),
		randomHumidityEvent(deviceId),
		randomWindSpeedEvent(deviceId),
		randomWindDirectionEvent(deviceId),
	}
}

func addSample(ctx context.Context, client *dynamodb.Client, event WeatherEvent) error {
	// There is also attributevalue.MarshalMap(event) in aws-sdk-go-v2/feature/dynamodb/attributevalue but
	// which reduce the boiler plate, but needs to be completed with "PK", "SK" and potentially
	// adjusted for any DynamoDB expression attribute name
	asMap := map[string]types.AttributeValue{
		// we need to pass a pointer here since isAttributeValue() is attached to a pointer receiver
		"PK": &types.AttributeValueMemberS{
			Value: fmt.Sprintf("DeviceId#%d", event.DeviceId),
		},
		"SK": &types.AttributeValueMemberS{
			// considering there can be maximum one event of any type for a given device at any timestamp
			Value: fmt.Sprintf("Time#%d#Type%s", event.Time.Unix(), event.EventType),
		},
		"DeviceId": &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", event.DeviceId),
		},
		"EventType": &types.AttributeValueMemberS{
			Value: event.EventType,
		},
		"Value": &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%f", event.Value),
		},
		"Time": &types.AttributeValueMemberN{
			Value: fmt.Sprintf("%d", event.Time.Unix()),
		},
	}

	_, err := client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			Item:      asMap,
			TableName: dynamoTable,
		},
	)
	if err != nil {
		return fmt.Errorf("error while inserting event %v in DyanmoDB %w", event, err)
	}

	return nil
}

func genData() error {
	log.Println("generating random weather event")

	events := []WeatherEvent{}
	for i := range 10 {
		deviceId := int64(1000 + i)
		events = append(events, randomEvents(deviceId)...)
	}

	for _, weatherEvent := range events {
		err := addSample(ctx, dynamodbClient, weatherEvent)
		if err != nil {
			return err
		}
	}

	log.Println("done")
	return nil
}

// cf https://github.com/aws/aws-lambda-go/blob/main/events/README_EventBridge_Events.md
func handler(ctx context.Context, request events.EventBridgeEvent) {
	err := genData()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	lambda.Start(handler)
}
