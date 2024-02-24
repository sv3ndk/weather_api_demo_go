package main

import (
	"context"
	"flag"
	"io"
	"log"

	"nhooyr.io/websocket"
)

func main() {
	apiUrl := flag.String("url", "", "URL of the REST endpoint")
	flag.Parse()
	if len(*apiUrl) == 0 {
		flag.Usage()
	}

	log.Println("Listening to WS service at ", *apiUrl)
	// ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	ctx := context.Background()
	// defer cancel()

	c, _, err := websocket.Dial(ctx, *apiUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		log.Println("closing socket now")
		c.Close(websocket.StatusNormalClosure, "")
	}()

	for {
		_, reader, err := c.Reader(ctx)
		if err != nil {
			log.Fatal(err)
		} else {
			data, err := io.ReadAll(reader)
			if err != nil {
				log.Fatal(err)
			} else {
				log.Println(string(data))
			}
		}
	}
}
