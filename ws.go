package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
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

func createBoardResponse(board Board) BaseResponse[Board] {
	return BaseResponse[Board]{
		Type: "board",
		Data: Board{
			Pins: board.Pins,
		},
	}
}

func createPortsResponse(ports []string) BaseResponse[PortsResponseData] {
	return BaseResponse[PortsResponseData]{
		Type: "ports",
		Data: PortsResponseData{
			Ports: ports,
		},
	}
}

var upgrader = websocket.Upgrader{}

func serveWs(hub *WSHub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}

	client := &WSClient{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	sendBoardState(client)

	go client.Write()
	go client.Read()
}

func sendBoardState(c *WSClient) {
	json, err := json.Marshal(createBoardResponse(*board))
	if err != nil {
		log.Fatal(err)
	}

	c.send <- json
}
