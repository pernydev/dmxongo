package functions

import (
	"dmxongo/events"
)

var blackout = NewFunction("blackout", blackoutFunction, "basic")

func blackoutFunction(stopCh <-chan struct{}) {
	for fixture := range *Fixtures {
		(*Fixtures)[fixture].Brightness = 0
		(*Fixtures)[fixture].Update()
	}
	events.FixturesChanged()
}
