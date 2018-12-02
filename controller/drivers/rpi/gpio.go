package rpi

import (
	"github.com/kidoman/embd"
	"github.com/pkg/errors"
	"github.com/reef-pi/reef-pi/controller/drivers"
)

var (
	validGPIOPins = map[int]bool{
		2:  true,
		3:  true,
		4:  true,
		5:  true,
		6:  true,
		7:  true,
		8:  true,
		9:  true,
		10: true,
		11: true,
		12: true,
		13: true,
		14: true,
		15: true,
		16: true,
		17: true,
		18: true,
		19: true,
		20: true,
		21: true,
		22: true,
		23: true,
		24: true,
		25: true,
		26: true,
		27: true,
	}
)

type rpiPin struct {
	pin  int
	name string

	digitalPin embd.DigitalPin
}

func (p *rpiPin) Name() string {
	return p.name
}

func (p *rpiPin) Read() (bool, error) {
	err := p.digitalPin.SetDirection(embd.In)
	if err != nil {
		return false, errors.Wrapf(err, "can't read input from pin %d", p.pin)
	}

	v, err := p.digitalPin.Read()
	if err != nil {
		return false, err
	}
	return v == 1, nil
}

func (r *rpiDriver) InputPins() []drivers.InputPin {
	var pins []drivers.InputPin
	for _, pin := range r.pins {
		pins = append(pins, pin)
	}
	return pins
}
