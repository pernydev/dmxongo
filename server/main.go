package main

import (
	"dmxongo/api"
	"dmxongo/functions"
	"dmxongo/objects"
	"fmt"
)

var universe objects.Universe
var fixtures []objects.Fixture = []objects.Fixture{
	objects.MakeFixture("PAR", 0, 255, &universe),
	objects.MakeFixture("PAR", 5, 255, &universe),
}

func main() {
	universe = objects.NewUniverse()

	// chips.Init()

	for _, fixture := range fixtures {
		fixture.Update()
	}

	functions.Fixtures = &fixtures

	go api.HTTPAPI(&universe, &fixtures)

	fmt.Scanln()
}
