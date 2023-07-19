package vehicle

import (
	"fmt"

	"periph.io/x/conn/v3/i2c"

	"github.com/parlaynu/go-jetbot-101/internal/i2cmotor"
)

type WheelSpeeds struct {
	LSpeed float64
	RSpeed float64
}

type Vehicle struct {
	ichan      <-chan *WheelSpeeds
	ochan      chan<- *WheelSpeeds
	bus        i2c.BusCloser
	controller *i2cmotor.Controller
	left       *i2cmotor.Motor
	right      *i2cmotor.Motor
	maxSpeed   float64
}

func NewDefault(ichan <-chan *WheelSpeeds) (<-chan *WheelSpeeds, error) {
	return New(ichan, "/dev/i2c-1", 1, false)
}

func New(ichan <-chan *WheelSpeeds, i2cdriver string, maxspeed int, swapmotors bool) (<-chan *WheelSpeeds, error) {

	// open the i2c bus
	bus, err := i2cmotor.NewBus(i2cdriver)
	if err != nil {
		return nil, err
	}

	// the controller
	ctl, err := i2cmotor.NewController(bus, 0x60)
	if err != nil {
		bus.Close()
		return nil, err
	}

	// define the motors
	in1L, in2L, pwmL := 11, 12, 13
	in1R, in2R, pwmR := 10, 9, 8

	if swapmotors == true {
		in1L, in1R = in1R, in1L
		in2L, in2R = in2R, in2L
		pwmL, pwmR = pwmR, pwmL
	}

	// build the motors
	left, err := i2cmotor.NewMotor(ctl, in1L, in2L, pwmL)
	if err != nil {
		bus.Close()
		return nil, err
	}

	right, err := i2cmotor.NewMotor(ctl, in1R, in2R, pwmR)
	if err != nil {
		bus.Close()
		return nil, err
	}

	ochan := make(chan *WheelSpeeds)

	veh := Vehicle{
		ichan:      ichan,
		ochan:      ochan,
		bus:        bus,
		controller: ctl,
		left:       left,
		right:      right,
	}
	veh.setMaxSpeed(maxspeed)

	// set the vehicle running
	go veh.run()

	return ochan, nil
}

func (j *Vehicle) run() {
	defer close(j.ochan)

	for speeds := range j.ichan {
		j.setSpeed(speeds.LSpeed, speeds.RSpeed)
		j.ochan <- speeds
	}

	// shutdown cleanly
	j.stop()
	j.close()
}

func (j *Vehicle) setMaxSpeed(mspeed int) {
	// clamp to the range that works
	if mspeed < 512 {
		mspeed = 512
	}
	if mspeed > 4095 {
		mspeed = 4095
	}

	j.maxSpeed = float64(mspeed)
}

func (j *Vehicle) mapSpeed(ivalue float64) int {
	// motor needs at least 400 to run... this leaves a large dead patch where
	//   the motor is straining especially for low max speeds. rescale the values
	//   to skip over this dead area and remap the values

	ispeed := ivalue * j.maxSpeed

	// keep a small no-speed zone
	if ispeed < 100.0 && ispeed > -100.0 {
		return 0
	}

	fmt.Printf("ispeed: %0.2f\n", ispeed)

	// map the value
	sign := 1.0
	if ispeed < 0 {
		sign = -1.0
		ispeed *= -1.0
	}

	ospeed := sign * (384 + (ispeed-100.0)/(j.maxSpeed-100.0)*(j.maxSpeed-384.0))
	fmt.Printf("ospeed: %0.2f\n", ospeed)

	return int(ospeed)
}

func (j *Vehicle) setSpeed(lspeed, rspeed float64) {
	// translate speed into register speed
	lmotor := j.mapSpeed(lspeed)
	rmotor := j.mapSpeed(rspeed)

	switch {
	case lmotor > 0:
		j.left.Forward(lmotor)
	case lmotor < 0:
		j.left.Backward(-lmotor)
	default:
		j.left.Stop()
	}

	switch {
	case rmotor > 0:
		j.right.Forward(rmotor)
	case rmotor < 0:
		j.right.Backward(-rmotor)
	default:
		j.right.Stop()
	}
}

func (j *Vehicle) stop() {
	j.left.Stop()
	j.right.Stop()
}

func (j *Vehicle) close() {
	j.bus.Close()
}
