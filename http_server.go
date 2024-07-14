package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPServer struct{}

func NewHTTPServer() *HTTPServer {
	return &HTTPServer{}
}

func enableCORS(e *gin.Engine) {
	e.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

}

// (GET /ping)
func (s *HTTPServer) GetPing(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// (POST /connect)
func (s *HTTPServer) PostConnect(c *gin.Context) {
	var req ConnectRequest
	log.Print("PostConnect", req.Port)
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	err := portServer.OpenPort(req.Port)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	c.JSON(http.StatusOK, ConnectResponse{Port: &req.Port})
}

// (POST /delay)
func (s *HTTPServer) PostDelay(c *gin.Context) {
	var req PostDelayJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	command := Command{
		Command: "delay",
		Delay:   req.Delay,
	}

	json, err := json.Marshal(command)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: &err})
		return
	}

	commandWithNewline := append(json, '\n')

	err = portServer.Write(commandWithNewline)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: &err})
		return
	}

	c.JSON(http.StatusOK, command)
}

// (POST /setup_pin)
func (s *HTTPServer) PostSetupPin(c *gin.Context) {
	var req SetupPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	c.JSON(http.StatusOK, SetupPinResponse{Mode: &req.Mode, Pin: &req.Pin})
}

// (POST /write_pin)
func (s *HTTPServer) PostWritePin(c *gin.Context) {
	var req WritePinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	c.JSON(http.StatusOK, WritePinResponse{Pin: &req.Pin, Value: &req.Value})
}
