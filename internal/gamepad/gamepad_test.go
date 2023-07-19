package gamepad_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/require"
	"testing"

	"github.com/parlaynu/go-jetbot-101/internal/gamepad"
)

func TestOpenAndClose(t *testing.T) {
	// get the available devices
	devices, err := gamepad.Devices()
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	// open and close the first device
	ch, closer, err := gamepad.New(devices[0].Driver)
	require.NoError(t, err)
	require.NotNil(t, ch)
	require.NotNil(t, closer)

	err = closer()
	require.NoError(t, err)

	for i := 0; true; i++ {
		ev, ok := <-ch
		if ok == false {
			fmt.Println("channel is closed")
			break
		}
		fmt.Printf("event %03d: %s\n", i, ev)
	}
}

func TestTimeout(t *testing.T) {
	// get the available devices
	devices, err := gamepad.Devices()
	require.NoError(t, err)
	require.NotEmpty(t, devices)

	// open the first one...
	ch, closer, err := gamepad.New(devices[0].Driver)
	require.NoError(t, err)
	require.NotNil(t, ch)
	require.NotNil(t, closer)

	fmt.Println("closing in 5 seconds...")
	ticker := time.Tick(5 * time.Second)

	running := true
	for i := 0; running; i++ {
		select {
		case ev, ok := <-ch:
			if ok == false {
				fmt.Println("channel is closed")
				running = false
			} else {
				fmt.Printf("event %03d: %s\n", i, ev)
			}
		case <-ticker:
			err := closer()
			require.NoError(t, err)
		}
	}
}
