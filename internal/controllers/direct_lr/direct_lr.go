package direct_lr

import (
	"github.com/holoplot/go-evdev"

	"github.com/parlaynu/go-jetbot-101/internal/vehicle"
)

func New(closer func() error, ichan <-chan *evdev.InputEvent) <-chan *vehicle.WheelSpeeds {
	// put the filter into the process chain
	ichan = newFilter(ichan)

	// create the controller
	ochan := make(chan *vehicle.WheelSpeeds)
	go run(closer, ichan, ochan)
	return ochan
}

func run(closer func() error, ichan <-chan *evdev.InputEvent, ochan chan<- *vehicle.WheelSpeeds) {
	defer close(ochan)

	lSpeed := 0.0
	rSpeed := 0.0

	for ev := range ichan {
		if ev.Type == evdev.EV_KEY && ev.Code == evdev.BTN_X {
			// closing the closer causes the channels to close...
			//   keep processing events until that happens so we don't deadlock
			closer()

		} else if ev.Type == evdev.EV_ABS && ev.Code == evdev.ABS_Y {
			lSpeed = float64(128-ev.Value) / 128.0

		} else if ev.Type == evdev.EV_ABS && ev.Code == evdev.ABS_RZ {
			rSpeed = float64(128-ev.Value) / 128.0
		}

		ws := vehicle.WheelSpeeds{
			LSpeed: lSpeed,
			RSpeed: rSpeed,
		}
		ochan <- &ws
	}
}
