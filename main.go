package main

import (
	"dmxongo/api"
	"dmxongo/chips"
	"dmxongo/events"
	"dmxongo/fixtureTypes"
	"dmxongo/objects"
	"dmxongo/utils"
	"math/rand"
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

	chips.Init()

	for _, fixture := range fixtures {
		fixture.Update()
	}

	go api.HTTPAPI(&universe, fixtures)

	fireFrames := []objects.Color{
		{Red: 252, Green: 186, Blue: 3, White: 0},
		{Red: 255, Green: 162, Blue: 48, White: 0},
		{Red: 235, Green: 133, Blue: 23, White: 0},
		{Red: 255, Green: 162, Blue: 48, White: 0},
		{Red: 252, Green: 236, Blue: 17, White: 0},
		{Red: 252, Green: 186, Blue: 3, White: 0},
	}

	frames := utils.FadeColorArray(fireFrames, 100*time.Millisecond, 10)

	for fixture := range fixtures {
		go runFrames(fixture, frames)
	}

	for {
	}

}

func runFrames(fixture int, frames []objects.Color) {
	time.Sleep(time.Duration(rand.Intn(155)+25) * time.Millisecond)
	for {
		for _, frame := range frames {
			start := time.Now()
			fixtures[fixture].SetColor(frame)
			events.FixturesChanged(fixtures, &universe)
			elapsed := time.Since(start)
			api.SendUpdateSpeed(elapsed)
			time.Sleep(time.Duration(rand.Intn(155)+25) * time.Millisecond)
		}
	}
}
