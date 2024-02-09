package main

import (
	"dmxongo/api"
	"dmxongo/functions"
	"dmxongo/objects"
	"fmt"
)

var universe objects.Universe
var fixtures []objects.Fixture = []objects.Fixture{
	objects.MakeFixture("MOV", 121, 0, &universe, 15),
	objects.MakeFixture("MOV", 91, 0, &universe, 15),
	objects.MakeFixture("NMOV", 30, 0, &universe, 15),
	objects.MakeFixture("NMOV", 154, 0, &universe, 15),
	objects.MakeFixture("DIM", 1, 0, &universe, 1),
}

func main() {
	universe = objects.NewUniverse()

	// chips.Init()

	for k, fixture := range fixtures {
		if fixture.Type == "NMOV" {
			fixture.Pan = 10
			fixture.Tilt = 155
		} else if fixture.Type == "MOV" {
			fixture.Pan = 0
			fixture.Tilt = 127
		}
		fixtures[k] = fixture
		fixture.Update()
	}

	functions.Fixtures = &fixtures

	go api.HTTPAPI(&universe, &fixtures)

	fmt.Scanln()
}
