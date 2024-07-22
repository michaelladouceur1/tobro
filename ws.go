package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

// func createSuccessResponse(message string) []byte {
// 	res := SuccessResponse{
// 		Message: message,
// 	}

// 	json, err := json.Marshal(res)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return json
// }

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
			DigitalPins: board.DigitalPins,
			AnalogPins:  board.AnalogPins,
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

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	// sendBoardState(c)
	go watchPorts(c)

	for {
		// _, _, err := c.ReadMessage()
		// if err != nil {
		// 	log.Print("read:", err)
		// 	break
		// }
		time.Sleep(1 * time.Second)
	}
}

func sendBoardState(c *websocket.Conn) {
	boardState := board.GetState()
	json, err := json.Marshal(createBoardResponse(boardState))
	if err != nil {
		log.Fatal(err)
	}

	err = c.WriteMessage(websocket.TextMessage, json)
	if err != nil {
		log.Fatal(err)
	}
}

func watchPorts(c *websocket.Conn) {
	go portServer.WatchPorts()

	for ports := range portServer.AvaiblePorts {
		json, err := json.Marshal(createPortsResponse(ports))
		if err != nil {
			log.Fatal(err)
			return
		}

		err = c.WriteMessage(websocket.TextMessage, json)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}
