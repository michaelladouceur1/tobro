package models

// Circuit

type SupportedBoards string

var SupportedBoardsList = []string{
	string(ArduinoNano),
}

const (
	ArduinoNano SupportedBoards = "arduino_nano"
)

// Pin

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

// Sketch

type SketchAction string

const (
	DigitalWrite SketchAction = "digital_write"
	AnalogWrite  SketchAction = "analog_write"
)
