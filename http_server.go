package main

import (
	"encoding/json"
	"net/http"
)

type HTTPServer struct {
	session *Session
	dal     *DAL
	circuit *Circuit
	sketch  *Sketch
}

func NewHTTPServer(session *Session, dal *DAL, circuit *Circuit, sketch *Sketch) *HTTPServer {
	return &HTTPServer{
		session: session,
		dal:     dal,
		circuit: circuit,
		sketch:  sketch,
	}
}

func pinResponseFromPin(pin Pin) PinResponse {
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

func pinResponseFromPins(pins []Pin) []PinResponse {
	pinResponses := make([]PinResponse, 0)
	for _, pin := range pins {
		pinResponses = append(pinResponses, pinResponseFromPin(pin))
	}
	return pinResponses
}

func circuitResponseFromCircuit(circuit *Circuit) CircuitResponse {
	pinResponses := pinResponseFromPins(circuit.Pins)
	return CircuitResponse{
		Id:    circuit.ID,
		Name:  circuit.Name,
		Board: string(circuit.Board),
		Pins:  pinResponses,
	}
}

func (s *HTTPServer) GetPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SuccessResponse{Message: "pong"})
}

func (s *HTTPServer) PostConnect(w http.ResponseWriter, r *http.Request) {
	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := s.circuit.Connect(req.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.session.UpdatePort(req.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ConnectResponse{Port: &req.Port})
}

func (s *HTTPServer) GetBoards(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(BoardsResponse{Boards: supportedBoards})
}

// Circuit

func (s *HTTPServer) GetCircuit(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(circuitResponseFromCircuit(s.circuit))
}

func (s *HTTPServer) PostCircuit(w http.ResponseWriter, r *http.Request) {
	var req CreateCircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newCircuit, err := s.dal.CreateCircuit(req.Name, req.Board)
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

	newCircuit, err := s.dal.SaveCircuit(*s.circuit)
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

	pin, err := s.circuit.GetPin(req.PinNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetMode(req.Mode)
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

// Sketch
