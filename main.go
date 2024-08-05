package main

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"tobro/db"

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
var dbCtx context.Context
var dbClient *db.PrismaClient
var board *Board

func main() {
	portServer = NewPortServer()
	httpServer = NewHTTPServer()

	dbCtx = context.Background()
	dbClient = db.NewClient()
	if err := dbClient.Connect(); err != nil {
		log.Fatal(err)
	}
	defer dbClient.Disconnect()

	hub := NewWSHub()
	go hub.Run()

	board = NewBoard(ArduinoNano, portServer)

	monitor := NewMonitor(hub, portServer, board)
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
