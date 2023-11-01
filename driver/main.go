package main

import (
	"driver/dmx"
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func updateUniverse(dmx *dmx.DMX, message []byte) {
	start := time.Now()
	dmx.Data = json.RawMessage(message)
	dmx.Send()
	log.Printf("Time to send: %s", time.Since(start))
}

func main() {
	dmx, err := dmx.NewDMX(512, "AL05ZDF6")
	if err != nil {
		log.Fatal(err)
	}
	defer dmx.Close()

	url := "ws://127.0.0.1:8080/ws/universe"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Dial error:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			updateUniverse(dmx, message)
		}
	}()

	log.Println("### opened ###")
	<-done
	log.Println("### closed ###")
}
