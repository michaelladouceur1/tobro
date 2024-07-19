package main

import (
	"time"
)

type PinType string
type PinMode int

const (
	PinAnalog  PinType = "analog"
	PinDigital PinType = "digital"
)

const (
	PinInput  PinMode = 0
	PinOutput PinMode = 1
)

const DigitalPinLow = 0
const DigitalPinHigh = 1

const AnalogPinMin = 0
const AnalogPinMax = 255

type Pin struct {
	Pin          int
	PinType      PinType
	Mode         PinMode
	State        int
	PWMSupported bool
	PortServer   *PortServer
	// DutyCycle float32
	// Period    int

}

// TODO: Separate the pin into digital and analog pins
func NewPin(pin int, ps *PortServer) *Pin {
	return &Pin{
		Pin:        pin,
		PinType:    PinDigital,
		Mode:       PinInput,
		State:      DigitalPinLow,
		PortServer: ps,
	}
}

func (p *Pin) SetPinType(pinType PinType) {
	p.PinType = pinType
}

func (p *Pin) SetPinMode(mode SetupPinRequestMode) error {
	var pinMode PinMode
	switch mode {
	case Input:
		pinMode = PinInput
	case Output:
		pinMode = PinOutput
	default:
		return &InvalidModeError{}
	}

	err := p.PortServer.SetupPin(p.Pin, pinMode)
	if err != nil {
		return err
	}

	p.Mode = pinMode

	return nil
}

func (p *Pin) SetDigitalPinState(state int) error {
	var err error
	switch state {
	case DigitalPinLow:
		err = p.Low()
	case DigitalPinHigh:
		err = p.High()
	default:
		return &InvalidDigitalStateError{}
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Pin) SetAnalogPinState(state int) error {
	if state < AnalogPinMin || state > AnalogPinMax {
		return &InvalidAnalogStateError{}
	}

	err := p.PortServer.WriteAnalogPin(p.Pin, state)
	if err != nil {
		return err
	}

	return nil
}

func (p *Pin) High() error {
	err := p.PortServer.WriteDigitalPin(p.Pin, DigitalPinHigh)
	if err != nil {
		return err
	}

	p.State = DigitalPinHigh

	return nil
}

func (p *Pin) Low() error {
	err := p.PortServer.WriteDigitalPin(p.Pin, DigitalPinLow)
	if err != nil {
		return err
	}

	p.State = DigitalPinLow

	return nil
}

func (p *Pin) PWM(dutyCycle int, period int, duration int) error {
	if !p.PWMSupported {
		return &PWMNotSupportedError{}
	}

	dutyCycleFloat := float32(dutyCycle) / 100
	timeHigh := int(float32(period) * dutyCycleFloat)
	timeLow := period - timeHigh

	for {

		if duration == 0 {
			break
		}

		err := p.High()
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(timeHigh) * time.Millisecond)

		err = p.Low()
		if err != nil {
			return err
		}

		time.Sleep(time.Duration(timeLow) * time.Millisecond)

		if duration > 0 {
			duration -= period
		}
	}

	return nil
}

type InvalidModeError struct{}

func (e *InvalidModeError) Error() string {
	return "invalid mode"
}

type InvalidDigitalStateError struct{}

func (e *InvalidDigitalStateError) Error() string {
	return "Invalid digital state. Must be 0 or 1"
}

type InvalidAnalogStateError struct{}

func (e *InvalidAnalogStateError) Error() string {
	return "Invalid analog value. Must be between 0 and 255"
}

type PWMNotSupportedError struct{}

func (e *PWMNotSupportedError) Error() string {
	return "PWM not supported"
}
