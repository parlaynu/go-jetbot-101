package i2cmotor_test

import (
	"fmt"

	"github.com/stretchr/testify/require"
	"testing"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/host/v3"

	"github.com/parlaynu/go-jetbot-101/internal/i2cmotor"
)

func TestHostState(t *testing.T) {

	// display the host state
	state, err := host.Init()
	require.NoError(t, err)

	fmt.Printf("Loaded Drivers:\n")
	for _, driver := range state.Loaded {
		fmt.Printf("- %s\n", driver)
	}

	fmt.Printf("Skipped Drivers:\n")
	for _, failure := range state.Skipped {
		fmt.Printf("- %s: %s\n", failure.D, failure.Err)
	}

	if len(state.Failed) > 0 {
		fmt.Printf("Failed Drivers:\n")
		for _, failure := range state.Failed {
			fmt.Printf("- %s: %v\n", failure.D, failure.Err)
		}
	}
}

func TestBus(t *testing.T) {
	bus, err := i2cmotor.NewBus("/dev/i2c-1")
	require.NoError(t, err)

	defer bus.Close()

	fmt.Printf("Bus: %s\n", bus.String())
	if p, ok := bus.(i2c.Pins); ok {
		fmt.Printf("SDA: %s\n", p.SDA())
		fmt.Printf("SCL: %s\n", p.SCL())
	}
}
