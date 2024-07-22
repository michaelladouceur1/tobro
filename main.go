package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
)

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
var httpServer *HTTPServer
var board *Board

func main() {
	portServer = NewPortServer()
	httpServer = NewHTTPServer()

	board = NewBoard(ArduinoNano, portServer)

	route := mux.NewRouter()

	route.Use(enableCORS)

	h := HandlerFromMux(httpServer, route)
	route.HandleFunc("/ws", websocketHandler)
	route.PathPrefix("/").Handler(http.FileServer(http.FS(uiFS)))

	s := &http.Server{
		Handler: h,
		Addr:    ":8080",
	}

	s.ListenAndServe()
}
