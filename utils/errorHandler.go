package utils

import (
	"log"

	"github.com/gin-gonic/gin"
)

// CustomError ,Client Error Message
type CustomError struct {
	Msg     string
	Command string
}

// PushError ,Simple Error Handler
func PushError(code int, clientMes CustomError, err error, ctx *gin.Context) {
	// Get Log To Console
	log.Printf("Error:\n\n%v\n\nStackTrace:\n\n%v", clientMes, err)

	// Back To Client
	ctx.JSON(code, struct {
		Msg string `json:"msg"`
	}{
		Msg: clientMes.Msg,
	})
}
