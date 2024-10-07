package monitor

import (
	"encoding/json"
	"log"
	"time"

	"tobro/api/ws"
	"tobro/internal/models/circuit"
	"tobro/internal/models/pin"
	"tobro/pkg/arduino"
)

type Monitor struct {
	hub     *ws.WSHub
	ps      *arduino.PortServer
	circuit *circuit.Circuit
}

func New(hub *ws.WSHub, ps *arduino.PortServer, circuit *circuit.Circuit) *Monitor {
	return &Monitor{
		hub:     hub,
		ps:      ps,
		circuit: circuit,
	}
}

func (m *Monitor) Run() {
	go m.watchPorts()
	go m.watchPortConnection()
	go m.watchPinState()
}

func (m *Monitor) watchPorts() {
	for ports := range m.ps.AvaiblePorts {
		portNames := make([]string, 0, len(ports))
		for _, port := range ports {
			portNames = append(portNames, port.Name)
		}
		json, err := json.Marshal(ws.CreatePortsResponse(portNames))
		if err != nil {
			log.Fatal(err)
			return
		}

		m.hub.Broadcast <- json
	}
}

func (m *Monitor) watchPortConnection() {
	for {
		connected := <-m.ps.Connected
		json, err := json.Marshal(ws.CreatePortConnectionResponse(connected, m.ps.PortName))
		if err != nil {
			log.Fatal(err)
			return
		}

		m.hub.Broadcast <- json
	}
}

func (m *Monitor) watchPinState() {
	for {
		if m.circuit == nil {
			log.Print("Board is nil")
			continue
		}

		for _, p := range m.circuit.Pins {
			go func(p pin.Pin) {
				select {
				case state, ok := <-p.State:
					if !ok {
						log.Printf("Pin %d state channel is closed", p.PinNumber)
						return
					}

					json, err := json.Marshal(ws.CreatePinStateResponse(p.PinNumber, state))
					if err != nil {
						log.Fatal(err)
						return
					}

					m.hub.Broadcast <- json
				}
			}(p)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
