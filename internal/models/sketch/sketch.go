package sketch

import (
	"encoding/json"
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

func (s *Sketch) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"id":    s.ID,
		"name":  s.Name,
		"steps": s.Steps,
	})
}

func (s *Sketch) GetSteps() []SketchStep {
	return s.Steps
}

func (s *Sketch) GetStep(id int) (*SketchStep, error) {
	for i, step := range s.Steps {
		if step.ID == id {
			return &s.Steps[i], nil
		}
	}

	return nil, &StepNotFoundError{}
}

func (s *Sketch) AddStep(start, end int, pin *pin.Pin, action models.SketchAction) {
	s.Steps = append(s.Steps, SketchStep{
		Start:  start,
		End:    end,
		Pin:    pin,
		Action: action,
	})
}

func (s *Sketch) RemoveStep(id int) error {
	for i, step := range s.Steps {
		if step.ID == id {
			s.Steps = append(s.Steps[:i], s.Steps[i+1:]...)
			return nil
		}
	}

	return &StepNotFoundError{}
}
