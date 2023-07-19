package direct_lw

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	"testing"

	"github.com/holoplot/go-evdev"

	"github.com/parlaynu/go-jetbot-101/internal/gamepad"
	"github.com/parlaynu/go-jetbot-101/internal/vehicle"
)

type Event struct {
	Input    evdev.InputEvent
	Expected vehicle.WheelSpeeds
}

func TestDirect(t *testing.T) {
	// create the events to test
	evLfwd := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_Y,
			Value: 64,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: 0.5,
			RSpeed: 0.0,
		},
	}
	evLrev := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_Y,
			Value: 192,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: -0.5,
			RSpeed: 0.0,
		},
	}
	evLstop := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_Y,
			Value: 128,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: 0.0,
			RSpeed: 0.0,
		},
	}
	evRfwd := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_RZ,
			Value: 64,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: 0.0,
			RSpeed: 0.5,
		},
	}
	evRrev := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_RZ,
			Value: 192,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: 0.0,
			RSpeed: -0.5,
		},
	}
	evRstop := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_ABS,
			Code:  evdev.ABS_RZ,
			Value: 128,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: 0.0,
			RSpeed: 0.0,
		},
	}
	evQuit := &Event{
		Input: evdev.InputEvent{
			Type:  evdev.EV_KEY,
			Code:  evdev.BTN_X,
			Value: 0,
		},
		Expected: vehicle.WheelSpeeds{
			LSpeed: 0.0,
			RSpeed: 0.0,
		},
	}

	events := make([]*Event, 0)
	events = append(events, evLfwd)
	events = append(events, evLrev)
	events = append(events, evLstop)
	events = append(events, evRfwd)
	events = append(events, evRrev)
	events = append(events, evRstop)
	events = append(events, evQuit)
	events = append(events, evRstop)

	// push all the events onto a channel
	ichan := make(chan *evdev.InputEvent, len(events))
	for _, ev := range events {
		ichan <- &ev.Input
	}

	// the closer function
	closer := func() error {
		fmt.Println("closing...")
		close(ichan)
		return nil
	}

	// create the controller
	ochan := New(closer, ichan)
	i := 0
	for ; true; i++ {
		oev, ok := <-ochan
		if ok == false {
			break
		}

		iev := events[i]
		require.Equal(t, oev.LSpeed, iev.Expected.LSpeed)
		require.Equal(t, oev.RSpeed, iev.Expected.RSpeed)
	}

	require.Equal(t, len(events), i)
}

func TestDirectGamepad(t *testing.T) {
	// get the available devices
	devices, err := gamepad.Devices()
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	// open the first one...
	ichan, closer, err := gamepad.New(devices[0].Driver)
	require.NoError(t, err)
	require.NotNil(t, ichan)
	require.NotNil(t, closer)

	// create the controller
	ochan := New(closer, ichan)

	// timer for the test
	fmt.Println("closing in 10 seconds...")
	ticker := time.Tick(10 * time.Second)

	// start testing
	running := true
	for i := 0; running; i++ {
		select {
		case ev, ok := <-ochan:
			if ok == false {
				fmt.Println("channel is closed")
				running = false
			} else {
				fmt.Printf("speeds %03d: %0.4f %0.4f\n", i, ev.LSpeed, ev.RSpeed)
			}
		case <-ticker:
			err := closer()
			require.NoError(t, err)
		}
	}
}
