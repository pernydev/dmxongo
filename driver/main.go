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

var failPattern = [][]int{
	{255, 0, 0, 200, 0},
	{255, 0, 0, 0, 0},
}

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
		fmt.Println("Dial error:", err)
		Fail(dmx)
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

	fmt.Println("### opened ###")
	<-done

	// display fail pattern
	Fail(dmx)

	fmt.Println("### closed ###")
}

func Fail(dmx *dmx.DMX) {
	for i := 1; i <= 3; i++ {
		for _, frame := range failPattern {
			for i, channel := range frame {
				dmx.SetData(i+1, byte(channel))
			}
			dmx.Send()
			time.Sleep(500 * time.Millisecond)
		}
	}
	main()
}
