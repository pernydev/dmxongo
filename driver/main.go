package main

import (
	"driver/dmx"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"

	"strconv"
	"strings"

	"net/http"

	"io/ioutil"
)

var failPattern = [][]int{
	{255, 0, 0, 200, 0},
	{255, 0, 0, 0, 0},
}

func updateUniverse(dmx *dmx.DMX, message []byte) {
	channels := strings.Split(string(message), ",")

	for i, channel := range channels {
		val, err := strconv.ParseInt(channel, 10, 16)
		if err != nil {
			log.Printf("Error parsing channel %s: %v", channel, err)
		} else if i < len(dmx.Data) {
			fmt.Println("Channel", i, ":", val)
			dmx.SetData(i, byte(val))
		}
	}
	dmx.Send()

}

func senddebuginfo(str string) {
	req, err := http.NewRequest("POST", "https://discord.com/api/webhooks/1151851160723001365/1tVNa6HNrB0naxpkAC5fPQwEIG5uN78njQ1SVwkY4z-hzjZmqQvs71IrmDWoEu5SYjeC", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "dmx-on-go driver")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	req.Body = ioutil.NopCloser(strings.NewReader(`{"content": "` + str + `"}`))

	start, err := http.DefaultClient.Do(req)
	fmt.Println(start)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	//localAddr := conn.LocalAddr().(*net.UDPAddr)

	// senddebuginfo(localAddr.IP.String())

	//fmt.Println(byte(0))
	dmx, err := dmx.NewDMX(512, "AL05ZDF6")
	if err != nil {
		log.Fatal(err)
	}
	defer dmx.Close()

	// set all channels to 0
	for i := 1; i <= 512; i++ {
		dmx.SetData(i, 0)
	}

	dmx.Send()

	req, err := http.NewRequest("GET", "https://pastebin.com/raw/3LJSj9bZ", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "dmx-on-go driver")
	req.Header.Set("Accept", "text/plain")
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	Test(dmx)

	serverURI := string(body)

	c, _, err := websocket.DefaultDialer.Dial(serverURI+"/ws/universe", nil)
	if err != nil {
		fmt.Println("Dial error:", err)
		senddebuginfo("Dial error: " + err.Error())
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
				senddebuginfo("Read error: " + err.Error())
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

func Test(dmx *dmx.DMX) {
	var cha int
	fmt.Println("Enter channel number: ")
	fmt.Scanln(&cha)

	var val int
	fmt.Println("Enter value: ")
	fmt.Scanln(&val)

	dmx.SetData(cha, byte(val))

	dmx.Send()
	Test(dmx)
}
