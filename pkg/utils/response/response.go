package response

import (
	"strings"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SuccessResponse(ctx *gin.Context, statusCode int, message string, data ...interface{}) {

	response := Response{
		Success: true,
		Message: message,
		Error:   nil,
		Data:    data,
	}
	ctx.JSON(statusCode, response)
}

func ErrorResponse(ctx *gin.Context, statusCode int, message string, err error, data interface{}) {
	errFields := strings.Split(err.Error(), "\n")
	response := Response{
		Success: false,
		Message: message,
		Error:   errFields,
		Data:    data,
	}

	ctx.JSON(statusCode, response)
}
