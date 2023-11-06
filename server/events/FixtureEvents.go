package events

import "fmt"

var FixtureChangeListeners []func()

func FixturesChanged() {
	fmt.Println("Fixtures changed")
	for _, listener := range FixtureChangeListeners {
		fmt.Println("Listener called")
		listener()
	}
	UniverseChanged()
}
