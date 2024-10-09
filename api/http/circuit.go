package http_server

import (
	"encoding/json"
	"net/http"
	"tobro/internal/models"
	"tobro/internal/models/circuit"
	"tobro/internal/models/pin"
)

func pinResponseFromPin(pin pin.Pin) PinResponse {
	return PinResponse{
		PinNumber:    pin.PinNumber,
		Type:         string(pin.PinType),
		Mode:         int(pin.Mode),
		Min:          pin.Min,
		Max:          pin.Max,
		DigitalRead:  pin.DigitalRead,
		DigitalWrite: pin.DigitalWrite,
		AnalogRead:   pin.AnalogRead,
		AnalogWrite:  pin.AnalogWrite,
	}
}

func pinResponseFromPins(pins []pin.Pin) []PinResponse {
	pinResponses := make([]PinResponse, 0)
	for _, pin := range pins {
		pinResponses = append(pinResponses, pinResponseFromPin(pin))
	}
	return pinResponses
}

func circuitResponseFromCircuit(circuit *circuit.Circuit) CircuitResponse {
	pinResponses := pinResponseFromPins(circuit.Pins)
	return CircuitResponse{
		Id:    circuit.ID,
		Name:  circuit.Name,
		Board: string(circuit.Board),
		Pins:  pinResponses,
	}
}

func (s *HTTPServer) GetCircuit(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(circuitResponseFromCircuit(s.circuit))
}

func (s *HTTPServer) PostCircuit(w http.ResponseWriter, r *http.Request) {
	var req CreateCircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCircuit, err := s.store.CreateCircuit(req.Name, req.Board)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.circuit.UpdateFromDBModel(newCircuit)

	json.NewEncoder(w).Encode(circuitResponseFromCircuit(s.circuit))
}

func (s *HTTPServer) PostSaveCircuit(w http.ResponseWriter, r *http.Request) {
	var req SaveCircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if req.Id != s.circuit.ID {
	// 	http.Error(w, "invalid circuit id", http.StatusBadRequest)
	// 	return
	// }

	newCircuit, err := s.store.SaveCircuit(*s.circuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.circuit.UpdateFromDBModel(newCircuit)

	json.NewEncoder(w).Encode(circuitResponseFromCircuit(s.circuit))
}

func (s *HTTPServer) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	var req SetupPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	p, err := s.circuit.GetPin(req.PinNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch req.Mode {
	case Input:
		err = p.SetMode(models.PinInput)
	case Output:
		err = p.SetMode(models.PinOutput)
	default:
		http.Error(w, "invalid mode", http.StatusBadRequest)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(SetupPinResponse{Mode: string(req.Mode), PinNumber: req.PinNumber})
}

func (s *HTTPServer) PostDigitalWritePin(w http.ResponseWriter, r *http.Request) {
	var req DigitalWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dwPin, err := s.circuit.GetDigitalWritePin(req.PinNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = dwPin.SetDigitalState(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DigitalWritePinResponse{PinNumber: req.PinNumber, Value: req.Value})
}

func (s *HTTPServer) PostAnalogWritePin(w http.ResponseWriter, r *http.Request) {
	var req AnalogWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := s.circuit.GetAnalogWritePin(req.PinNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetAnalogState(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AnalogWritePinResponse{PinNumber: &req.PinNumber, Value: &req.Value})
}
