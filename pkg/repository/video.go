package repository

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
	"gorm.io/gorm"
)

type videDatabase struct {
	db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) interfaces.VideoRepository {
	return &videDatabase{
		db: db,
	}
}

func (c *videDatabase) Save(ctx context.Context, video domain.Video) error {

	uploadedAt := time.Now()
	query := `INSERT INTO videos (id, name, uploaded_at) VALUES ($1, $2, $3)`
	err := c.db.Exec(query, video.ID, video.Name, uploadedAt).Error

	return err
}

func (c *videDatabase) FindAll(ctx context.Context, pagination request.Pagination) (videos []domain.Video, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT id, name, uploaded_at FROM videos LIMIT $1 OFFSET $2`
	err = c.db.Raw(query, limit, offset).Scan(&videos).Error

	return
}
