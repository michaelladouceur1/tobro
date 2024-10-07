package sketch

import (
	"tobro/db"
	"tobro/internal/models"
	"tobro/internal/models/circuit"
	"tobro/internal/models/pin"
)

type Sketch struct {
	Circuit *circuit.Circuit
	ID      int
	Name    string
	Steps   []SketchStep
}

type SketchStep struct {
	ID     int
	Start  int
	End    int
	Pin    *pin.Pin
	Action models.SketchAction
}

func New(id int, name string, c *circuit.Circuit) *Sketch {
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
			Action: models.SketchAction(step.Action),
		})
	}
}
