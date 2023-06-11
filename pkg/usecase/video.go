package usecase

import (
	"context"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	repo "github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
)

type videoUseCase struct {
	repo repo.VideoRepository
}

func NewVideoUseCase(repo repo.VideoRepository) interfaces.VideUseCase {
	return &videoUseCase{
		repo: repo,
	}
}

func (c *videoUseCase) Save(ctx context.Context, video domain.Video) (string, error) {

	video.ID = utils.NewVideUniqueID()
	err := c.repo.Save(ctx, video)

	return video.ID, err
}
func (c *videoUseCase) FindAll(ctx context.Context, pagination request.Pagination) ([]domain.Video, error) {

	video, err := c.repo.FindAll(ctx, pagination)

	return video, err
}
