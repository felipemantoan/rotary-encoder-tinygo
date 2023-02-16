package main

import (
	"machine"
	"time"
)

type Clock struct {
	Now        int
	Previous   int
	Transition int
}

type Encoder struct {
	Name      string
	Degrees   int32
	ClockPin  machine.Pin
	DataPin   machine.Pin
	SwitchPin machine.Pin
	Clock     *Clock
}

func (encoder *Encoder) Pressed() {
	if !encoder.SwitchPin.Get() {
		time.Sleep(time.Millisecond * 250)
	}
}

func (encoder *Encoder) CheckEncoder() {

	dta := 0
	if encoder.DataPin.Get() {
		dta = 1
	}

	clk := 0
	if encoder.ClockPin.Get() {
		clk = 2
	}

	encoder.Clock.Now = dta<<1 | clk
	encoder.Clock.Transition = encoder.Clock.Now<<2 | clk

	if encoder.Clock.Now == encoder.Clock.Previous {
		return
	}

	if encoder.Clock.Transition == 10 {
		encoder.Degrees += 20
		println("Turing clockwise", encoder.Name, encoder.Degrees)
	}

	if encoder.Clock.Transition == 8 {
		encoder.Degrees -= 20
		println("Turing anti-clockwise", encoder.Name, encoder.Degrees)
	}

	if encoder.Degrees >= 360 || encoder.Degrees == -360 {
		encoder.Degrees = 0
	}

	encoder.Clock.Previous = encoder.Clock.Now

	time.Sleep(time.Millisecond * 1)
}

func main() {

	machine.GP17.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.GP16.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.GP18.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	machine.GP19.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.GP20.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	machine.GP21.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	encoders := []*Encoder{
		{
			Name:      "Encoder 1",
			Degrees:   0,
			ClockPin:  machine.GP17,
			DataPin:   machine.GP16,
			SwitchPin: machine.GP18,
			Clock: &Clock{
				Now:        0,
				Previous:   0,
				Transition: 0,
			},
		},
		{
			Name:      "Encoder 2",
			Degrees:   20,
			ClockPin:  machine.GP19,
			DataPin:   machine.GP20,
			SwitchPin: machine.GP21,
			Clock: &Clock{
				Now:        0,
				Previous:   0,
				Transition: 0,
			},
		},
	}

	for {
		for _, encoder := range encoders {
			encoder.Pressed()

			// machine.GP17.SetInterrupt(machine.PinFalling|machine.PinRising, CheckEncoder())
			encoder.CheckEncoder()
		}
	}
}
