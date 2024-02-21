package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

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

// I've copy/pasted this everywhere because I have no shame :)
type WeatherEvent struct {
	DeviceId  int64
	Time      time.Time
	EventType string
	Value     float64
}

// prints any stuct. This is convenient since it dereference all pointers
func logAny(v any) {
	j, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(j))
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

	log.Printf("Inserting event into DynamoDB: %+v\n", event)

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
	logAny(asMap)

	output, err := client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			Item:      asMap,
			TableName: dynamoTable,
		},
	)
	if err != nil {
		return err
	}

	log.Println("insertion result:")
	logAny(output)
	return nil
}

// TODO: move this to a lambda triggered from Event bridge
func main() {
	log.Println("generating random weather event")

	events := []WeatherEvent{}

	for i := range 10 {
		deviceId := int64(1000 + i)
		events = append(events, randomEvents(deviceId)...)
	}

	log.Println("Created weather event")
	for _, weatherEvent := range events {
		err := addSample(ctx, dynamodbClient, weatherEvent)
		if err != nil {
			log.Fatal("Could not insert event into Dynamodb", err)
		}
	}
}
