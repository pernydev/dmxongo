package objects

type Fixture struct {
	Type            string
	StartingChannel int   // DMX Channel the fixture is configured for
	Color           Color // Color of the fixture
	Brightness      int   // Brightness of the fixture
	Strobe          int   // Strobe value of the fixture
	ChannelValues   []int // Array of channel values
	UniverseCTX     *Universe
	Pan             int
	Tilt            int
}

func (p *Fixture) Update() {
	p.ChannelValues = fixtures[p.Type](*p)
	p.UniverseCTX.Update(p.StartingChannel, p.ChannelValues)
}

func (p *Fixture) JSON() map[string]interface{} {
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
		"position": map[string]int{
			"pan":  p.Pan,
			"tilt": p.Tilt,
		},
	}

	return response
}

func (p *Fixture) SetColor(color Color) {
	p.Color = color
	p.Update()
}

func MakeFixture(typeName string, startingChannel int, brightness int, universeCTX *Universe) Fixture {
	par := Fixture{
		Type:            typeName,
		StartingChannel: startingChannel,
		Color:           Color{Red: 0, Green: 0, Blue: 0, White: 0},
		Brightness:      brightness,
		Strobe:          0,
		ChannelValues:   make([]int, 5),
		UniverseCTX:     universeCTX,
	}

	return par
}
