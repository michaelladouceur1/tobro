package pin

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
