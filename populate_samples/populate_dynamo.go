package main

// Simple script to insert some sample data into Dynamodb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

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

func addSample(ctx context.Context, client *dynamodb.Client, event WeatherEvent, iteration int) error {

	// mutating this copy of event to create some variation in the data
	event.Time = time.Unix(event.Time.Unix()+int64(iteration*10), 0)
	event.Value = event.Value + float64(iteration)
	logAny(event)

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
	}
	logAny(asMap)

	output, err := client.PutItem(
		ctx,
		&dynamodb.PutItemInput{
			Item:      asMap,
			TableName: aws.String("weather-api-demo-WeatherDynamoTable-SUD97CERSXYR"),
		},
	)
	if err != nil {
		return err
	}

	log.Println("insertion result:")
	logAny(output)
	return nil
}

func main() {
	log.Println("Populating DynamoDB with dummy weather data")
	ctx := context.Background()

	sdkConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Could not connect to dynamo ", err)
	}
	dynamoClient := dynamodb.NewFromConfig(sdkConfig)

	some_events := []WeatherEvent{
		{DeviceId: 123, Time: time.Unix(1708176000, 0), EventType: "Pressure", Value: 1037},
		{DeviceId: 123, Time: time.Unix(1708176000, 0), EventType: "Temperature", Value: 20.5},
		{DeviceId: 123, Time: time.Unix(1708176000, 0), EventType: "Humidity", Value: 0.5},
		{DeviceId: 123, Time: time.Unix(1708176000, 0), EventType: "WindSpeed", Value: 5.5},
		{DeviceId: 123, Time: time.Unix(1708176000, 0), EventType: "WindDirection", Value: 180},

		{DeviceId: 124, Time: time.Unix(1708176000, 0), EventType: "Pressure", Value: 1037},
		{DeviceId: 124, Time: time.Unix(1708176000, 0), EventType: "Temperature", Value: 20.5},
		{DeviceId: 124, Time: time.Unix(1708176000, 0), EventType: "Humidity", Value: 0.5},
		{DeviceId: 124, Time: time.Unix(1708176000, 0), EventType: "WindSpeed", Value: 5.5},
		{DeviceId: 124, Time: time.Unix(1708176000, 0), EventType: "WindDirection", Value: 180},
	}

	for _, event := range some_events {
		for i := range 3 {
			err = addSample(ctx, dynamoClient, event, i)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
