package i2cmotor

import (
	"fmt"
	"math"
	"time"

	"periph.io/x/conn/v3/i2c"
)

type Controller struct {
	i2c.Dev
}

// registers
const (
	__MODE1         byte = 0x00
	__MODE2              = 0x01
	__SUBADR1            = 0x02
	__SUBADR2            = 0x03
	__SUBADR3            = 0x04
	__LED0_ON_L          = 0x06
	__LED0_ON_H          = 0x07
	__LED0_OFF_L         = 0x08
	__LED0_OFF_H         = 0x09
	__ALL_LED_ON_L       = 0xFA
	__ALL_LED_ON_H       = 0xFB
	__ALL_LED_OFF_L      = 0xFC
	__ALL_LED_OFF_H      = 0xFD
	__PRESCALE           = 0xFE
)

// mode 1
const (
	__ALLCALL byte = 0x01
	__SLEEP        = 0x10
	__RESTART      = 0x80
)

// mode 2
const (
	__INVRT  byte = 0x10
	__OUTDRV      = 0x04
	__OCH         = 0x80
)

func NewController(bus i2c.Bus, address uint16) (*Controller, error) {
	ctl := &Controller{
		Dev: i2c.Dev{
			Bus:  bus,
			Addr: address,
		},
	}

	err := ctl.Init()
	if err != nil {
		return nil, err
	}

	err = ctl.SetFrequency(1600)
	if err != nil {
		return nil, err
	}

	return ctl, nil
}

func clamp(value int) int {
	switch {
	case value < 0:
		return 0
	case value > 4096:
		return 4096
	default:
		return value
	}
}

func (c *Controller) SetPWM(pin int, on, off int) error {
	on = clamp(on)
	err := c.Write(byte((__LED0_ON_L+4*pin)&0xff), byte(on&0xff))
	if err != nil {
		return err
	}
	err = c.Write(byte((__LED0_ON_H+4*pin)&0xff), byte((on>>8)&0xff))
	if err != nil {
		return err
	}

	off = clamp(off)
	err = c.Write(byte((__LED0_OFF_L+4*pin)&0xff), byte(off&0xff))
	if err != nil {
		return err
	}
	err = c.Write(byte((__LED0_OFF_H+4*pin)&0xff), byte((off>>8)&0xff))
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) Write(values ...byte) error {
	err := c.Tx(
		values,
		nil,
	)
	return err
}

func (c *Controller) Read8(values ...byte) (byte, error) {
	result := make([]byte, 1)
	err := c.Tx(
		values,
		result,
	)
	return result[0], err
}

func (c *Controller) Init() error {

	// set all outputs to off
	err := c.Write(__ALL_LED_ON_L, 0)
	if err != nil {
		return fmt.Errorf("Controller Init: __ALL_LED_ON_L: %v", err)
	}

	err = c.Write(__ALL_LED_ON_H, 0)
	if err != nil {
		return fmt.Errorf("Controller Init: __ALL_LED_ON_H: %v", err)
	}

	err = c.Write(__ALL_LED_OFF_L, 0)
	if err != nil {
		return fmt.Errorf("Controller Init: __ALL_LED_OFF_L: %v", err)
	}

	err = c.Write(__ALL_LED_OFF_H, 0x80) // 12th bit set: all off
	if err != nil {
		return fmt.Errorf("Controller Init: __ALL_LED_OFF_H: %v", err)
	}

	// configure mode1 register
	err = c.Write(__MODE1, __ALLCALL)
	if err != nil {
		return fmt.Errorf("Controller Init: __MODE1: __ALLCALL: %v", err)
	}
	time.Sleep(5 * time.Millisecond)

	// configure mode2 register
	err = c.Write(__MODE2, __OUTDRV|__OCH)
	if err != nil {
		return fmt.Errorf("Controller Init: __MODE2: __OUTDRV: %v", err)
	}
	time.Sleep(5 * time.Millisecond)

	return nil
}

func (c *Controller) SetFrequency(freq int) error {

	// calculate the prescale value to set
	prescaleval := 25000000.0 // 25MHz
	prescaleval /= 4096.0     // 12-bit
	prescaleval /= float64(freq)
	prescaleval -= 1.0

	prescale := uint8(math.Floor(prescaleval + 0.5))

	// turn off the oscillator
	oldmode, err := c.Read8(__MODE1)
	if err != nil {
		return fmt.Errorf("Controller SetFrequency: Read Mode1: %v", err)
	}
	newmode := (oldmode | __SLEEP) & (^__RESTART & 0xff) // don't write 1 to the RESTART bit...

	err = c.Write(__MODE1, newmode)
	if err != nil {
		return fmt.Errorf("Controller SetFrequency: Disable Oscillator: 0x%x: %v", newmode, err)
	}
	time.Sleep(5 * time.Millisecond)

	// set the frequency
	err = c.Write(__PRESCALE, prescale)
	if err != nil {
		return fmt.Errorf("Controller SetFrequency: __PRESCALE: %v", err)
	}
	time.Sleep(5 * time.Millisecond)

	// turn off sleep
	err = c.Write(__MODE1, oldmode)
	if err != nil {
		return fmt.Errorf("Controller SetFrequency: Enable Oscillator: 0x%x: %v", oldmode, err)
	}
	time.Sleep(5 * time.Millisecond)

	// restart the oscillator
	err = c.Write(__MODE1, oldmode|__RESTART)
	if err != nil {
		return fmt.Errorf("Controller SetFrequency: Restart Oscillator: 0x%x: %v", oldmode, err)
	}
	time.Sleep(5 * time.Millisecond)

	return nil
}
