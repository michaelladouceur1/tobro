package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
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
var httpServer *HTTPServer

func main() {
	portServer = NewPortServer()
	httpServer = NewHTTPServer()

	r := gin.Default()

	enableCORS(r)

	RegisterHandlers(r, httpServer)

	s := &http.Server{
		Handler: r,
		Addr:    ":8081",
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("HTTP server started")

	mux := http.NewServeMux()
	mux.HandleFunc("/", staticHandler)
	mux.HandleFunc("/ws", websocketHandler)

	go func() {
		if err := http.ListenAndServe(":8080", mux); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Print("Websocket server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Print("Shutting down servers...")
}
