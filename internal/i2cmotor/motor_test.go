package i2cmotor_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/parlaynu/go-jetbot-101/internal/i2cmotor"
)

func TestMotor(t *testing.T) {
	bus, err := i2cmotor.NewBus("/dev/i2c-1")
	require.NoError(t, err)
	defer bus.Close()

	ctl, err := i2cmotor.NewController(bus, 0x60)
	require.NoError(t, err)

	_, err = i2cmotor.NewMotor(ctl, 1, 2, 3)
	require.NoError(t, err)
}
