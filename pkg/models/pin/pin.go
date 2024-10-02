package pin

import (
	"encoding/json"
	"tobro/db"
	"tobro/pkg/arduino"
	"tobro/pkg/models"
)

type Pin struct {
	PortServer   *arduino.PortServer
	PinNumber    int
	PinType      models.PinType
	Min          int
	Max          int
	DigitalRead  bool
	DigitalWrite bool
	AnalogRead   bool
	AnalogWrite  bool
	Mode         models.PinMode
	State        chan int
}

type PinConfig struct {
	PinType      models.PinType
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

func NewPin(ps *arduino.PortServer, pinNumber int, config PinConfig) *Pin {
	var min, max int
	if config.AnalogWrite || config.AnalogRead {
		min, max = models.AnalogPinMin, models.AnalogPinMax
	} else {
		min, max = models.DigitalPinLow, models.DigitalPinHigh
	}

	return &Pin{
		PortServer:   ps,
		PinNumber:    pinNumber,
		PinType:      config.PinType,
		Min:          min,
		Max:          max,
		DigitalRead:  config.DigitalRead,
		DigitalWrite: config.DigitalWrite,
		AnalogRead:   config.AnalogRead,
		AnalogWrite:  config.AnalogWrite,
		Mode:         models.PinInput,
		State:        make(chan int),
	}
}

func (p *Pin) UpdateFromDBModel(model *db.PinDBModel) {
	p.Mode = models.PinMode(model.Mode)
}

func (p *Pin) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"pinNumber":    p.PinNumber,
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

func (p *Pin) SetMode(mode models.PinMode) error {
	err := p.PortServer.SetupPin(p.PinNumber, mode)
	if err != nil {
		return err
	}

	p.Mode = mode

	return nil
}

func (p *Pin) High() error {
	if !p.DigitalWrite {
		return &DigitalWriteNotSupportedError{}
	}

	err := p.PortServer.WriteDigitalPin(p.PinNumber, models.DigitalPinHigh)
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

	err := p.PortServer.WriteDigitalPin(p.PinNumber, models.DigitalPinLow)
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

	err := p.PortServer.WriteAnalogPin(p.PinNumber, state)
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
	case models.DigitalPinLow:
		err = p.Low()
	case models.DigitalPinHigh:
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
