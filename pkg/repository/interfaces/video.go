package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/response"
)

type VideoRepository interface {
	Save(ctx context.Context, video domain.Video) error
	FindAll(ctx context.Context, pagination request.Pagination) (videos []response.VideoDetails, err error)
}
