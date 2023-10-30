package fixtureTypes

import (
	"dmxongo/objects"
)

type SchoolPAR struct {
	StartingChannel int           // DMX Channel the fixture is configured for
	Color           objects.Color // Color of the fixture
	Brightness      int           // Brightness of the fixture
	Strobe          int           // Strobe value of the fixture
	ChannelValues   []int         // Array of channel values
}

func (p *SchoolPAR) Update() {
	// TODO: Find the light specs and make this function work
}
