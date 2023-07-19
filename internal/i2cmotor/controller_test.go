package i2cmotor_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/parlaynu/go-jetbot-101/internal/i2cmotor"
)

func TestNoController(t *testing.T) {
	bus, err := i2cmotor.NewBus("/dev/i2c-1")
	require.NoError(t, err)
	defer bus.Close()

	_, err = i2cmotor.NewController(bus, 0x16)
	require.Error(t, err)
}

func TestController(t *testing.T) {
	bus, err := i2cmotor.NewBus("/dev/i2c-1")
	require.NoError(t, err)
	defer bus.Close()

	_, err = i2cmotor.NewController(bus, 0x60)
	require.NoError(t, err)
}
