package api

import (
	"github.com/gin-gonic/gin"
	_ "github.com/nikhilnarayanan623/video-streaming-clean-arch/docs"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/config"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func NewServerHTTP(cfg *config.Config, videoHandler interfaces.VideHandler) *Server {

	engine := gin.Default()

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	video := engine.Group("/video")
	{
		video.POST("/", videoHandler.Upload)
		video.GET("/all", videoHandler.FindAll)
		video.GET("/stream", videoHandler.Play)
	}

	return &Server{
		engine: engine,
		port:   cfg.Port,
	}
}

func (c *Server) Start() {
	c.engine.Run(c.port)
}
