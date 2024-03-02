// Lambda writing random weather events to Dynamodb.
// Meant to be triggered every minute by an EventBridge scheduler.
package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
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

// randomEvents creates one random weather event of each type for the given deviceID
func randomEvents(deviceId int64) []WeatherEvent {
	return []WeatherEvent{
		randomPressureEvent(deviceId),
		randomTemperatureEvent(deviceId),
		randomHumidityEvent(deviceId),
		randomWindSpeedEvent(deviceId),
		randomWindDirectionEvent(deviceId),
	}
}

// addSamples inserts the given weather events into DynamoDB as one single batch
func addSamples(ctx context.Context, weatherEvents []WeatherEvent) error {
	log.Println("inserting batch")

	if len(weatherEvents) > 0 && len(weatherEvents) < 26 {

		putRequests := make([]types.WriteRequest, 0, len(weatherEvents))

		for _, weatherEvent := range weatherEvents {
			putRequest := types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: map[string]types.AttributeValue{
						"PK": &types.AttributeValueMemberS{
							Value: fmt.Sprintf("DeviceId#%d", weatherEvent.DeviceId),
						},
						"SK": &types.AttributeValueMemberS{
							Value: fmt.Sprintf("Time#%d#Type%s", weatherEvent.Time.Unix(), weatherEvent.EventType),
						},
						"DeviceId": &types.AttributeValueMemberN{
							Value: fmt.Sprintf("%d", weatherEvent.DeviceId),
						},
						"EventType": &types.AttributeValueMemberS{
							Value: weatherEvent.EventType,
						},
						"Value": &types.AttributeValueMemberN{
							Value: fmt.Sprintf("%f", weatherEvent.Value),
						},
						"Time": &types.AttributeValueMemberN{
							Value: fmt.Sprintf("%d", weatherEvent.Time.Unix()),
						},
					},
				},
			}
			putRequests = append(putRequests, putRequest)
		}

		input := dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				*dynamoTable: putRequests,
			},
		}

		if _, err := dynamodbClient.BatchWriteItem(ctx, &input); err != nil {
			return fmt.Errorf("error while inserting events in DyanmoDB %w", err)
		}

	} else {
		log.Printf("refusing to insert a batch of size %d", len(weatherEvents))
	}

	return nil
}

// addAllSamples slices the given array into batches of 25 (i.e. the maximum allowed
// by DynamoDB) and sends them to addSamples.
// (in theory we should check if keys overlap, although here we know they never do)
func addAllSamples(ctx context.Context, weatherEvents []WeatherEvent) {
	log.Println("sending generated data to DB")
	var waiter = sync.WaitGroup{}
	for i := 0; i < len(weatherEvents); i += 25 {
		fromIdx := i
		toIdx := min(i+25, len(weatherEvents))
		waiter.Add(1)
		go func() {
			defer waiter.Done()
			if err := addSamples(ctx, weatherEvents[fromIdx:toIdx]); err != nil {
				log.Println("failed to insert data in Dynamo", err)
			}
		}()
	}
	waiter.Wait()
}

func handler(ctx context.Context, request events.EventBridgeEvent) {
	log.Println("generating random weather event")
	events := make([]WeatherEvent, 0, 50)
	for i := range 10 {
		deviceId := int64(1000 + i)
		events = append(events, randomEvents(deviceId)...)
	}
	addAllSamples(ctx, events)
	log.Println("done")
}

func main() {
	lambda.Start(handler)
}
