package main

import (
	"encoding/json"
	"log"
	"time"
)

type Monitor struct {
	hub   *WSHub
	ps    *PortServer
	board *Board
}

func NewMonitor(hub *WSHub, ps *PortServer, board *Board) *Monitor {
	return &Monitor{
		hub:   hub,
		ps:    ps,
		board: board,
	}
}

func (m *Monitor) Run() {
	go m.watchPorts()
	go m.watchPinState()
}

func (m *Monitor) watchPorts() {
	go m.ps.WatchPorts()

	for ports := range m.ps.AvaiblePorts {
		json, err := json.Marshal(createPortsResponse(ports))
		if err != nil {
			log.Fatal(err)
			return
		}

		m.hub.broadcast <- json
	}
}

func (m *Monitor) watchPinState() {
	for {
		if m.board == nil {
			log.Print("Board is nil")
			continue
		}

		for _, pin := range m.board.Pins {
			go func(pin Pin) {
				select {
				case state, ok := <-pin.State:
					if !ok {
						log.Printf("Pin %d state channel is closed", pin.ID)
						return
					}

					json, err := json.Marshal(createPinStateResponse(pin.ID, state))
					if err != nil {
						log.Fatal(err)
						return
					}

					log.Print("Pin state: ", string(json))

					m.hub.broadcast <- json
				}
			}(pin)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
