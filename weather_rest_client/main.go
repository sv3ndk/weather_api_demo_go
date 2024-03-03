// CLI app fetching fetching the weather events from the weather REST API
package main

import (
	"flag"
	"log"
	"time"

	"weather_rest_client/weather_client"
)

/*
	./weather_rest_client  \
		-url https://rest.weather-api-demo.poc.svend.xyz/weather  \
		-deviceId 1001 \
		-timeDelta 13 \
		-apiKey to_be_fetched_from_aws \
		-certFile certificates/clientCert.pem \
		-keyFile certificates/clientKey.pem
*/
func main() {
	deviceId := flag.Int("deviceId", -1, "Id of the device")
	timeDelta := flag.Int("timeDelta", -1, "Duration in minutes of the queried period, ending now")
	apiUrl := flag.String("url", "", "URL of the REST endpoint")
	apiKey := flag.String("apiKey", "", "API key")
	certFile := flag.String("certFile", "", "PEM file containing the client public certificate")
	keyFile := flag.String("keyFile", "", "PEM file containing the client private key")
	toTime := time.Now()
	fromTime := toTime.Add(-10 * time.Minute)

	if flag.Parse(); *deviceId == -1 || *timeDelta == -1 || len(*apiUrl) == 0 {
		flag.Usage()
	}

	client := weather_client.New(*apiUrl, *apiKey, *certFile, *keyFile)
	events, err := client.QueryEvents(*deviceId, fromTime, toTime)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("\n\nlast %d minutes of weather events of device %d:\n\n", *timeDelta, *deviceId)
	for _, event := range events {
		log.Println(event)
	}
}
