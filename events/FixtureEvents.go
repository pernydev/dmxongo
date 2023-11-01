package events

import (
	"dmxongo/api"
	"dmxongo/fixtureTypes"
	"dmxongo/objects"
)

func FixturesChanged(fixtures []fixtureTypes.PAR, universe *objects.Universe) {
	api.FixturesChanged(fixtures)
	UniverseChanged(universe)
}
