package main

import (
	"time"
)

type PinType string
type PinMode int
type PinState int

const (
	PinAnalog  PinType = "analog"
	PinDigital PinType = "digital"
)

const (
	PinInput  PinMode = 0
	PinOutput PinMode = 1
)

const (
	PinLow  PinState = 0
	PinHigh PinState = 1
)

type Pin struct {
	Pin          int
	PinType      PinType
	Mode         PinMode
	State        PinState
	PWMSupported bool
	PortServer   *PortServer
	// DutyCycle float32
	// Period    int

}

func NewPin(pin int, ps *PortServer) *Pin {
	return &Pin{
		Pin:        pin,
		PinType:    PinDigital,
		Mode:       PinInput,
		State:      PinLow,
		PortServer: ps,
	}
}

func (p *Pin) SetPinType(pinType PinType) {
	p.PinType = pinType
}

func (p *Pin) SetPinMode(mode SetupPinRequestMode) error {
	var pinMode PinMode
	if mode == Input {
		pinMode = PinInput
	} else if mode == Output {
		pinMode = PinOutput
	} else {
		return &InvalidModeError{}
	}

	err := p.PortServer.SetupPin(p.Pin, pinMode)
	if err != nil {
		return err
	}

	p.Mode = pinMode

	return nil
}

func (p *Pin) High() error {
	err := p.PortServer.WriteDigitalPin(p.Pin, PinHigh)
	if err != nil {
		return err
	}

	p.State = PinHigh

	return nil
}

func (p *Pin) Low() error {
	err := p.PortServer.WriteDigitalPin(p.Pin, PinLow)
	if err != nil {
		return err
	}

	p.State = PinLow

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

type PWMNotSupportedError struct{}

func (e *PWMNotSupportedError) Error() string {
	return "PWM not supported"
}
