package main

import (
	"encoding/json"
	"log"
	"time"

	"go.bug.st/serial"
)

type PortServer struct {
	Port         serial.Port
	AvaiblePorts []string
	Settings     serial.Mode
}

type PortServerResult[T any] struct {
	Data  T
	Error error
}

type WatchPortsResult = PortServerResult[[]string]

type ListenToPortResult = PortServerResult[string]

type Command string

const (
	SetupPinCommandType        Command = "setup_pin"
	DigitalWritePinCommandType Command = "digital_write_pin"
)

type SetupPinCommand struct {
	Command string `json:"command"`
	Pin     int    `json:"pin"`
	Mode    int    `json:"mode"`
}

type DigitalWritePinCommand struct {
	Command string `json:"command"`
	Pin     int    `json:"pin"`
	Value   int    `json:"value"`
}

func NewPortServer() *PortServer {
	return &PortServer{
		Port:         nil,
		AvaiblePorts: []string{},
		Settings: serial.Mode{
			BaudRate: 9600,
		},
	}
}

func (ps *PortServer) PopulatePorts() error {
	ports, err := serial.GetPortsList()
	if err != nil {
		return err
	}

	ps.AvaiblePorts = ports

	return nil
}

func (ps *PortServer) WatchPorts() chan WatchPortsResult {
	ch := make(chan WatchPortsResult)

	go func() {
		for {
			ports, err := serial.GetPortsList()
			if err != nil {
				ch <- WatchPortsResult{Error: err}
				continue
			}

			if len(ports) != len(ps.AvaiblePorts) {
				ps.AvaiblePorts = ports
				ch <- WatchPortsResult{Data: ports}
			}

			ch <- WatchPortsResult{Data: nil}
		}
	}()

	return ch
}

func (ps *PortServer) OpenPort(port string) error {
	var err error

	err = ps.PortExists(port)
	if err != nil {
		return err
	}

	if ps.Port != nil {
		err = ps.ClosePort()
		if err != nil {
			return err
		}
	}

	err = ps.attemptOpenPort(10, port)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) attemptOpenPort(attempts int, port string) error {
	var err error

	for i := 0; i < attempts; i++ {
		ps.Port, err = serial.Open(port, &ps.Settings)
		if err != nil {
			if i == attempts-1 {
				return &PortOpenTimeoutError{}
			}

			time.Sleep(100 * time.Millisecond)
			continue
		}

		break
	}

	return nil
}

func (ps *PortServer) ClosePort() error {
	if ps.Port == nil {
		return nil
	}

	err := ps.Port.Close()
	if err != nil {
		return err
	}

	ps.Port = nil

	return nil
}

func (ps *PortServer) SetupPin(pin int, mode PinMode) error {
	command := SetupPinCommand{
		Command: string(SetupPinCommandType),
		Pin:     pin,
		Mode:    int(mode),
	}

	json, err := json.Marshal(command)
	if err != nil {
		return err
	}

	log.Print("Setup pin command to ps.Write: ", string(json))
	err = ps.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) WriteDigitalPin(pin int, value PinState) error {
	command := DigitalWritePinCommand{
		Command: string(DigitalWritePinCommandType),
		Pin:     pin,
		Value:   int(value),
	}

	json, err := json.Marshal(command)
	if err != nil {
		return err
	}

	log.Print("Write digital pin command to ps.Write: ", string(json))
	err = ps.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) Write(data []byte) error {
	if ps.Port == nil {
		return &PortNotOpenError{}
	}

	commandWithNewline := append(data, '\n')

	log.Print("Writing to port: ", string(commandWithNewline))

	_, err := ps.Port.Write([]byte(commandWithNewline))
	if err != nil {
		return err
	}

	err = ps.Port.ResetOutputBuffer()
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) ListenToPort() chan ListenToPortResult {
	ch := make(chan ListenToPortResult)
	buf := make([]byte, 100)
	var input string

	go func() {
		for {
			if ps.Port == nil {
				continue
			}

			n, err := ps.Port.Read(buf)
			if err != nil {
				ch <- ListenToPortResult{Error: err}
				continue
			}

			input += string(buf[:n])

			log.Print(input)

			if input[len(input)-1] == '\n' {
				ch <- ListenToPortResult{Data: input}
				log.Print(input)
				input = ""
			}
		}
	}()

	return ch
}

func (ps *PortServer) PortExists(port string) error {
	for _, p := range ps.AvaiblePorts {
		if p == port {
			return nil
		}
	}

	return &PortDoesNotExistError{}
}

type PortDoesNotExistError struct{}

func (e *PortDoesNotExistError) Error() string {
	return "Port does not exist"
}

type PortNotOpenError struct{}

func (e *PortNotOpenError) Error() string {
	return "Port is not open"
}

type PortOpenTimeoutError struct{}

func (e *PortOpenTimeoutError) Error() string {
	return "Port open timeout"
}

type InvalidPinModeError struct{}

func (e *InvalidPinModeError) Error() string {
	return "Invalid pin mode"
}
