package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

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

func main() {
	portServer := NewPortServer()

	dal := NewDAL()
	if err := dal.Connect(); err != nil {
		log.Fatal(err)
	}
	defer dal.Disconnect()

	circuit := NewCircuit(0, "Default Circuit", ArduinoNano, portServer)
	dbCircuit, err := dal.InitCircuit(circuit)
	if err != nil {
		log.Fatal(err)
	}

	circuit.UpdateFromDBModel(dbCircuit)

	httpServer := NewHTTPServer(circuit, dal, portServer)

	hub := NewWSHub()
	go hub.Run()

	monitor := NewMonitor(hub, portServer, circuit)
	go monitor.Run()

	route := mux.NewRouter()

	route.Use(enableCORS)
	route.Use(logRequest)

	h := HandlerFromMux(httpServer, route)
	route.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if os.Getenv("GO_ENV") != "dev" {
		log.Print("Serving UI from embed.FS")
		route.PathPrefix("/").Handler(http.FileServer(http.FS(uiFS)))
	} else {
		log.Print("Serving UI from web/build")
		route.PathPrefix("/").Handler(http.FileServer(http.Dir("web/build")))
	}

	s := &http.Server{
		Handler: h,
		Addr:    ":8080",
	}

	s.ListenAndServe()
}
