package monitor

import (
	"encoding/json"
	"log"
	"time"

	"tobro/pkg/arduino"
	"tobro/pkg/models"
	"tobro/pkg/ws"
)

type Monitor struct {
	hub     *ws.WSHub
	ps      *arduino.PortServer
	circuit *models.Circuit
}

func NewMonitor(hub *ws.WSHub, ps *arduino.PortServer, circuit *models.Circuit) *Monitor {
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
		json, err := json.Marshal(ws.CreatePortsResponse(ports))
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

		for _, pin := range m.circuit.Pins {
			go func(pin models.Pin) {
				select {
				case state, ok := <-pin.State:
					if !ok {
						log.Printf("Pin %d state channel is closed", pin.PinNumber)
						return
					}

					json, err := json.Marshal(ws.CreatePinStateResponse(pin.PinNumber, state))
					if err != nil {
						log.Fatal(err)
						return
					}

					m.hub.Broadcast <- json
				}
			}(pin)
		}

		time.Sleep(10 * time.Millisecond)
	}
}
