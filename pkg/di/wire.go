//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/config"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/db"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/repository"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase"
)

func InitializeApi(cfg *config.Config) (*api.Server, error) {

	wire.Build(
		db.ConnectDatabase,
		repository.NewVideoRepository,
		usecase.NewVideoUseCase,
		handler.NewVideoHandler,
		api.NewServerHTTP,
	)
	return &api.Server{}, nil
}
