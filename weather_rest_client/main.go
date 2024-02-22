package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type WeatherEvent struct {
	DeviceId  int64
	Time      time.Time
	EventType string
	Value     float64
}

func main() {
	fromTime := time.Now().Add(-1 * time.Hour)
	toTime := time.Now()
	deviceId := 1003

	apiUrl := "https://rzwn13e2ak.execute-api.eu-central-1.amazonaws.com/Prod/weather"

	log.Printf("Looking for weather events for device %v from %s to %v", deviceId, fromTime, toTime)

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	iso8601Format := "2006-01-02T15:04:05-0700"

	q := url.Values{}
	q.Add("device_id", fmt.Sprint(deviceId))
	q.Add("from", fromTime.Format(iso8601Format))
	q.Add("to", toTime.Format(iso8601Format))

	req.URL.RawQuery = q.Encode()
	log.Println(req.URL.String())

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Error response from server: %s", resp.Status)
	}

	log.Println("got response code", resp.StatusCode)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data []WeatherEvent
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Found weather events")
	for _, event := range data {
		// logAny(event)
		log.Println(event)
	}

}
