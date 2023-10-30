package fixtureTypes

import (
	"dmxongo/objects"
)

type PAR struct {
	StartingChannel int           // DMX Channel the fixture is configured for
	Color           objects.Color // Color of the fixture
	Brightness      int           // Brightness of the fixture
	Strobe          int           // Strobe value of the fixture
	ChannelValues   []int         // Array of channel values
}

func (p *PAR) Update() {
	p.ChannelValues[0] = p.Brightness
	p.ChannelValues[1] = p.Color.Red
	p.ChannelValues[2] = p.Color.Green
	p.ChannelValues[3] = p.Color.Blue
	p.ChannelValues[4] = p.Strobe
}
