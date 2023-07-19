package main

import (
	"fmt"
	"log"

	"github.com/parlaynu/go-jetbot-101/internal/vehicle"
)

func main() {
	// create the channel to send the message
	ichan := make(chan *vehicle.WheelSpeeds)

	vchan, err := vehicle.NewDefault(ichan)
	if err != nil {
		log.Fatal(err)
	}

	// send a zero speed message
	ichan <- &vehicle.WheelSpeeds{
		LSpeed: 0.0,
		RSpeed: 0.0,
	}

	// and shutdown...
	close(ichan)

	// and wait...
	for v := range vchan {
		fmt.Printf("left: %0.4f, right: %0.4f\n", v.LSpeed, v.RSpeed)
	}
}
