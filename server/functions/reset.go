package functions

import (
	"dmxongo/events"
	"dmxongo/objects"
)

var reset = NewFunction("reset", resetFunction, "basic")

func resetFunction(stopCh <-chan struct{}) {
	for fixture := range *Fixtures {
		(*Fixtures)[fixture].Color = objects.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
			White: 0,
		}
		(*Fixtures)[fixture].Update()
	}
	events.FixturesChanged()
}
