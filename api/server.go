package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// optional code omitted

type Server struct{}

func NewServer() Server {
	return Server{}
}

var axes []Axis

// // (GET /ping)
// func (Server) GetPing(ctx *gin.Context) {
// 	resp := Pong{
// 		Ping: "pong",
// 	}

// 	ctx.JSON(http.StatusOK, resp)
// }

// (GET /axis)
func (Server) GetAxis(ctx *gin.Context) {
	resp := Axis{
		Max:      10,
		MaxSpeed: 100,
		Min:      0,
		MinSpeed: 0,
		Type:     Angular,
	}

	ctx.JSON(http.StatusOK, resp)
}

// (POST /axis)
func (Server) PostAxis(ctx *gin.Context) {
	var req Axis
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	axes = append(axes, req)

	ctx.JSON(http.StatusCreated, req)
}

// (POST /rotate_axis)
func (Server) PostRotateAxis(ctx *gin.Context) {
	var req RotateAxisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	targetAxis := findAxisById(req.AxisId)

	rotateAxis(targetAxis, req.Angle)

	ctx.JSON(http.StatusOK, gin.H{})
}

func findAxisById(id *string) *Axis {
	for i := range axes {
		if axes[i].Id == id {
			return &axes[i]
		}
	}

	return nil
}

func rotateAxis(axis *Axis, angle float32) {
	if axis == nil {
		return
	}

	*axis.Value = angle
}
