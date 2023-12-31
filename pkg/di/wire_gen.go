// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/db"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/repository"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase"
)

// Injectors from wire.go:

func InitializeApi(cfg *config.Config) (*api.Server, error) {
	gormDB, err := db.ConnectDatabase(cfg)
	if err != nil {
		return nil, err
	}
	videoRepository := repository.NewVideoRepository(gormDB)
	videUseCase := usecase.NewVideoUseCase(videoRepository)
	videHandler := handler.NewVideoHandler(videUseCase)
	server := api.NewServerHTTP(cfg, videHandler)
	return server, nil
}
