package gamepad

import (
	"errors"
	"fmt"
	"os"

	"github.com/holoplot/go-evdev"
)

func New(driver string) (<-chan *evdev.InputEvent, func() error, error) {
	device, err := evdev.Open(driver)
	if err != nil {
		return nil, nil, err
	}

	ch := make(chan *evdev.InputEvent)
	go run(ch, device)

	return ch, device.Close, nil
}

func run(ch chan<- *evdev.InputEvent, device *evdev.InputDevice) {
	defer close(ch)

	// set the device to non-block mode so a call to close will
	//   interrupt reading and shutdown nicely
	device.NonBlock()

	for {
		ev, err := device.ReadOne()
		if err != nil {
			if errors.Is(err, os.ErrClosed) == false {
				fmt.Printf("Failed to read event: %v\n", err)
			}
			break
		}
		ch <- ev
	}
}
