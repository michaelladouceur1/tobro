package main

import "encoding/json"

type PinType string
type PinMode int

const (
	PinAnalog  PinType = "analog"
	PinDigital PinType = "digital"

	PinInput  PinMode = 0
	PinOutput PinMode = 1

	DigitalPinLow  = 0
	DigitalPinHigh = 1

	AnalogPinMin = 0
	AnalogPinMax = 255
)

type Pin struct {
	PortServer   *PortServer
	ID           int
	PinType      PinType
	Min          int
	Max          int
	DigitalRead  bool
	DigitalWrite bool
	AnalogRead   bool
	AnalogWrite  bool
	Mode         PinMode
	State        chan int
}

type PinConfig struct {
	PinType      PinType
	DigitalRead  bool
	DigitalWrite bool
	AnalogRead   bool
	AnalogWrite  bool
}

type DigitalWritePin interface {
	High() error
	Low() error
	SetDigitalState(state int) error
}

type AnalogWritePin interface {
	SetAnalogState(state int) error
}

func NewPin(ps *PortServer, id int, config PinConfig) *Pin {
	var min, max int
	if config.AnalogWrite || config.AnalogRead {
		min, max = AnalogPinMin, AnalogPinMax
	} else {
		min, max = DigitalPinLow, DigitalPinHigh
	}

	return &Pin{
		PortServer:   ps,
		ID:           id,
		PinType:      config.PinType,
		Min:          min,
		Max:          max,
		DigitalRead:  config.DigitalRead,
		DigitalWrite: config.DigitalWrite,
		AnalogRead:   config.AnalogRead,
		AnalogWrite:  config.AnalogWrite,
		Mode:         PinInput,
		State:        make(chan int),
	}
}

func (p *Pin) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":           p.ID,
		"type":         p.PinType,
		"min":          p.Min,
		"max":          p.Max,
		"digitalRead":  p.DigitalRead,
		"digitalWrite": p.DigitalWrite,
		"analogRead":   p.AnalogRead,
		"analogWrite":  p.AnalogWrite,
		"mode":         p.Mode,
	})
}

func (p *Pin) SetMode(mode SetupPinRequestMode) error {
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

func (p *Pin) High() error {
	if !p.DigitalWrite {
		return &DigitalWriteNotSupportedError{}
	}

	err := p.PortServer.WriteDigitalPin(p.ID, DigitalPinHigh)
	if err != nil {
		return err
	}

	p.State <- p.Max

	return nil
}

func (p *Pin) Low() error {
	if !p.DigitalWrite {
		return &DigitalWriteNotSupportedError{}
	}

	err := p.PortServer.WriteDigitalPin(p.ID, DigitalPinLow)
	if err != nil {
		return err
	}

	p.State <- p.Min

	return nil
}

func (p *Pin) SetAnalogState(state int) error {
	if !p.AnalogWrite {
		return &AnalogWriteNotSupportedError{}
	}

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

func (p *Pin) SetDigitalState(state int) error {
	if !p.DigitalWrite {
		return &DigitalWriteNotSupportedError{}
	}

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

type DigitalWriteNotSupportedError struct{}

func (e *DigitalWriteNotSupportedError) Error() string {
	return "Digital write not supported"
}

type AnalogWriteNotSupportedError struct{}

func (e *AnalogWriteNotSupportedError) Error() string {
	return "Analog write not supported"
}
