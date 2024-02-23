package main

import (
	"flag"
	"log"
	"time"

	"weather_rest_client/weather_client"
)

func main() {

	deviceId := flag.Int("deviceId", -1, "Id of the device")
	timeDelta := flag.Int("timeDelta", -1, "Duration in minutes of the queried period, ending now")
	apiUrl := flag.String("url", "", "URL of the REST endpoint")
	flag.Parse()
	if *deviceId == -1 || *timeDelta == -1 || len(*apiUrl) == 0 {
		flag.Usage()
	}
	fromTime := time.Now().Add(-10 * time.Minute)
	toTime := time.Now()

	client := weather_client.New(*apiUrl)

	events, err := client.QueryEvents(*deviceId, fromTime, toTime)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n\nlast %d minutes of events of device %d:\n\n", *timeDelta, *deviceId)
	for _, event := range events {
		log.Println(event)
	}
}
