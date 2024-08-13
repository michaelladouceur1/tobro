package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type HTTPServer struct{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
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

	err := portServer.OpenPort(req.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	portServer.ListenToPort()

	json.NewEncoder(w).Encode(ConnectResponse{Port: &req.Port})
}

func (s *HTTPServer) GetBoards(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(BoardsResponse{Boards: supportedBoards})
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

	newCircuit, err := dal.CreateCircuit(*NewCircuit(0, req.Name, SupportedBoards(req.Board), portServer))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	circuit.UpdateCircuit(newCircuit)

	log.Print(circuit.ID)
	log.Print(circuit.Name)
	log.Print(circuit.Board)

	json.NewEncoder(w).Encode(circuitResponseFromCircuit(circuit))
}

func (s *HTTPServer) PostSaveCircuit(w http.ResponseWriter, r *http.Request) {
	var req SaveCircuitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Print(req.Id)
	log.Print(circuit.ID)
	log.Print(circuit.Name)
	log.Print(circuit.Board)

	if req.Id != circuit.ID {
		http.Error(w, "invalid circuit id", http.StatusBadRequest)
		return
	}

	for _, pin := range circuit.Pins {
		log.Print(pin.Mode)
	}

	newCircuit, err := dal.SaveCircuit(*circuit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// for _, pin := range newCircuit.Pins {
	// 	log.Print(pin.Mode)
	// }

	circuit.UpdateCircuit(newCircuit)

	json.NewEncoder(w).Encode(circuitResponseFromCircuit(circuit))
}

func (s *HTTPServer) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	var req SetupPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := circuit.GetPin(req.PinNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetMode(req.Mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, pin := range circuit.Pins {
		log.Print(pin.Mode)
	}

	json.NewEncoder(w).Encode(SetupPinResponse{Mode: string(req.Mode), PinNumber: req.PinNumber})
}

func (s *HTTPServer) PostDigitalWritePin(w http.ResponseWriter, r *http.Request) {
	var req DigitalWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := circuit.GetDigitalWritePin(*req.PinNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetDigitalState(req.Value)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(DigitalWritePinResponse{PinNumber: req.PinNumber, Value: &req.Value})
}

func (s *HTTPServer) PostAnalogWritePin(w http.ResponseWriter, r *http.Request) {
	var req AnalogWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := circuit.GetAnalogWritePin(req.PinNumber)
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
