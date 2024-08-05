package main

import (
	"encoding/json"
	"net/http"
	"tobro/db"
)

type HTTPServer struct{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
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

func (s *HTTPServer) PostCreateCircuit(w http.ResponseWriter, r *http.Request) {
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

	board = NewBoard(boardType, portServer)

	createdBoard, err := dbClient.Circuit.CreateOne(db.Circuit.Name.Equals(req.Name), db.Circuit.Board.Set(string(boardType))).Exec(dbCtx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, pin := range board.Pins {
		pinModeInt := int(pin.Mode)
		dbClient.Pin.CreateOne(db.Pin.Pin.Set(pin.ID), db.Pin.Circuit.Link(db.Circuit.ID.Equals(createdBoard.ID)), db.Pin.Mode.Set(pinModeInt)).Exec(dbCtx)
	}

	json.NewEncoder(w).Encode(createBoardResponse(*board))
}

func (s *HTTPServer) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	var req SetupPinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := board.GetPin(req.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = pin.SetMode(req.Mode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	modeStr := string(req.Mode)
	json.NewEncoder(w).Encode(SetupPinResponse{Mode: &modeStr, Pin: &req.Pin})
}

func (s *HTTPServer) PostDigitalWritePin(w http.ResponseWriter, r *http.Request) {
	var req DigitalWritePinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pin, err := board.GetDigitalWritePin(req.Pin)
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

	pin, err := board.GetAnalogWritePin(req.Pin)
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
