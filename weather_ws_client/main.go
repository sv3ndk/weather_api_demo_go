package main

import (
	"context"
	"log"
	"time"

	"nhooyr.io/websocket"
)

func main() {
	apiUrl := "wss://pfd41jd2b8.execute-api.eu-central-1.amazonaws.com/v1"
	log.Println("Querying WS service at ", apiUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.CloseNow()

	// // err = wsjson.Write(ctx, c, "hi")
	// // if err != nil {
	// // 	// ...
	// // }

	c.Close(websocket.StatusNormalClosure, "")
}
