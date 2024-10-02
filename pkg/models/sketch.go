package models

import "tobro/db"

type Sketch struct {
	Circuit *Circuit
	ID      int
	Name    string
	Steps   []SketchStep
}

type SketchStep struct {
	ID     int
	Start  int
	End    int
	Pin    *Pin
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

func (s *Sketch) UpdateFromDBModel(model *db.SketchDBModel) {
	s.ID = model.ID
	s.Name = model.Name
	s.Steps = make([]SketchStep, 0)

	steps := model.Steps()
	for _, step := range steps {
		pin, err := s.Circuit.GetPin(step.Pin().PinNumber)
		if err != nil {
			continue
		}
		s.Steps = append(s.Steps, SketchStep{
			ID:     step.ID,
			Start:  step.Start,
			End:    step.End,
			Pin:    pin,
			Action: SketchAction(step.Action),
		})
	}
}
