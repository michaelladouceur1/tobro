package main

import (
	"tobro/db"
)

// arduino nano pinout
// digital pins: 2-13
// analog/digital pins: A0-A5
// analog pins: A6, A7
// pwm pins: 3, 5, 6, 9, 10, 11

type SupportedBoards string

var supportedBoards = []string{
	string(ArduinoNano),
}

const (
	ArduinoNano SupportedBoards = "arduino_nano"
)

type Circuit struct {
	portServer *PortServer
	ID         int
	Name       string
	Board      SupportedBoards
	Pins       []Pin
}

func NewCircuit(id int, name string, boardType SupportedBoards, ps *PortServer) *Circuit {
	switch boardType {
	case ArduinoNano:
		digitalPinConfig := PinConfig{
			PinType:      PinDigital,
			DigitalRead:  true,
			DigitalWrite: true,
			AnalogRead:   false,
			AnalogWrite:  false,
		}

		digitalPwmPinConfig := PinConfig{
			PinType:      PinDigital,
			DigitalRead:  true,
			DigitalWrite: true,
			AnalogRead:   false,
			AnalogWrite:  true,
		}

		analogDigitalPinConfig := PinConfig{
			PinType:      PinAnalog,
			DigitalRead:  true,
			DigitalWrite: true,
			AnalogRead:   true,
			AnalogWrite:  true,
		}

		analogPinConfig := PinConfig{
			PinType:      PinAnalog,
			DigitalRead:  false,
			DigitalWrite: false,
			AnalogRead:   true,
			AnalogWrite:  false,
		}

		c := &Circuit{
			portServer: ps,
			ID:         id,
			Name:       name,
			Board:      ArduinoNano,
			Pins: []Pin{
				*NewPin(ps, 2, digitalPinConfig),
				*NewPin(ps, 3, digitalPwmPinConfig),
				*NewPin(ps, 4, digitalPinConfig),
				*NewPin(ps, 5, digitalPwmPinConfig),
				*NewPin(ps, 6, digitalPwmPinConfig),
				*NewPin(ps, 7, digitalPinConfig),
				*NewPin(ps, 8, digitalPinConfig),
				*NewPin(ps, 9, digitalPwmPinConfig),
				*NewPin(ps, 10, digitalPwmPinConfig),
				*NewPin(ps, 11, digitalPwmPinConfig),
				*NewPin(ps, 12, digitalPinConfig),
				*NewPin(ps, 13, digitalPinConfig),
				*NewPin(ps, 14, analogDigitalPinConfig),
				*NewPin(ps, 15, analogDigitalPinConfig),
				*NewPin(ps, 16, analogDigitalPinConfig),
				*NewPin(ps, 17, analogDigitalPinConfig),
				*NewPin(ps, 18, analogDigitalPinConfig),
				*NewPin(ps, 19, analogDigitalPinConfig),
				*NewPin(ps, 20, analogPinConfig),
				*NewPin(ps, 21, analogPinConfig),
			},
		}

		go c.watchPortConnection()

		return c
	default:
		return nil
	}
}

func SupportedBoardPins(board string) ([]Pin, error) {
	circuit := NewCircuit(0, "", SupportedBoards(board), nil)
	if circuit == nil {
		return nil, &UnsupportedBoardError{}
	}
	return circuit.Pins, nil
}

func (c *Circuit) Connect(port string) error {
	err := c.portServer.OpenPort(port)
	if err != nil {
		return err
	}

	c.portServer.ListenToPort()

	return nil
}

func (c *Circuit) UpdateFromDBModel(model *db.CircuitDBModel) {
	c.ID = model.ID
	c.Name = model.Name
	c.Board = SupportedBoards(model.Board)

	pins := model.Pins()
	for _, pin := range pins {
		cPin, err := c.GetPin(pin.PinNumber)
		if err != nil {
			continue
		}

		cPin.UpdateFromDBModel(&pin)
	}
}

func (c *Circuit) GetPin(pinNumber int) (*Pin, error) {
	for i, p := range c.Pins {
		if p.PinNumber == pinNumber {
			return &c.Pins[i], nil
		}
	}

	return nil, &PinNotFoundError{}
}

func (c *Circuit) GetPins() []*Pin {
	var pins []*Pin
	for _, p := range c.Pins {
		pins = append(pins, &p)
	}

	return pins
}

func (c *Circuit) GetDigitalWritePin(pinNumber int) (DigitalWritePin, error) {
	for i, p := range c.Pins {
		if p.PinNumber == pinNumber {
			if p.DigitalWrite {
				return &c.Pins[i], nil
			}

			return nil, &PinNotDigitalError{}
		}

	}

	return nil, &PinNotFoundError{}
}

func (c *Circuit) GetAnalogWritePin(pinNumber int) (AnalogWritePin, error) {
	for i, p := range c.Pins {
		if p.PinNumber == pinNumber {
			if p.AnalogWrite {
				return &c.Pins[i], nil
			}

			return nil, &PinNotAnalogError{}
		}
	}

	return nil, &PinNotFoundError{}
}

func (c *Circuit) Reset() {
	for i := range c.Pins {
		c.Pins[i].State <- 0
	}
}

func (c *Circuit) watchPortConnection() {
	for {
		connected := <-c.portServer.Connected
		if !connected {
			c.Reset()
		}
	}
}

type PinNotFoundError struct{}

func (e *PinNotFoundError) Error() string {
	return "Pin not found"
}

type PinNotDigitalError struct{}

func (e *PinNotDigitalError) Error() string {
	return "Pin is not digital"
}

type PinNotAnalogError struct{}

func (e *PinNotAnalogError) Error() string {
	return "Pin is not analog"
}

type UnsupportedBoardError struct{}

func (e *UnsupportedBoardError) Error() string {
	return "Unsupported board"
}
