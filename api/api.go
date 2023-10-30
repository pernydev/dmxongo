package api

import (
	"dmxongo/objects"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"

	"dmxongo/fixtureTypes"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var universeClients = make(map[uuid.UUID]*websocket.Conn)
var fixtureClients = make(map[uuid.UUID]*websocket.Conn)

var universeLatest *objects.Universe
var fixturesLatest []fixtureTypes.PAR

func API(universeCTX *objects.Universe, fixturesCTX []fixtureTypes.PAR) {
	universeLatest = universeCTX
	fixturesLatest = fixturesCTX

	router := gin.Default()
	router.GET("/listen/universe", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		clientUUID := uuid.New()
		universeClients[clientUUID] = conn

		conn.WriteJSON(universeLatest)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				delete(universeClients, clientUUID)
				return
			}
		}
	})
	router.GET("/listen/fixtures", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		clientUUID := uuid.New()
		universeClients[clientUUID] = conn

		JSONfixtures := make([]map[string]interface{}, len(fixturesLatest))
		for i, fixture := range fixturesLatest {
			JSONfixtures[i] = fixture.JSON()
		}

		conn.WriteJSON(JSONfixtures)

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				delete(universeClients, clientUUID)
				return
			}
		}
	})
	router.Run(":8080")
}

func UniverseChanged(universe *objects.Universe) {
	universeLatest = universe
	for _, client := range universeClients {
		client.WriteJSON(universe)
	}
}

func FixturesChanged(fixtures []fixtureTypes.PAR) {
	fixturesLatest = fixtures
	JSONfixtures := make([]map[string]interface{}, len(fixtures))
	for i, fixture := range fixtures {
		JSONfixtures[i] = fixture.JSON()
	}

	for _, client := range fixtureClients {
		client.WriteJSON(JSONfixtures)
	}
}
