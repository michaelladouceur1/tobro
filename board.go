package main

// arduino nano pinout
// digital pins: 2-13
// analog/digital pins: A0-A5
// analog pins: A6, A7
// pwm pins: 3, 5, 6, 9, 10, 11

type SupportedBoards string

const (
	ArduinoNano SupportedBoards = "arduino_nano"
)

type Board struct {
	Pins []Pin
}

func NewBoard(boardType SupportedBoards, ps *PortServer) *Board {
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
			AnalogRead:  true,
			AnalogWrite: false,
		}
		
		return &Board{
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
	default:
		return nil
	}
}


func (b *Board) PinCount() int {
	return len(b.Pins)
}

func (b *Board) GetState() Board {
	return *b
}

func (b *Board) GetPin(id int) (*Pin, error) {
	for _, p := range b.Pins {
		if p.ID == id {
			return &p, nil
		}
	}

	return nil, &PinNotFoundError{}
}

func (b *Board) GetPins() []*Pin {
	var pins []*Pin
	for _, p := range b.Pins {
		pins = append(pins, &p)
	}

	return pins
}

func (b *Board) GetDigitalWritePin(id int) (DigitalWritePin, error) {
	for _, p := range b.Pins {
		if p.ID == id {
			if p.DigitalWrite {
				return &p, nil
			}

			return nil, &PinNotDigitalError{}
		}
	
	}

	return nil, &PinNotFoundError{}
}

func (b *Board) GetAnalogWritePin(id int) (AnalogWritePin, error) {
	for _, p := range b.Pins {
		if p.ID == id {
			if p.AnalogWrite {
				return &p, nil
			}

			return nil, &PinNotAnalogError{}
		}
	}

	return nil, &PinNotFoundError{}
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
