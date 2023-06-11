package main

import (
	"log"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/di"
)

func main() {

	cfg, err := config.LoadConfigs()
	if err != nil {
		log.Fatalf("failed load configs \nerror:%w", err)
	}

	server, err := di.InitializeApi(cfg)
	if err != nil {
		log.Fatalf("failed initialize api \nerror:%w", err)
	}

	server.Start()
}
