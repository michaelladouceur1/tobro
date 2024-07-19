package main

import (
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

	portServer.ListenToPort()

	c.JSON(http.StatusOK, ConnectResponse{Port: &req.Port})
}

// (POST /setup_pin)
func (s *HTTPServer) PostSetupPin(c *gin.Context) {
	var req SetupPinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	pin, err := board.GetPin(req.Pin)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	err = pin.SetPinMode(req.Mode)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: &err})
		return
	}

	modeStr := string(req.Mode)
	c.JSON(http.StatusOK, SetupPinResponse{Mode: &modeStr, Pin: &req.Pin})
}

// (POST /digital_write_pin)
func (s *HTTPServer) PostDigitalWritePin(c *gin.Context) {
	var req WritePinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	pin, err := board.GetPin(req.Pin)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	err = pin.SetDigitalPinState(req.Value)

	if err != nil {
		err := err.Error()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: &err})
		return
	}

	c.JSON(http.StatusOK, DigitalWritePinResponse{Pin: &req.Pin, Value: &req.Value})
}

// (POST /pwm)
func (s *HTTPServer) PostPwm(c *gin.Context) {
	var req PWMRequest
	// json, err := json.Marshal(req)
	// log.Println("PostPwm", string(json))
	if err := c.ShouldBindJSON(&req); err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	log.Println("PWM", req)

	pin, err := board.GetPin(req.Pin)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusBadRequest, ErrorResponse{Message: &err})
		return
	}

	log.Println("PWM", pin)

	err = pin.PWM(req.DutyCycle, req.Period, 5000)
	if err != nil {
		err := err.Error()
		c.JSON(http.StatusInternalServerError, ErrorResponse{Message: &err})
		return
	}

	value := 1
	c.JSON(http.StatusOK, DigitalWritePinResponse{Pin: &req.Pin, Value: &value})

	// c.JSON(http.StatusOK, PWM{Pin: &req.Pin, DutyCycle: &req.DutyCycle, Period: &req.Period})
}
