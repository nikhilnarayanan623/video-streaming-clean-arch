package interfaces

import "github.com/gin-gonic/gin"

type VideHandler interface {
	Upload(ctx *gin.Context)
	FindAll(ctx *gin.Context)
	Play(ctx *gin.Context)
}
