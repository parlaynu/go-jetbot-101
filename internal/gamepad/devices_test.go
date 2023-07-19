package gamepad_test

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/parlaynu/go-jetbot-101/internal/gamepad"
)

func TestListDevices(t *testing.T) {
	devices, err := gamepad.Devices()
	require.NoError(t, err)
	require.NotEmpty(t, devices)
}
