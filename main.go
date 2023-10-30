package main

import (
	"dmxongo/api"
	"dmxongo/events"
	"dmxongo/fixtureTypes"
	"dmxongo/objects"
	"time"
)

func main() {
	universe := objects.NewUniverse()
	fixtures := []fixtureTypes.PAR{
		fixtureTypes.MakePAR(1, objects.Color{Red: 255, Green: 0, Blue: 0, White: 0}, 255, 0, universe),
	}

	for _, fixture := range fixtures {
		fixture.Update()
	}

	go api.API(&universe, fixtures)

	// fade blue channel in and out
	for {
		for i := 0; i < 255; i++ {
			fixtures[0].Color.Blue = i
			fixtures[0].Update()
			time.Sleep(10 * time.Millisecond)
			events.FixturesChanged(fixtures, &universe)
		}

		for i := 255; i > 0; i-- {
			fixtures[0].Color.Blue = i
			fixtures[0].Update()
			time.Sleep(10 * time.Millisecond)
			events.FixturesChanged(fixtures, &universe)
		}
	}

}
