package direct_lw

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	"testing"

	"github.com/holoplot/go-evdev"

	"github.com/parlaynu/go-jetbot-101/internal/gamepad"
)

func TestFilter(t *testing.T) {

	// create the filter
	ichan := make(chan *evdev.InputEvent)
	ochan := newFilter(ichan)

	// create the events to test
	var events = map[evdev.EvType]map[evdev.EvCode]bool{
		evdev.EV_ABS: map[evdev.EvCode]bool{
			evdev.ABS_RZ: true,
			evdev.ABS_X:  false,
			evdev.ABS_Y:  true,
			evdev.ABS_Z:  false,
		},
		evdev.EV_KEY: map[evdev.EvCode]bool{
			evdev.BTN_X: true,
			evdev.BTN_Y: false,
			evdev.BTN_Z: false,
		},
		evdev.EV_REL: map[evdev.EvCode]bool{
			evdev.REL_X: false,
			evdev.REL_Y: false,
			evdev.REL_Z: false,
		},
	}

	// goroutine to send all the events into the channel
	go func() {
		defer close(ichan)

		for t, codes := range events {
			for c, _ := range codes {
				ev := evdev.InputEvent{
					Type:  t,
					Code:  c,
					Value: 0,
				}
				ichan <- &ev
			}
		}
	}()

	for ev := range ochan {
		codes, ok := events[ev.Type]
		require.True(t, ok)

		_, ok = codes[ev.Code]
		require.True(t, ok)
	}
}

func TestFilterGamepad(t *testing.T) {
	// get the available devices
	devices, err := gamepad.Devices()
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	// open the first one...
	ichan, closer, err := gamepad.New(devices[0].Driver)
	require.NoError(t, err)
	require.NotNil(t, ichan)
	require.NotNil(t, closer)

	// create the filter
	ochan := newFilter(ichan)

	var valid = map[evdev.EvType]map[evdev.EvCode]bool{
		evdev.EV_ABS: map[evdev.EvCode]bool{
			evdev.ABS_RZ: true,
			evdev.ABS_Y:  true,
		},
		evdev.EV_KEY: map[evdev.EvCode]bool{
			evdev.BTN_X: true,
		},
	}

	// timer for the test
	fmt.Println("closing in 5 seconds...")
	ticker := time.Tick(5 * time.Second)

	// start testing
	running := true
	for i := 0; running; i++ {
		select {
		case ev, ok := <-ochan:
			if ok == false {
				fmt.Println("channel is closed")
				running = false
			} else {
				fmt.Printf("event %03d: %s\n", i, ev)

				codes, ok := valid[ev.Type]
				require.True(t, ok)

				_, ok = codes[ev.Code]
				require.True(t, ok)
			}
		case <-ticker:
			err := closer()
			require.NoError(t, err)
		}
	}
}
