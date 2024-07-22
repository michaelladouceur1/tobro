package main

import (
	"encoding/json"
	"log"
)

type Monitor struct {
	hub *WSHub
	ps  *PortServer
}

func NewMonitor(hub *WSHub, ps *PortServer) *Monitor {
	return &Monitor{
		hub: hub,
		ps:  ps,
	}
}

func (m *Monitor) Run() {
	go m.watchPorts()
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
