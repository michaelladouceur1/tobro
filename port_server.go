package main

import (
	"encoding/json"
	"log"
	"time"

	"go.bug.st/serial"
)

type PortServer struct {
	PortName     string
	Port         serial.Port
	AvaiblePorts chan []string
	Settings     serial.Mode
}

type PortServerResult[T any] struct {
	Data  T
	Error error
}

type WatchPortsResult = PortServerResult[[]string]

type ListenToPortResult = PortServerResult[string]

type Command int

const (
	SetupPinCommandType        Command = 1
	DigitalWritePinCommandType Command = 2
	AnalogWritePinCommandType  Command = 3
)

// cmd: 1 (setup_pin)
type SetupPinCommand struct {
	Command uint `json:"c"`
	Pin     uint `json:"p"`
	Mode    uint `json:"m"`
}

// cmd: 2 (digital_write_pin)
type DigitalWritePinCommand struct {
	Command uint `json:"c"`
	Pin     uint `json:"p"`
	Value   uint `json:"v"`
}

// cmd: 3 (analog_write_pin)
type AnalogWritePinCommand struct {
	Command uint `json:"c"`
	Pin     uint `json:"p"`
	Value   uint `json:"v"`
}

func NewPortServer() *PortServer {
	return &PortServer{
		Port:         nil,
		AvaiblePorts: make(chan []string),
		Settings: serial.Mode{
			BaudRate: 115200,
		},
	}
}

func (ps *PortServer) WatchPorts() {
	for {
		ports, err := serial.GetPortsList()
		if err != nil {
			continue
		}

		ps.AvaiblePorts <- ports

		// if ps.Port != nil && ps.PortName != "" {
		// 	if err = ps.PortExists(ps.PortName); err != nil {
		// 		log.Print("Port does not exist")
		// 		ps.ClosePort()
		// 	}

		// 	continue
		// }

		time.Sleep(1 * time.Second)
	}
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

	ps.PortName = port

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
	ps.PortName = ""

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

func (ps *PortServer) SetupPin(pin int, mode PinMode) error {
	command := SetupPinCommand{
		Command: uint(SetupPinCommandType),
		Pin:     uint(pin),
		Mode:    uint(mode),
	}

	json, err := json.Marshal(command)
	if err != nil {
		return err
	}

	err = ps.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) WriteDigitalPin(pin int, value int) error {
	command := DigitalWritePinCommand{
		Command: uint(DigitalWritePinCommandType),
		Pin:     uint(pin),
		Value:   uint(value),
	}

	json, err := json.Marshal(command)
	if err != nil {
		return err
	}

	err = ps.Write(json)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) WriteAnalogPin(pin int, value int) error {
	command := AnalogWritePinCommand{
		Command: uint(AnalogWritePinCommandType),
		Pin:     uint(pin),
		Value:   uint(value),
	}

	json, err := json.Marshal(command)
	if err != nil {
		return err
	}

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

	// bits := len(commandWithNewline) * 8
	// log.Printf("Bits: %d", bits)

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
	availablePorts := <-ps.AvaiblePorts
	log.Print(availablePorts)
	for _, p := range availablePorts {
		log.Print(p)
		if p == port {
			log.Print("Port exists")
			return nil
		}
		log.Print("Port does not exist")
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
