package constants

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
