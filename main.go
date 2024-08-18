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
	s, err := NewSession("session.json")
	if err != nil {
		log.Fatal(err)
	}

	ps := NewPortServer(s)
	log.Print("PortServer created")

	dal := NewDAL()
	if err := dal.Connect(); err != nil {
		log.Fatal(err)
	}
	defer dal.Disconnect()

	c := NewCircuit(0, "Default Circuit", ArduinoNano, ps)
	dbCircuit, err := dal.InitCircuit(c)
	if err != nil {
		log.Fatal(err)
	}

	c.UpdateFromDBModel(dbCircuit)

	hs := NewHTTPServer(s, dal, c)

	hub := NewWSHub()
	go hub.Run()

	m := NewMonitor(hub, ps, c)
	m.Run()

	r := mux.NewRouter()

	r.Use(enableCORS)
	r.Use(logRequest)

	h := HandlerFromMux(hs, r)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, ps, w, r)
	})

	if os.Getenv("GO_ENV") != "dev" {
		log.Print("Serving UI from embed.FS")
		r.PathPrefix("/").Handler(http.FileServer(http.FS(uiFS)))
	} else {
		log.Print("Serving UI from web/build")
		r.PathPrefix("/").Handler(http.FileServer(http.Dir("web/build")))
	}

	srv := &http.Server{
		Handler: h,
		Addr:    ":8080",
	}

	srv.ListenAndServe()
}
