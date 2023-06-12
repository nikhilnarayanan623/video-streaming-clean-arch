package interfaces

import (
	"context"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/response"
)

type VideUseCase interface {
	Save(ctx context.Context, uploadReq request.UploadVideo) (videoID string, err error)
	FindByID(ctx context.Context, id string) (domain.Video, error)
	FindAll(ctx context.Context, pagination request.Pagination) (videos []response.VideoDetails, err error)

	Stream(ctx context.Context, filePath string, buffer chan []byte, flagChan chan bool, errChan chan error)
}
