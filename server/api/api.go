package api

import (
	"dmxongo/objects"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"net/http"

	"dmxongo/functions"
	"strings"
	"sync"

	"dmxongo/events"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all connections
	},
}

var universe *objects.Universe
var fixtures *[]objects.Fixture

var universeClients = make(map[string]*websocket.Conn)
var fixtureClients = make(map[string]*websocket.Conn)
var statsClients = make(map[string]*websocket.Conn)

var socketMutex = &sync.Mutex{}

var functionsState = make(map[string]bool)

var stats = make(map[string]interface{})

func HTTPAPI(universePointer *objects.Universe, fixturesPointer *[]objects.Fixture) {
	stats["updateSpeeds"] = make([]float64, 0)

	events.UniverseChangeListeners = append(events.UniverseChangeListeners, UniverseChanged)
	events.FixtureChangeListeners = append(events.FixtureChangeListeners, FixturesChanged)

	universe = universePointer
	fixtures = fixturesPointer

	r := gin.Default()
	r.GET("/ws/universe", func(c *gin.Context) {
		ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		id := uuid.New().String()
		universeClients[id] = ws

		// universe.ChannelValues separated by a comma
		csvString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(universe.ChannelValues)), ","), "[]")
		ws.WriteMessage(websocket.TextMessage, []byte(csvString))
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

		JSONFixtures := make([]map[string]interface{}, len(*fixtures))
		for i, fixture := range *fixtures {
			JSONFixtures[i] = fixture.JSON()
		}

		ws.WriteJSON(JSONFixtures)

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
		c.JSON(200, universe)
	})
	r.GET("/fixtures", func(c *gin.Context) {
		JSONFixtures := make([]map[string]interface{}, len(*fixtures))
		for i, fixture := range *fixtures {
			JSONFixtures[i] = fixture.JSON()
		}

		c.JSON(200, JSONFixtures)
	})
	r.GET("/stats", func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(200, stats)
	})
	r.POST("/function/:function", func(c *gin.Context) {
		functionName := c.Param("function")
		// check if function is running
		if functionsState[functionName] {
			functions.Functions[functionName].Stop()
			delete(functionsState, functionName)
			c.JSON(200, gin.H{"status": "stopped"})
			return
		}

		// start function
		functions.Functions[functionName].Start()
		if functions.Functions[functionName].Type != "basic" {
			functionsState[functionName] = true
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(200, gin.H{"status": "started"})

	})

	r.GET("/function", func(c *gin.Context) {
		JSONObject := make(map[string]interface{})
		for functionName, function := range functions.Functions {
			JSONObject[functionName] = map[string]interface{}{
				"running": functionsState[functionName],
				"type":    function.Type,
			}
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.JSON(200, JSONObject)
	})

	r.Run(":8080")
}

func UniverseChanged() {
	socketMutex.Lock()
	defer socketMutex.Unlock()

	for _, client := range universeClients {
		// universe.ChannelValues separated by a comma
		csvString := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(universe.ChannelValues)), ","), "[]")
		client.WriteMessage(websocket.TextMessage, []byte(csvString))
	}
}

func FixturesChanged() {
	fmt.Println(*fixtures)
	JSONFixtures := make([]map[string]interface{}, len(*fixtures))
	for i, fixture := range *fixtures {
		JSONFixtures[i] = fixture.JSON()
	}

	socketMutex.Lock()
	defer socketMutex.Unlock()

	for _, client := range fixtureClients {
		client.WriteJSON(JSONFixtures)
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
