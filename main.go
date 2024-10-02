package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"tobro/pkg/arduino"
	tobroHTTP "tobro/pkg/http"
	"tobro/pkg/models"
	"tobro/pkg/models/circuit"
	"tobro/pkg/models/sketch"
	"tobro/pkg/monitor"
	"tobro/pkg/store"
	"tobro/pkg/ws"
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
	ps := arduino.NewServer()
	log.Print("PortServer created")

	st := store.New()
	if err := st.Connect(); err != nil {
		log.Fatal(err)
	}
	defer st.Disconnect()

	c := circuit.New(0, "Default Circuit", models.ArduinoNano, ps)
	dbCircuit, err := st.InitCircuit(c)
	if err != nil {
		log.Fatal(err)
	}

	c.UpdateFromDBModel(dbCircuit)

	sk := sketch.New(0, "Default Sketch", c)

	hs := tobroHTTP.NewHTTPServer(st, c, sk)

	hub := ws.NewWSHub()
	go hub.Run()

	m := monitor.New(hub, ps, c)
	m.Run()

	r := mux.NewRouter()

	r.Use(tobroHTTP.EnableCORS)
	r.Use(tobroHTTP.LogRequest)

	h := tobroHTTP.HandlerFromMux(hs, r)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, ps, w, r)
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
