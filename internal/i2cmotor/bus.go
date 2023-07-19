package i2cmotor

import (
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func NewBus(driver string) (i2c.BusCloser, error) {

	// load all the devices/drivers
	_, err := host.Init()
	if err != nil {
		return nil, err
	}

	// open the bus
	b, err := i2creg.Open(driver)
	if err != nil {
		return nil, err
	}

	return b, nil
}
