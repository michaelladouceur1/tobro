package main

import (
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

func (ps *PortServer) Write(data []byte) error {
	if ps.Port == nil {
		return &PortNotOpenError{}
	}

	_, err := ps.Port.Write([]byte(data))
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

			if input[len(input)-1] == '\n' {
				ch <- ListenToPortResult{Data: input}
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
