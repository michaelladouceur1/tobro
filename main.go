package main

import (
	"bufio"
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"go.bug.st/serial"
)

type Input struct {
	Delay int `json:"delay"`
}

type Command struct {
	Command string `json:"command"`
	Delay   int    `json:"delay"`
}

//go:embed web/build
var UI embed.FS

var uiFS fs.FS

func init() {
	var err error
	uiFS, err = fs.Sub(UI, "web/build")
	if err != nil {
		panic(err)
	}
}

// var port serial.Port
var portServer *PortServer

func main() {
	portServer = NewPortServer()

	// portServer.PopulatePorts()
	// portsChan := portServer.WatchPorts()

	// for result := range portsChan {
	// 	if result.Error != nil {
	// 		log.Fatal(result.Error)
	// 		return
	// 	}

	// 	if result.Data != nil {
	// 		log.Println("Avaible ports:", result.Data)
	// 		log.Println("Opening port:", result.Data[0])
	// 		err := portServer.OpenPort(portServer.AvaiblePorts[0])
	// 		if err != nil {
	// 			log.Fatal(err)
	// 			return
	// 		}
	// 		portServer.ListenToPort()
	// 		break
	// 	}

	// }

	// portServer.OpenPort(portServer.AvaiblePorts[0])
	// portServer.ListenToPort()

	// go listenAsync(port)
	go waitForKey(portServer.Port)

	mux := http.NewServeMux()
	mux.HandleFunc("/", staticHandler)
	mux.HandleFunc("/ws", websocketHandler)
	mux.HandleFunc("/connect", connectHandler)
	mux.HandleFunc("/delay", delayHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
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
