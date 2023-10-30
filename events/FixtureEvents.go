package events

import (
	"dmxongo/api"
	"dmxongo/fixtureTypes"
	"dmxongo/objects"
	"fmt"
)

func FixturesChanged(fixtures []fixtureTypes.PAR, universe *objects.Universe) {
	fmt.Println(fixtures[0].Color)
	api.FixturesChanged(fixtures)
	UniverseChanged(universe)
}
