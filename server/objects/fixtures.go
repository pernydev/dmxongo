package objects

var fixtures = map[string]func(fixture Fixture) []int{
	"PAR": PAR,
}

func PAR(fixture Fixture) []int {
	fixture.ChannelValues[0] = fixture.Color.Red
	fixture.ChannelValues[1] = fixture.Color.Green
	fixture.ChannelValues[2] = fixture.Color.Blue
	fixture.ChannelValues[3] = fixture.Brightness
	fixture.ChannelValues[4] = fixture.Strobe
	return fixture.ChannelValues
}
