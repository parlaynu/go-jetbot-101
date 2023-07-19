package direct_lw

import (
	"github.com/holoplot/go-evdev"
)

func newFilter(ichan <-chan *evdev.InputEvent) <-chan *evdev.InputEvent {
	ochan := make(chan *evdev.InputEvent)
	go filter(ichan, ochan)
	return ochan
}

func filter(ichan <-chan *evdev.InputEvent, ochan chan<- *evdev.InputEvent) {
	defer close(ochan)

	var filter = map[evdev.EvType]map[evdev.EvCode]bool{
		evdev.EV_ABS: map[evdev.EvCode]bool{
			evdev.ABS_RZ: true,
			evdev.ABS_Z:  true,
		},
		evdev.EV_KEY: map[evdev.EvCode]bool{
			evdev.BTN_X: true,
		},
	}

	for ev := range ichan {
		if codes, ok := filter[ev.Type]; ok {
			if _, ok := codes[ev.Code]; ok {
				ochan <- ev
			}
		}
	}
}
