package direct_lw

import (
	"math"

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

	fSpeed := 0.0
	wSpeed := 0.0

	l := 0.5

	for ev := range ichan {
		if ev.Type == evdev.EV_KEY && ev.Code == evdev.BTN_X {
			// closing the closer causes the channels to close...
			//   keep processing events until that happens so we don't deadlock
			closer()

		} else if ev.Type == evdev.EV_ABS && ev.Code == evdev.ABS_RZ {
			fSpeed = float64(128-ev.Value) / 128.0

		} else if ev.Type == evdev.EV_ABS && ev.Code == evdev.ABS_Z {
			wSpeed = float64(127-ev.Value) / 128.0
		}

		// fmt.Printf("fwd: %0.4f, rot: %0.4f\n", fSpeed, wSpeed)

		if fSpeed > 0 {
			lSpeed = fSpeed - wSpeed*l/2.0
			rSpeed = fSpeed + wSpeed*l/2.0

			lSpeed = math.Min(1.0, lSpeed)
			rSpeed = math.Min(1.0, rSpeed)

		} else if fSpeed < 0 {
			lSpeed = fSpeed + wSpeed*l/2.0
			rSpeed = fSpeed - wSpeed*l/2.0

			lSpeed = math.Max(-1.0, lSpeed)
			rSpeed = math.Max(-1.0, rSpeed)

		} else {
			lSpeed = -wSpeed
			rSpeed = wSpeed
		}

		// fmt.Printf("left: %0.4f, right: %0.4f\n", lSpeed, rSpeed)

		ws := vehicle.WheelSpeeds{
			LSpeed: lSpeed,
			RSpeed: rSpeed,
		}
		ochan <- &ws
	}
}
