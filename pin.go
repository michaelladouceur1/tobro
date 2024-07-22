package main

import "encoding/json"

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

type Pin interface {
	SetMode(mode SetupPinRequestMode) error
	High() error
	Low() error
	// SetState(state int) error
	// SetAnalogState(state int) error
}

type DigitalWritePin interface {
	SetDigitalState(state int) error
}

type AnalogWritePin interface {
	SetAnalogState(state int) error
}

type pin struct {
	PortServer *PortServer
	ID         int
	PinType    PinType
	Mode       PinMode
	Min        int
	Max        int
	State      chan int
}

type DigitalPin struct {
	PWM bool
	pin
}

type AnalogPin struct {
	pin
}

func NewDigitalPin(id int, pwm bool, ps *PortServer) *DigitalPin {
	var min, max int
	if pwm {
		min, max = AnalogPinMin, AnalogPinMax
	} else {
		min, max = DigitalPinLow, DigitalPinHigh
	}

	return &DigitalPin{
		pin: pin{
			PortServer: ps,
			ID:         id,
			PinType:    PinDigital,
			Mode:       PinInput,
			Min:        min,
			Max:        max,
			State:      make(chan int),
		},
		PWM: pwm,
	}
}

func NewAnalogPin(id int, ps *PortServer) *AnalogPin {
	return &AnalogPin{
		pin: pin{
			PortServer: ps,
			ID:         id,
			PinType:    PinAnalog,
			Mode:       PinInput,
			Min:        AnalogPinMin,
			Max:        AnalogPinMax,
			State:      make(chan int),
		},
	}
}

func (p *DigitalPin) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":   p.ID,
		"type": p.PinType,
		"mode": p.Mode,
		"pwm":  p.PWM,
	})
}

func (p *AnalogPin) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":   p.ID,
		"type": p.PinType,
		"mode": p.Mode,
	})
}

func (p *pin) SetMode(mode SetupPinRequestMode) error {
	var pinMode PinMode
	switch mode {
	case Input:
		pinMode = PinInput
	case Output:
		pinMode = PinOutput
	default:
		return &InvalidModeError{}
	}

	err := p.PortServer.SetupPin(p.ID, pinMode)
	if err != nil {
		return err
	}

	p.Mode = pinMode

	return nil
}

func (p *pin) DigitalValid() bool {
	return p.PinType == PinDigital
}

func (p *DigitalPin) High() error {
	err := p.PortServer.WriteDigitalPin(p.ID, p.Max)
	if err != nil {
		return err
	}

	p.State <- p.Max

	return nil
}

func (p *DigitalPin) Low() error {
	err := p.PortServer.WriteDigitalPin(p.ID, p.Min)
	if err != nil {
		return err
	}

	p.State <- p.Min

	return nil
}

func (p *AnalogPin) High() error {
	err := p.PortServer.WriteAnalogPin(p.ID, p.Max)
	if err != nil {
		return err
	}

	p.State <- p.Max

	return nil
}

func (p *AnalogPin) Low() error {
	err := p.PortServer.WriteAnalogPin(p.ID, p.Min)
	if err != nil {
		return err
	}

	p.State <- p.Min

	return nil
}

func (p *DigitalPin) SetAnalogState(state int) error {
	if !p.PWM {
		return &PWMNotSupportedError{}
	}
	return p.setAnalogState(state)
}

func (p *AnalogPin) SetAnalogState(state int) error {
	return p.setAnalogState(state)
}

func (p *DigitalPin) SetDigitalState(state int) error {
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

func (p *pin) setAnalogState(state int) error {
	if state < p.Min || state > p.Max {
		return &InvalidAnalogStateError{}
	}

	err := p.PortServer.WriteAnalogPin(p.ID, state)
	if err != nil {
		return err
	}

	p.State <- state

	return nil
}

// func (p *DigitalPin) PWM(dutyCycle int, period int, duration int) error {
// 	if !p.PWM {
// 		return &PWMNotSupportedError{}
// 	}

// 	dutyCycleFloat := float32(dutyCycle) / 100
// 	timeHigh := int(float32(period) * dutyCycleFloat)
// 	timeLow := period - timeHigh

// 	for {

// 		if duration == 0 {
// 			break
// 		}

// 		err := p.High()
// 		if err != nil {
// 			return err
// 		}

// 		time.Sleep(time.Duration(timeHigh) * time.Millisecond)

// 		err = p.Low()
// 		if err != nil {
// 			return err
// 		}

// 		time.Sleep(time.Duration(timeLow) * time.Millisecond)

// 		if duration > 0 {
// 			duration -= period
// 		}
// 	}

// 	return nil
// }

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
