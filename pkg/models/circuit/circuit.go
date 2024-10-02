package circuit

import (
	"tobro/db"
	"tobro/pkg/arduino"
	"tobro/pkg/models"
	"tobro/pkg/models/pin"
)

// arduino nano pinout
// digital pins: 2-13
// analog/digital pins: A0-A5
// analog pins: A6, A7
// pwm pins: 3, 5, 6, 9, 10, 11

type Circuit struct {
	portServer *arduino.PortServer
	ID         int
	Name       string
	Board      models.SupportedBoards
	Pins       []pin.Pin
}

func NewCircuit(id int, name string, boardType models.SupportedBoards, ps *arduino.PortServer) *Circuit {
	switch boardType {
	case models.ArduinoNano:
		digitalPinConfig := pin.PinConfig{
			PinType:      models.PinDigital,
			DigitalRead:  true,
			DigitalWrite: true,
			AnalogRead:   false,
			AnalogWrite:  false,
		}

		digitalPwmPinConfig := pin.PinConfig{
			PinType:      models.PinDigital,
			DigitalRead:  true,
			DigitalWrite: true,
			AnalogRead:   false,
			AnalogWrite:  true,
		}

		analogDigitalPinConfig := pin.PinConfig{
			PinType:      models.PinAnalog,
			DigitalRead:  true,
			DigitalWrite: true,
			AnalogRead:   true,
			AnalogWrite:  true,
		}

		analogPinConfig := pin.PinConfig{
			PinType:      models.PinAnalog,
			DigitalRead:  false,
			DigitalWrite: false,
			AnalogRead:   true,
			AnalogWrite:  false,
		}

		c := &Circuit{
			portServer: ps,
			ID:         id,
			Name:       name,
			Board:      models.ArduinoNano,
			Pins: []pin.Pin{
				*pin.NewPin(ps, 2, digitalPinConfig),
				*pin.NewPin(ps, 3, digitalPwmPinConfig),
				*pin.NewPin(ps, 4, digitalPinConfig),
				*pin.NewPin(ps, 5, digitalPwmPinConfig),
				*pin.NewPin(ps, 6, digitalPwmPinConfig),
				*pin.NewPin(ps, 7, digitalPinConfig),
				*pin.NewPin(ps, 8, digitalPinConfig),
				*pin.NewPin(ps, 9, digitalPwmPinConfig),
				*pin.NewPin(ps, 10, digitalPwmPinConfig),
				*pin.NewPin(ps, 11, digitalPwmPinConfig),
				*pin.NewPin(ps, 12, digitalPinConfig),
				*pin.NewPin(ps, 13, digitalPinConfig),
				*pin.NewPin(ps, 14, analogDigitalPinConfig),
				*pin.NewPin(ps, 15, analogDigitalPinConfig),
				*pin.NewPin(ps, 16, analogDigitalPinConfig),
				*pin.NewPin(ps, 17, analogDigitalPinConfig),
				*pin.NewPin(ps, 18, analogDigitalPinConfig),
				*pin.NewPin(ps, 19, analogDigitalPinConfig),
				*pin.NewPin(ps, 20, analogPinConfig),
				*pin.NewPin(ps, 21, analogPinConfig),
			},
		}

		go c.watchPortConnection()

		return c
	default:
		return nil
	}
}

func GetSupportedBoards() []string {
	return models.SupportedBoardsList
}

func SupportedBoardPins(board string) ([]pin.Pin, error) {
	circuit := NewCircuit(0, "", models.SupportedBoards(board), nil)
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
	c.Board = models.SupportedBoards(model.Board)

	pins := model.Pins()
	for _, pin := range pins {
		cPin, err := c.GetPin(pin.PinNumber)
		if err != nil {
			continue
		}

		cPin.UpdateFromDBModel(&pin)
	}
}

func (c *Circuit) GetPin(pinNumber int) (*pin.Pin, error) {
	for i, p := range c.Pins {
		if p.PinNumber == pinNumber {
			return &c.Pins[i], nil
		}
	}

	return nil, &PinNotFoundError{}
}

func (c *Circuit) GetPins() []*pin.Pin {
	var pins []*pin.Pin
	for _, p := range c.Pins {
		pins = append(pins, &p)
	}

	return pins
}

func (c *Circuit) GetDigitalWritePin(pinNumber int) (pin.DigitalWritePin, error) {
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

func (c *Circuit) GetAnalogWritePin(pinNumber int) (pin.AnalogWritePin, error) {
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
