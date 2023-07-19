package vehicle_test

import (
	"fmt"

	"github.com/stretchr/testify/require"
	"testing"

	"github.com/parlaynu/go-jetbot-101/internal/vehicle"
)

func TestVehicle(t *testing.T) {
	ichan := make(chan *vehicle.WheelSpeeds)

	vchan, err := vehicle.NewDefault(ichan)
	require.NoError(t, err)

	close(ichan)

	// wait for vehicle to finish
	for v := range vchan {
		fmt.Printf("left: %0.4f, right: %0.4f\n", v.LSpeed, v.RSpeed)
	}
}
