package http_server

import (
	"encoding/json"
	"net/http"
	"tobro/internal/models/circuit"
	"tobro/internal/models/sketch"
	"tobro/pkg/store"
)

type HTTPServer struct {
	store   *store.Store
	circuit *circuit.Circuit
	sketch  *sketch.Sketch
}

func NewHTTPServer(store *store.Store, circuit *circuit.Circuit, sketch *sketch.Sketch) *HTTPServer {
	return &HTTPServer{
		store:   store,
		circuit: circuit,
		sketch:  sketch,
	}
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

	json.NewEncoder(w).Encode(ConnectResponse{Port: &req.Port})
}

func (s *HTTPServer) GetBoards(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(BoardsResponse{Boards: circuit.GetSupportedBoards()})
}
