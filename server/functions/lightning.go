package functions

import (
	"dmxongo/events"
	"dmxongo/objects"
	"fmt"
	"math/rand"
	"time"
)

var lightning = NewFunction("lightning", lightningFunction, "basic")

func lightningFunction(_ <-chan struct{}) {
	fmt.Println("lightning!")
	// flash a random fixture
	fixtureID := rand.Intn(len(*Fixtures))
	originalColor := (*Fixtures)[fixtureID].Color
	originalBrightness := (*Fixtures)[fixtureID].Brightness

	// set the fixture to white
	(*Fixtures)[fixtureID].Brightness = 255
	(*Fixtures)[fixtureID].SetColor(objects.Color{Red: 255, Green: 255, Blue: 255, White: 255})
	(*Fixtures)[fixtureID].Update()
	events.FixturesChanged()

	// sleep for 300ms
	<-time.After(200 * time.Millisecond)
	(*Fixtures)[fixtureID].Brightness = originalBrightness
	(*Fixtures)[fixtureID].SetColor(originalColor)
	(*Fixtures)[fixtureID].Update()
	events.FixturesChanged()
}
