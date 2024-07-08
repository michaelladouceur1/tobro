package main

// import (
// 	"log"
// 	"net/http"

// 	"tobro-server/api"

// 	"github.com/gin-gonic/gin"
// )

// func main() {
// 	// create a type that satisfies the `api.ServerInterface`, which contains an implementation of every operation from the generated code
// 	server := api.NewServer()

// 	r := gin.Default()

// 	api.RegisterHandlers(r, server)

// 	// And we serve HTTP until the world ends.

// 	s := &http.Server{
// 		Handler: r,
// 		Addr:    "0.0.0.0:8080",
// 	}

// 	// And we serve HTTP until the world ends.
// 	log.Fatal(s.ListenAndServe())
// }

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"go.bug.st/serial"
)

type Input struct {
	Data string `json:"data"`
}

type Command struct {
	Command string `json:"command"`
	Delay   int    `json:"delay"`
}

func (i *Input) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Data string `json:"data"`
	}{
		Data: "",
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	i.Data = aux.Data
	return nil
}

func main() {
	ports, err := serial.GetPortsList()
	if err != nil {
		panic(err)
	}

	if len(ports) == 0 {
		fmt.Println("No ports found")
		return
	}

	for _, port := range ports {
		fmt.Println(port)
		fmt.Println()
	}

	mode := &serial.Mode{
		BaudRate: 9600,
	}

	port, err := serial.Open(ports[0], mode)
	if err != nil {
		panic(err)
	}
	// defer port.Close()

	go listenAsync(port)
	waitForKey(port)
}

func waitForKey(port serial.Port) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		if s == "exit" {
			return
		}

		// command := "{\"command\":\"" + s + "\"}"
		sInt, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
			return
		}

		command := Command{
			Command: "delay",
			Delay:   sInt,
		}

		json, err := json.Marshal(command)
		if err != nil {
			log.Fatal(err)
			return
		}

		// fmt.Print("Command: ")
		// fmt.Println(string(json))

		// sizeInBytes := len(command)
		// fmt.Printf("Size in bytes: %d\n", sizeInBytes)

		// commandByteArr := make([]byte, len(command))
		// copy(commandByteArr[:], command)

		commandWithNewline := append(json, '\n')

		_, err = port.Write([]byte(commandWithNewline))
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func listenAsync(port serial.Port) {
	var input string

	buff := make([]byte, 1000)
	for {
		// Reads up to 100 bytes
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}

		input += string(buff[:n])

		if strings.Contains(input, "\n") {
			fmt.Printf("Input: %s\n", input)
			input = ""
		}

		// If we receive a newline stop reading
		// if strings.Contains(string(buff[:n]), "\n") {
		// 	fmt.Println("\nNewline detected")
		// 	break
		// }
	}
}
