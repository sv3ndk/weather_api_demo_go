package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type InputParams struct {
	DeviceId string
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
	deviceId, ok1 := params["device_id"]
	fromTimeIso, ok2 := params["from"]
	toTimeIso, ok3 := params["to"]
	if !(ok1 && ok2 && ok3) {
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
		DeviceId: deviceId,
		FromTime: fromTime,
		ToTime:   toTime,
	}, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	inputParams, err := parseParams(request.QueryStringParameters)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 400,
		}, nil
	}
	log.Println("input Params", inputParams)

	var greeting string
	sourceIP := request.RequestContext.Identity.SourceIP
	if sourceIP == "" {
		greeting = "Hello, world!\n"
	} else {
		greeting = fmt.Sprintf("Hello, %s!\n", sourceIP)
	}

	return events.APIGatewayProxyResponse{
		Body:       greeting,
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
