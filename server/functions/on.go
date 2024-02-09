package functions

import (
	"dmxongo/events"
)

var on = NewFunction("on", onFunction, "toggle")

func onFunction(stopCh <-chan struct{}) {
	for fixture := range *Fixtures {
		(*Fixtures)[fixture].Brightness = 255
		(*Fixtures)[fixture].Update()
	}
	events.FixturesChanged()
}
