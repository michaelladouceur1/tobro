package main

import (
	"encoding/json"
	"log"
	"time"

	"go.bug.st/serial"
)

type PortServer struct {
	session      *Session
	Port         serial.Port
	Connected    chan bool
	PortName     string
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

func NewPortServer(session *Session) *PortServer {
	ps := &PortServer{
		session:      session,
		Port:         nil,
		Connected:    make(chan bool),
		AvaiblePorts: make(chan []string),
		Settings: serial.Mode{
			BaudRate: 115200,
		},
	}

	go ps.watchPorts()
	go ps.autoConnect()

	return ps
}

func (ps *PortServer) OpenPort(port string) error {
	var err error

	err = ps.portExists(port)
	if err != nil {
		log.Print(err)
		return err
	}

	if ps.Port != nil {
		err = ps.closePort()
		if err != nil {
			log.Print(err)
			return err
		}
	}

	err = ps.attemptOpenPort(10, port)
	if err != nil {
		log.Print(err)
		return err
	}

	ps.PortName = port
	ps.Connected <- true

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

	err = ps.write(json)
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

	err = ps.write(json)
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

	err = ps.write(json)
	if err != nil {
		return err
	}

	return nil
}

func (ps *PortServer) write(data []byte) error {
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

func (ps *PortServer) watchPorts() {
	for {
		ports, err := serial.GetPortsList()
		if err != nil {
			continue
		}

		ps.AvaiblePorts <- ports

		time.Sleep(1 * time.Second)
	}
}

// this function assumes the arduino is connected to the first available port and the only connected serial device
// this is a temporary solution. a symlink to the arduino port should be created in the future
func (ps *PortServer) autoConnect() {
	for {
		ports := <-ps.AvaiblePorts
		if len(ports) == 0 {
			log.Print("No ports available")
			ps.resetPort()
			continue
		}

		if len(ports) > 0 && ps.Port == nil {
			log.Printf("Attempting to open port: %s", ports[0])
			err := ps.OpenPort(ports[0])
			if err != nil {
				log.Print(err)
			}
		}
	}
}

func (ps *PortServer) portExists(port string) error {
	availablePorts := <-ps.AvaiblePorts
	for _, p := range availablePorts {
		log.Print(p)
		if p == port {
			return nil
		}
	}

	return &PortDoesNotExistError{}
}

func (ps *PortServer) closePort() error {
	if ps.Port == nil {
		return nil
	}

	err := ps.Port.Close()
	if err != nil {
		return err
	}

	ps.resetPort()

	return nil
}

func (ps *PortServer) resetPort() {
	ps.Port = nil
	ps.PortName = ""
	ps.Connected <- false
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
