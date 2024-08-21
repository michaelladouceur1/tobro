package main

import "time"

type Sketch struct {
	Circuit *Circuit
	ID      int
	Name    string
	Steps   []SketchStep
}

type SketchStep struct {
	ID     int
	Start  time.Time
	End    time.Time
	Pin    Pin
	Action SketchAction
}

type SketchAction string

const (
	DigitalWrite SketchAction = "digital_write"
	AnalogWrite  SketchAction = "analog_write"
)

func NewSketch(id int, name string, c *Circuit) *Sketch {
	return &Sketch{
		ID:      id,
		Name:    name,
		Circuit: c,
	}
}
