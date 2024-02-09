package objects

var fixtures = map[string]func(fixture Fixture) []int{
	"PAR":  PAR,
	"DIM":  DIM,
	"MOV":  MOV,
	"NMOV": NMOV,
}

func PAR(fixture Fixture) []int {
	fixture.ChannelValues[0] = fixture.Color.Red
	fixture.ChannelValues[1] = fixture.Color.Green
	fixture.ChannelValues[2] = fixture.Color.Blue
	fixture.ChannelValues[3] = fixture.Brightness
	fixture.ChannelValues[4] = fixture.Strobe
	return fixture.ChannelValues
}

func DIM(fixture Fixture) []int {
	fixture.ChannelValues[0] = fixture.Brightness
	return fixture.ChannelValues
}

func MOV(fixture Fixture) []int {
	fixture.ChannelValues[0] = fixture.Focus
	switch fixture.Color { // It's a gobo wheel so it can't take RGB values
	case RGB(255, 0, 0):
		fixture.ChannelValues[1] = 131
	case RGB(0, 255, 0):
		fixture.ChannelValues[1] = 30
	case RGB(0, 0, 255):
		fixture.ChannelValues[1] = 165
	case RGB(255, 255, 0):
		fixture.ChannelValues[1] = 96
	case RGB(255, 0, 255):
		fixture.ChannelValues[1] = 62
	case RGB(0, 0, 0):
		fixture.ChannelValues[1] = 0
	}

	fixture.ChannelValues[2] = fixture.Gobo // shape gobo

	fixture.ChannelValues[4] = fixture.Pan
	fixture.ChannelValues[6] = fixture.Tilt
	fixture.ChannelValues[8] = fixture.Brightness
	fixture.ChannelValues[9] = 255 // rotation
	fixture.ChannelValues[10] = 0

	return fixture.ChannelValues
}

func NMOV(fixture Fixture) []int {
	switch fixture.Color { // It's a gobo wheel so it can't take RGB values
	case RGB(255, 0, 0):
		fixture.ChannelValues[1] = 59
	case RGB(0, 255, 0):
		fixture.ChannelValues[1] = 21
	case RGB(0, 0, 255):
		fixture.ChannelValues[1] = 50
	case RGB(255, 255, 0):
		fixture.ChannelValues[1] = 96
	case RGB(255, 0, 255):
		fixture.ChannelValues[1] = 39
	case RGB(0, 0, 0):
		fixture.ChannelValues[1] = 0
	}

	fixture.ChannelValues[2] = fixture.Gobo // shape gobo

	fixture.ChannelValues[4] = fixture.Pan
	fixture.ChannelValues[6] = fixture.Tilt
	fixture.ChannelValues[8] = fixture.Brightness
	fixture.ChannelValues[9] = 255 // rotation
	fixture.ChannelValues[10] = 255

	return fixture.ChannelValues
}
