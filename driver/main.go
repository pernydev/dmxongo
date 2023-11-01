package main

import (
	"driver/dmx"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"

	"strconv"
	"strings"
)

func updateUniverse(dmx *dmx.DMX, message []byte) {
	start := time.Now()
	channels := strings.Split(string(message), ",")
	fmt.Println(channels)

	for i, channel := range channels {
		val, err := strconv.ParseInt(channel, 10, 16)
		if err != nil {
			log.Printf("Error parsing channel %s: %v", channel, err)
		} else if i < len(dmx.Data) {
			dmx.Data[i] = byte(val)
		}
	}

	dmx.Send()
	log.Printf("Time to send: %s", time.Since(start))
}

func main() {
	fmt.Println(byte(0))
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
