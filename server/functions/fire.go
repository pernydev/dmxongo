package functions

import (
	"dmxongo/events"
	"dmxongo/objects"
	"dmxongo/utils"
	"time"

	"math/rand"
)

var fire = NewFunction("fire", fireFunction, "toggle")

func fireFunction(stopCh <-chan struct{}) {
	// objects.RGB(255, 100, 0),
	frames := []objects.Color{
		objects.RGB(255, 100, 10),
		objects.RGB(255, 100, 0),
		objects.RGB(255, 100, 15),
		objects.RGB(255, 100, 0),
	}

	frames = utils.Fade(frames, 10)

	for {
		for _, frame := range frames {
			select {
			case <-stopCh:
				return
			default:
				for fixture := range *Fixtures {
					(*Fixtures)[fixture].SetColor(frame)
				}
				events.FixturesChanged()
				time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
			}
		}
	}
}
