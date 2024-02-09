package functions

import (
	"dmxongo/events"
	"dmxongo/objects"
	"time"
)

var testanimation = NewFunction("testanimation", testanimationFunction, "toggle")

func testanimationFunction(stopCh <-chan struct{}) {
	frames := []objects.Color{
		objects.RGB(255, 0, 0),
		objects.RGB(0, 255, 0),
		objects.RGB(0, 0, 255),
	}

	for {
		for _, frame := range frames {
			select {
			case <-stopCh:
				// set all the fictures to black
				for _, fixture := range *Fixtures {
					fixture.SetColor(objects.Color{})
					fixture.Brightness = 255
				}
				events.FixturesChanged()

				return
			default:
				for fixture := range *Fixtures {
					(*Fixtures)[fixture].SetColor(frame)
				}
				events.FixturesChanged()
				time.Sleep(4 * time.Second)
			}
		}
	}
}
