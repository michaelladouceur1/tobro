package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"tobro/pkg/arduino"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

type BaseResponse[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type PortsResponseData struct {
	Ports []string `json:"ports"`
}

type PortConnectionResponseData struct {
	Connected bool   `json:"connected"`
	PortName  string `json:"portName"`
}

type PinStateResponseData struct {
	PinNumber int `json:"pinNumber"`
	State     int `json:"state"`
}

func CreatePortsResponse(ports []string) BaseResponse[PortsResponseData] {
	return BaseResponse[PortsResponseData]{
		Type: "ports",
		Data: PortsResponseData{
			Ports: ports,
		},
	}
}

func CreatePortConnectionResponse(connected bool, portName string) BaseResponse[PortConnectionResponseData] {
	return BaseResponse[PortConnectionResponseData]{
		Type: "port_connection",
		Data: PortConnectionResponseData{
			Connected: connected,
			PortName:  portName,
		},
	}
}

func CreatePinStateResponse(pinNumber int, state int) BaseResponse[PinStateResponseData] {
	return BaseResponse[PinStateResponseData]{
		Type: "pin_state",
		Data: PinStateResponseData{
			PinNumber: pinNumber,
			State:     state,
		},
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		allowedOrigins := map[string]bool{
			"http://localhost:3000": true,
			"http://localhost:8000": true,
			"http://localhost:8080": true,
		}
		origin := r.Header.Get("Origin")
		return allowedOrigins[origin]
	},
}

func ServeWs(hub *WSHub, ps *arduino.PortServer, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	client := &WSClient{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	SendPortName(hub, ps)

	go client.Write()
	go client.Read()
}

func SendPortName(hub *WSHub, ps *arduino.PortServer) error {
	connected := ps.PortName != ""
	json, err := json.Marshal(CreatePortConnectionResponse(connected, ps.PortName))
	if err != nil {
		return err
	}

	hub.Broadcast <- json

	return nil
}
