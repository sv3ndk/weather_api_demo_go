package weather_client

import (
	"crypto/tls"
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

type WeatherClient struct {
	ApiUrl     string
	httpClient *http.Client
	apiKey     string
}

func New(url, apiKey, certFile, keyFile string) WeatherClient {

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatal(err)
	}

	return WeatherClient{
		ApiUrl: url,
		httpClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					Certificates: []tls.Certificate{cert},
				},
			},
		},
		apiKey: apiKey,
	}
}

func (c WeatherClient) QueryEvents(deviceId int, fromTime time.Time, toTime time.Time) ([]WeatherEvent, error) {
	log.Printf("looking for weather events for device %v from %s to %v", deviceId, fromTime, toTime)

	req, err := http.NewRequest("GET", c.ApiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	iso8601Format := "2006-01-02T15:04:05-0700"

	q := url.Values{}
	q.Add("device_id", fmt.Sprint(deviceId))
	q.Add("from", fromTime.Format(iso8601Format))
	q.Add("to", toTime.Format(iso8601Format))
	req.URL.RawQuery = q.Encode()
	log.Printf("querying URL %s", req.URL.String())

	req.Header["X-API-Key"] = []string{c.apiKey}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Panicln(err)
		return nil, fmt.Errorf("could not query API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error response from server: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data []WeatherEvent
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse response body: %w", err)
	}

	return data, nil
}
