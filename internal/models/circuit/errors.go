package circuit

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
