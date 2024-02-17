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
			TableName: aws.String("weather-api-demo-WeatherDynamoTable-1N9WZZRBZNKLM"),
		},
	)
	if err != nil {
		return err
	}

	log.Println("insertion result:")
	logAny(output)
	return nil
}

func mainn() {
	log.Println("Populating DynamoDB with dummy weather data")
	ctx := context.Background()

	sdk_config, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal("Could not connect to dynamo ", err)
	}
	dynamoClient := dynamodb.NewFromConfig(sdk_config)

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


func main() {
	iso8601Timestamp := "2024-02-17T20:13:55+0000"

    layout := "2006-01-02T15:04:05-0700"

	// Parse ISO 8601 timestamp
	parsedTime, err := time.Parse(layout, iso8601Timestamp)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}

	fmt.Println("Parsed time:", parsedTime)
}
