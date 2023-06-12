package repository

import (
	"context"
	"time"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/response"
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
	query := `INSERT INTO videos (id, name, url, description, uploaded_at) VALUES ($1, $2, $3, $4, $5)`
	err := c.db.Exec(query, video.ID, video.Name, video.Url, video.Description, uploadedAt).Error

	return err
}

func (c *videDatabase) FindByID(ctx context.Context, id string) (video domain.Video, err error) {

	query := `SELECT id, name, url, uploaded_at FROM videos WHERE id = $1`
	err = c.db.Raw(query, id).Scan(&video).Error

	return
}

func (c *videDatabase) FindAll(ctx context.Context, pagination request.Pagination) (videos []response.VideoDetails, err error) {

	limit := pagination.Count
	offset := (pagination.PageNumber - 1) * limit

	query := `SELECT id, name, description, uploaded_at FROM videos LIMIT $1 OFFSET $2`
	err = c.db.Raw(query, limit, offset).Scan(&videos).Error

	return
}
