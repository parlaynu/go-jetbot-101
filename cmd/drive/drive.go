package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/parlaynu/go-jetbot-101/internal/controllers/direct_lr"
	"github.com/parlaynu/go-jetbot-101/internal/controllers/direct_lw"
	"github.com/parlaynu/go-jetbot-101/internal/gamepad"
	"github.com/parlaynu/go-jetbot-101/internal/vehicle"
)

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-v] [-x] [-s top-speed] [-c controller]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	verbose := flag.Bool("v", false, "verbose output")
	swapmotors := flag.Bool("x", false, "motors are wired the other way around")
	speedscale := flag.Int("s", 2, "set the top speed scale (1 - 8)")
	controller := flag.Int("c", 0, "select the controller to use (0 - lr, 1 - lw)")
	flag.Parse()

	if *speedscale < 1 || *speedscale > 8 {
		log.Fatalf("Error: topspeed parameter out of range")
	}
	maxspeed := *speedscale * 512

	// create the gamepad
	devices, err := gamepad.Devices()
	if err != nil {
		log.Fatal(err)
	}
	gchan, closer, err := gamepad.New(devices[0].Driver)
	if err != nil {
		log.Fatal(err)
	}

	// create the controller
	var cchan <-chan *vehicle.WheelSpeeds
	if *controller == 0 {
		cchan = direct_lr.New(closer, gchan)
	} else {
		cchan = direct_lw.New(closer, gchan)
	}

	// create the vehicle
	vchan, err := vehicle.New(cchan, "/dev/i2c-1", maxspeed, *swapmotors)
	if err != nil {
		log.Fatal(err)
	}

	// print something out...
	fmt.Println("jetbot is ready...")
	if *controller == 0 {
		fmt.Println("  right stick = right wheel")
		fmt.Println("   left stick = left wheel")
		fmt.Println("     x button = stop")
	} else {
		fmt.Println("  right stick fb = forward speed")
		fmt.Println("  right stick lr = angular speed")
		fmt.Println("        x button = stop")
	}

	// and wait...
	for v := range vchan {
		if *verbose {
			fmt.Printf("left: %0.4f, right: %0.4f\n", v.LSpeed, v.RSpeed)
		}
	}
	fmt.Println("done")
}
