package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	thttp "tobro/api/http"
	"tobro/api/ws"
	"tobro/db"
	"tobro/internal/models"
	"tobro/internal/models/circuit"
	"tobro/internal/models/sketch"
	"tobro/pkg/arduino"
	"tobro/pkg/monitor"
	"tobro/pkg/store"
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
	st := store.New(db.NewClient())
	if err := st.Connect(); err != nil {
		log.Fatal(err)
	}
	defer st.Disconnect()

	a := arduino.NewServer()

	c := circuit.New(0, "Default Circuit", models.ArduinoNano, a)
	dbCircuit, err := st.InitCircuit(c)
	if err != nil {
		log.Fatal(err)
	}

	c.UpdateFromDBModel(dbCircuit)

	sk := sketch.New(0, "Default Sketch", c)

	hs := thttp.NewHTTPServer(st, c, sk)

	hub := ws.NewWSHub()
	go hub.Run()

	m := monitor.New(hub, a, c)
	m.Run()

	r := mux.NewRouter()

	r.Use(thttp.EnableCORS)
	r.Use(thttp.LogRequest)

	h := thttp.HandlerFromMux(hs, r)
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, a, w, r)
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
