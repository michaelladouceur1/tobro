package http_server

import (
	"encoding/json"
	"net/http"
	"tobro/internal/models/sketch"
)

func apiSketchStepsFromSketchSteps(steps []sketch.SketchStep) []SketchStepAPI {
	apiSketchSteps := make([]SketchStepAPI, 0)
	for _, step := range steps {
		apiSketchSteps = append(apiSketchSteps, SketchStepAPI{
			Id:        step.ID,
			Start:     step.Start,
			End:       step.End,
			PinNumber: step.Pin.PinNumber,
			Action:    string(step.Action),
		})
	}

	return apiSketchSteps
}

func apiSketchFromSketch(sketch *sketch.Sketch) SketchAPI {
	apiSketchSteps := apiSketchStepsFromSketchSteps(sketch.Steps)
	return SketchAPI{
		Id:    sketch.ID,
		Name:  sketch.Name,
		Steps: apiSketchSteps,
	}
}

func (s *HTTPServer) GetSketch(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(apiSketchFromSketch(s.sketch))
}

func (s *HTTPServer) PostSketch(w http.ResponseWriter, r *http.Request) {
	var req SketchAPI
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newSketch, err := s.store.CreateSketch(s.circuit.ID, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.sketch.UpdateFromDBModel(newSketch)

	json.NewEncoder(w).Encode(apiSketchFromSketch(s.sketch))
}
