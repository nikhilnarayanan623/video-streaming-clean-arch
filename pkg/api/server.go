package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/config"
)

type Server struct {
	engine *gin.Engine
	port   string
}

func NewServerHTTP(cfg *config.Config, videHandler interfaces.VideHandler) *Server {

	engine := gin.Default()

	video := engine.Group("/video")
	{
		video.POST("/", videHandler.Upload)
		video.GET("/all", videHandler.FindAll)
		video.GET("/play", videHandler.Play)
	}

	return &Server{
		engine: engine,
		port:   cfg.Port,
	}
}

func (c *Server) Start() {
	c.engine.Run(c.port)
}
