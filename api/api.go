package api

import (
	"dmxongo/objects"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"dmxongo/fixtureTypes"
	"net/http"

	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173"
	},
}

var universeLatest *objects.Universe
var fixturesLatest []map[string]interface{}

var universeClients = make(map[string]*websocket.Conn)
var fixtureClients = make(map[string]*websocket.Conn)
var statsClients = make(map[string]*websocket.Conn)

var socketMutex = &sync.Mutex{}

var stats = make(map[string]interface{})

func HTTPAPI(universe *objects.Universe, fixtures []fixtureTypes.PAR) {
	stats["updateSpeeds"] = make([]float64, 0)

	universeLatest = universe
	JSONfixtures := make([]map[string]interface{}, len(fixtures))
	for i, fixture := range fixtures {
		JSONfixtures[i] = fixture.JSON()
	}
	fixturesLatest = JSONfixtures

	r := gin.Default()
	r.GET("/ws/universe", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		id := uuid.New().String()
		universeClients[id] = ws
		ws.WriteJSON(universeLatest)
		ws.SetCloseHandler(func(code int, text string) error {
			delete(universeClients, id)
			return nil
		})
	})
	r.GET("/ws/fixtures", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		id := uuid.New().String()
		fixtureClients[id] = ws
		ws.WriteJSON(fixturesLatest)
		ws.SetCloseHandler(func(code int, text string) error {
			delete(fixtureClients, id)
			return nil
		})
	})
	r.GET("/ws/stats", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		id := uuid.New().String()
		statsClients[id] = ws

		ws.WriteJSON(stats)

		ws.SetCloseHandler(func(code int, text string) error {
			delete(statsClients, id)
			return nil
		})
	})

	r.GET("/universe", func(c *gin.Context) {
		c.JSON(200, universeLatest)
	})
	r.GET("/fixtures", func(c *gin.Context) {
		c.JSON(200, fixturesLatest)
	})
	r.GET("/stats", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(200, stats)
	})

	r.Run(":8080")
}
func UniverseChanged(universe *objects.Universe) {
	universeLatest = universe
	socketMutex.Lock()
	defer socketMutex.Unlock()

	for _, client := range universeClients {
		client.WriteJSON(universe)
	}
}

func FixturesChanged(fixtures []fixtureTypes.PAR) {
	JSONfixtures := make([]map[string]interface{}, len(fixtures))
	for i, fixture := range fixtures {
		JSONfixtures[i] = fixture.JSON()
	}

	fixturesLatest = JSONfixtures

	socketMutex.Lock()
	defer socketMutex.Unlock()

	for _, client := range fixtureClients {
		client.WriteJSON(fixturesLatest)
	}
}

func SendUpdateSpeed(timeItTook time.Duration) {
	// if there are more than 50 values in stats["updateSpeeds"], remove the first one
	if len(stats["updateSpeeds"].([]float64)) > 50 {
		stats["updateSpeeds"] = stats["updateSpeeds"].([]float64)[1:]
	}

	stats["updateSpeeds"] = append(stats["updateSpeeds"].([]float64), float64(timeItTook.Nanoseconds()))

	socketMutex.Lock()
	defer socketMutex.Unlock()

	for _, client := range statsClients {
		client.WriteJSON(stats)
	}
}
