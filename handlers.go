package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/websocket"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

func createSuccessResponse(message string) []byte {
	res := SuccessResponse{
		Message: message,
	}

	json, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}

	return json
}

type BaseResponse[T any] struct {
	Type string `json:"type"`
	Data T      `json:"data"`
}

type PortsResponseData struct {
	Ports []string `json:"ports"`
}

func createPortsResponse(ports []string) BaseResponse[PortsResponseData] {
	return BaseResponse[PortsResponseData]{
		Type: "ports",
		Data: PortsResponseData{
			Ports: ports,
		},
	}
}

var upgrader = websocket.Upgrader{}

func staticHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)
	if path == "/" {
		path = "/index.html"
	}
	path = strings.TrimPrefix(path, "/")

	file, err := uiFS.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Println("file", path, "not found:", err)
			http.NotFound(w, r)
			return
		}
		log.Println("file", path, "cannot be read:", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	contentType := mime.TypeByExtension(filepath.Ext(path))
	w.Header().Set("Content-Type", contentType)
	if strings.HasPrefix(path, "static/") {
		w.Header().Set("Cache-Control", "public, max-age=31536000")
	}
	stat, err := file.Stat()
	if err == nil && stat.Size() > 0 {
		w.Header().Set("Content-Length", fmt.Sprintf("%d", stat.Size()))
	}

	n, _ := io.Copy(w, file)
	log.Println("file", path, "copied", n, "bytes")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	portsChan := portServer.WatchPorts()

	for result := range portsChan {
		if result.Error != nil {
			log.Fatal(result.Error)
			return
		}

		if result.Data != nil {
			json, err := json.Marshal(createPortsResponse(result.Data))
			if err != nil {
				log.Fatal(err)
				return
			}

			// err = portServer.OpenPort(portServer.AvaiblePorts[0])
			// if err != nil {
			// 	log.Fatal(err)
			// 	return
			// }

			// portServer.ListenToPort()

			err = c.WriteMessage(websocket.TextMessage, json)
			if err != nil {
				log.Fatal(err)
				return
			}
		}

	}
}

type ConnectData struct {
	Port string `json:"port"`
}

func connectHandler(w http.ResponseWriter, r *http.Request) {
	var data ConnectData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("connectHandler: %v", data)

	err := portServer.OpenPort(data.Port)
	if err != nil {
		log.Fatal(err)
		return
	}

	portServer.ListenToPort()

	w.WriteHeader(http.StatusNoContent)
	w.Write(createSuccessResponse("connected"))
}

func delayHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Delay handler called")
	var input Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	command := Command{
		Command: "delay",
		Delay:   input.Delay,
	}

	json, err := json.Marshal(command)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	commandWithNewline := append(json, '\n')

	err = portServer.Write(commandWithNewline)
	if err != nil {
		log.Fatal("ERRORRRRRR: ", err)
		return
	}
}
