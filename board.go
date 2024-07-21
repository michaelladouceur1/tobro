package main

// arduino nano pinout
// digital pins: 2-13
// analog pins: A0-A7
// pwm pins: 3, 5, 6, 9, 10, 11

type SupportedBoards string

const (
	ArduinoNano SupportedBoards = "arduino_nano"
)

type Board struct {
	DigitalPins []DigitalPin
	AnalogPins  []AnalogPin
}

func NewBoard(boardType SupportedBoards, ps *PortServer) *Board {
	switch boardType {
	case ArduinoNano:
		return &Board{
			DigitalPins: []DigitalPin{
				*NewDigitalPin(2, false, ps),
				*NewDigitalPin(3, true, ps),
				*NewDigitalPin(4, false, ps),
				*NewDigitalPin(5, true, ps),
				*NewDigitalPin(6, true, ps),
				*NewDigitalPin(7, false, ps),
				*NewDigitalPin(8, false, ps),
				*NewDigitalPin(9, true, ps),
				*NewDigitalPin(10, true, ps),
				*NewDigitalPin(11, true, ps),
				*NewDigitalPin(12, false, ps),
				*NewDigitalPin(13, false, ps),
			},
			AnalogPins: []AnalogPin{
				*NewAnalogPin(14, ps),
				*NewAnalogPin(15, ps),
				*NewAnalogPin(16, ps),
				*NewAnalogPin(17, ps),
				*NewAnalogPin(18, ps),
				*NewAnalogPin(19, ps),
				*NewAnalogPin(20, ps),
				*NewAnalogPin(21, ps),
			},
		}
	default:
		return nil
	}
}

func (b *Board) PinCount() int {
	return len(b.DigitalPins) + len(b.AnalogPins)
}

func (b *Board) GetState() Board {
	return *b
}

func (b *Board) GetPin(id int) (Pin, error) {
	digitalPin := b.GetDigitalPin(id)
	if digitalPin != nil {
		return digitalPin, nil
	}

	analogPin := b.GetAnalogPin(id)
	if analogPin != nil {
		return analogPin, nil
	}

	return nil, &PinNotFoundError{}

}

func (b *Board) GetDigitalWritePin(id int) (DigitalWritePin, error) {
	pin := b.GetDigitalPin(id)
	if pin == nil {
		return nil, &PinNotFoundError{}
	}

	return pin, nil
}

func (b *Board) GetAnalogWritePin(id int) (AnalogWritePin, error) {
	analogPin := b.GetAnalogPin(id)
	if analogPin != nil {
		return analogPin, nil
	}

	pwmPin := b.GetPWMPin(id)
	if pwmPin != nil {
		return pwmPin, nil
	}

	return nil, &PinNotFoundError{}
}

func (b *Board) GetDigitalPin(id int) *DigitalPin {
	for _, p := range b.DigitalPins {
		if p.ID == id {
			if p.PinType == PinDigital {
				return &p
			}

			return nil
		}
	}

	return nil
}

func (b *Board) GetAnalogPin(id int) *AnalogPin {
	for _, p := range b.AnalogPins {
		if p.ID == id {
			if p.PinType == PinAnalog {
				return &p
			}

			return nil
		}
	}

	return nil
}

func (b *Board) GetPWMPin(id int) *DigitalPin {
	for _, p := range b.DigitalPins {
		if p.ID == id {
			if p.PWM {
				return &p
			}

			return nil
		}
	}

	return nil

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
