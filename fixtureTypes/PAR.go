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
	UniverseCTX     objects.Universe
}

func (p *PAR) Update() {
	p.ChannelValues[0] = p.Brightness
	p.ChannelValues[1] = p.Color.Red
	p.ChannelValues[2] = p.Color.Green
	p.ChannelValues[3] = p.Color.Blue
	p.ChannelValues[4] = p.Strobe

	p.UniverseCTX.Update(p.StartingChannel, p.ChannelValues)
}

func (p *PAR) JSON() map[string]interface{} {
	// return a JSON representation of the fixture
	response := map[string]interface{}{
		"startingChannel": p.StartingChannel,
		"color": map[string]int{
			"red":   p.Color.Red,
			"green": p.Color.Green,
			"blue":  p.Color.Blue,
			"white": p.Color.White,
		},
		"brightness": p.Brightness,
		"strobe":     p.Strobe,
		"channels":   p.ChannelValues,
	}

	return response
}

func (p *PAR) SetColor(color objects.Color) {
	p.Color = color
	p.Update()
}

func MakePAR(startingChannel int, color objects.Color, brightness int, strobe int, universeCTX objects.Universe) PAR {
	par := PAR{
		StartingChannel: startingChannel,
		Color:           color,
		Brightness:      brightness,
		Strobe:          strobe,
		ChannelValues:   make([]int, 5),
		UniverseCTX:     universeCTX,
	}

	return par
}
