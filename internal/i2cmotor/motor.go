package i2cmotor

type Motor struct {
	C   *Controller
	In1 int
	In2 int
	Pwm int
}

func NewMotor(c *Controller, in1, in2, pwm int) (*Motor, error) {

	m := &Motor{
		C:   c,
		In1: in1,
		In2: in2,
		Pwm: pwm,
	}

	return m, nil
}

func (m *Motor) Forward(speed int) {
	m.C.SetPWM(m.Pwm, 0, speed)
	m.C.SetPWM(m.In1, 4096, 0) // bit 12 set: full on
	m.C.SetPWM(m.In2, 0, 4096) // bit 12 set: full off
}

func (m *Motor) Backward(speed int) {
	m.C.SetPWM(m.Pwm, 0, speed)
	m.C.SetPWM(m.In1, 0, 4096) // bit 12 set: full off
	m.C.SetPWM(m.In2, 4096, 0) // bit 12 set: full on
}

func (m *Motor) Stop() {
	m.C.SetPWM(m.Pwm, 4096, 0) // bit 12 set: full on - per the spec
	m.C.SetPWM(m.In1, 0, 4096) // bit 12 set: full off
	m.C.SetPWM(m.In2, 0, 4096) // bit 12 set: full off
}
