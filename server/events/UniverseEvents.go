package events

var UniverseChangeListeners []func()

func UniverseChanged() {
	for _, listener := range UniverseChangeListeners {
		listener()
	}
}
