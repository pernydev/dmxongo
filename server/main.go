package main

import (
	"dmxongo/api"
	"dmxongo/events"
	"dmxongo/fixtureTypes"
	"dmxongo/objects"
	"time"
)

var universe objects.Universe
var fixtures []fixtureTypes.PAR

func main() {
	universe = objects.NewUniverse()
	fixtures = []fixtureTypes.PAR{
		fixtureTypes.MakePAR(0, objects.Color{Red: 255, Green: 0, Blue: 0, White: 0}, 255, 0, universe),
		fixtureTypes.MakePAR(5, objects.Color{Red: 255, Green: 0, Blue: 0, White: 0}, 255, 0, universe),
	}

	// chips.Init()

	for _, fixture := range fixtures {
		fixture.Update()
	}

	go api.HTTPAPI(&universe, fixtures)

	frames := []objects.Color{
		{Red: 255, Green: 0, Blue: 0, White: 0},
		{Red: 255, Green: 255, Blue: 0, White: 0},
		{Red: 0, Green: 255, Blue: 0, White: 0},
		{Red: 255, Green: 255, Blue: 255, White: 0},
		{Red: 255, Green: 0, Blue: 255, White: 0},
	}

	for {
		for _, frame := range frames {
			for fixture := range fixtures {
				fixtures[fixture].SetColor(frame)
			}
			events.FixturesChanged(fixtures, &universe)
			time.Sleep(50 * time.Millisecond)
		}
	}
}
