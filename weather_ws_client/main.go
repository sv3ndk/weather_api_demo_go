package main

import (
	"context"
	"flag"
	"log"
	"time"

	"nhooyr.io/websocket"
)

func main() {
	apiUrl := flag.String("url", "", "URL of the REST endpoint")
	flag.Parse()
	if len(*apiUrl) == 0 {
		flag.Usage()
	}

	log.Println("Querying WS service at ", *apiUrl)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	c, _, err := websocket.Dial(ctx, *apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer c.CloseNow()

	// // err = wsjson.Write(ctx, c, "hi")
	// // if err != nil {
	// // 	// ...
	// // }
    
    time.Sleep(10 * time.Second)

	c.Close(websocket.StatusNormalClosure, "")
}
