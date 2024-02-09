package functions

import (
	"dmxongo/objects"
	"fmt"
	"sync"
)

var Fixtures *[]objects.Fixture

type Function struct {
	Name     string
	Type     string
	function func(stopCh <-chan struct{})
	stopCh   chan struct{}
	wg       sync.WaitGroup
	running  bool
	mu       sync.Mutex
}

func NewFunction(name string, function func(stopCh <-chan struct{}), functionType string) *Function {
	return &Function{
		Name:     name,
		function: function,
		Type:     functionType,
	}
}

func (f *Function) Start() {
	fmt.Println("of type", f.Type)
	if f.Type == "basic" {
		go f.function(nil)
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if f.running {
		return
	}

	f.stopCh = make(chan struct{}) // Create a new channel
	f.running = true
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.function(f.stopCh)
	}()
}

func (f *Function) Stop() {
	if f.Type == "basic" {
		return
	}
	f.mu.Lock()
	defer f.mu.Unlock()

	if !f.running {
		return
	}

	close(f.stopCh) // Close the current channel
	f.wg.Wait()
	f.running = false
}

var Functions = map[string]*Function{
	"testanimation": testanimation,
	"fire":          fire,
	"reset":         reset,
	"lightning":     lightning,
	"on":            on,
	"blackout":      blackout,
}
