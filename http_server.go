package main

import (
	"encoding/json"
	"net/http"
)

type HTTPServer struct{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func decodeRequestBody(w http.ResponseWriter, r *http.Request, target interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(target); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}

func (s *HTTPServer) GetPing(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(SuccessResponse{Message: "pong"})
}

func (s *HTTPServer) PostConnect(w http.ResponseWriter, r *http.Request) {
	var req ConnectRequest
	if err := decodeRequestBody(w, r, req); err != nil {
		return
	}

	err := portServer.OpenPort(req.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	portServer.ListenToPort()

	json.NewEncoder(w).Encode(ConnectResponse{Port: &req.Port})
}

func (s *HTTPServer) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	var req SetupPinRequest
	if err := decodeRequestBody(w, r, req); err != nil {
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
	if err := decodeRequestBody(w, r, req); err != nil {
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
	if err := decodeRequestBody(w, r, req); err != nil {
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
