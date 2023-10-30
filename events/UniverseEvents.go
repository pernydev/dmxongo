package events

import (
	"dmxongo/api"
	"dmxongo/objects"
)

func UniverseChanged(universe *objects.Universe) {
	api.UniverseChanged(universe)
}
