package objects

type Color struct {
	Red   int
	Green int
	Blue  int
	White int
}

func RGB(r int, g int, b int) Color {
	return Color{
		Red:   r,
		Green: g,
		Blue:  b,
		White: 0,
	}
}

type Universe struct {
	ChannelValues []int
}

func (u *Universe) Update(startingChannel int, channelValues []int) {
	// Update the Universe with the new values
	for i := 0; i < len(channelValues); i++ {
		u.ChannelValues[startingChannel+i] = channelValues[i]
	}
}

func NewUniverse() Universe {
	// return a new Universe
	universe := Universe{
		ChannelValues: make([]int, 512),
	}

	return universe
}
