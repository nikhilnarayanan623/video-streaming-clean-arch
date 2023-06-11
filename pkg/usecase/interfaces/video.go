package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
)

type VideUseCase interface {
	Save(ctx context.Context, video domain.Video) (videoID string, err error)
	FindAll(ctx context.Context, pagination request.Pagination) ([]domain.Video, error)
}
