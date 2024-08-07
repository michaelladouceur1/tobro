package main

import (
	"encoding/json"
	"net/http"
)

type HTTPServer struct{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func pinResponseFromPin(pin Pin) PinResponse {
	return PinResponse{
		Id:           pin.ID,
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

func (s *HTTPServer) GetPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SuccessResponse{Message: "pong"})
}

func (s *HTTPServer) PostConnect(w http.ResponseWriter, r *http.Request) {
	var req ConnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := portServer.OpenPort(req.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	portServer.ListenToPort()

	json.NewEncoder(w).Encode(ConnectResponse{Port: &req.Port})
}

func (s *HTTPServer) GetCircuit(w http.ResponseWriter, r *http.Request) {
	pinResponses := pinResponseFromPins(circuit.Pins)
	json.NewEncoder(w).Encode(CircuitResponse{Pins: pinResponses})
}

func (s *HTTPServer) PostCircuit(w http.ResponseWriter, r *http.Request) {
	var req CreateCircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	boardType, ok := req.Board.(SupportedBoards)
	if !ok {
		http.Error(w, "invalid board type", http.StatusBadRequest)
		return
	}

	circuit = NewCircuit(boardType, portServer)

	createdBoard, err := dal.CreateCircuit(req.Name, boardType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = dal.AddPins(createdBoard.ID, circuit.Pins)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pinResponses := pinResponseFromPins(circuit.Pins)
	json.NewEncoder(w).Encode(CircuitResponse{Pins: pinResponses})
}

func (s *HTTPServer) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	var req SetupPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := circuit.GetPin(req.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetMode(req.Mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(SetupPinResponse{Mode: string(req.Mode), Pin: req.Pin})
}

func (s *HTTPServer) PostDigitalWritePin(w http.ResponseWriter, r *http.Request) {
	var req DigitalWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := circuit.GetDigitalWritePin(req.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetDigitalState(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DigitalWritePinResponse{Pin: &req.Pin, Value: &req.Value})
}

func (s *HTTPServer) PostAnalogWritePin(w http.ResponseWriter, r *http.Request) {
	var req AnalogWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := circuit.GetAnalogWritePin(req.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetAnalogState(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(AnalogWritePinResponse{Pin: &req.Pin, Value: &req.Value})
}
