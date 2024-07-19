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
	Pins []Pin
}

func NewBoard(boardType SupportedBoards, ps *PortServer) *Board {
	// TODO: move pin configuration to a separate file
	switch boardType {
	case ArduinoNano:
		return &Board{
			Pins: []Pin{
				{
					Pin:          2,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          3,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          4,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          5,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          6,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          7,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          8,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          9,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          10,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          11,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          12,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          13,
					PinType:      PinDigital,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: true,
					PortServer:   ps,
				},
				{
					Pin:          14,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          15,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          16,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          17,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          18,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          19,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          20,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
				{
					Pin:          21,
					PinType:      PinAnalog,
					Mode:         PinInput,
					State:        DigitalPinLow,
					PWMSupported: false,
					PortServer:   ps,
				},
			},
		}
	default:
		return nil

	}
}

func (b *Board) PinCount() int {
	return len(b.Pins)
}

func (b *Board) GetPin(pin int) (*Pin, error) {
	for _, p := range b.Pins {
		if p.Pin == pin {
			return &p, nil
		}
	}

	return nil, &PinNotFoundError{}
}

func (b *Board) GetDigitalPin(pin int) (*Pin, error) {
	for _, p := range b.Pins {
		if p.Pin == pin {
			if p.PinType == PinDigital {
				return &p, nil
			}

			return nil, &PinNotDigitalError{}
		}
	}

	return nil, &PinNotFoundError{}
}

func (b *Board) GetAnalogPin(pin int) (*Pin, error) {
	for _, p := range b.Pins {
		if p.Pin == pin {
			if p.PinType == PinAnalog {
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
