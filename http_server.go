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

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
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
	log.Print("PostConnect")
	var req ConnectRequest
	if err := decodeRequestBody(w, r, req); err != nil {
		return
	}

	log.Print("PostConnect", req.Port)

	err := portServer.OpenPort(req.Port)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	portServer.ListenToPort()

	json.NewEncoder(w).Encode(ConnectResponse{Port: &req.Port})
}

func (s *HTTPServer) PostSetupPin(w http.ResponseWriter, r *http.Request) {
	log.Print("PostSetupPin")
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
	log.Print("PostDigitalWritePin")
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
	log.Print("PostAnalogWritePin")
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
