package usecase

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/domain"
	repo "github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/repository/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
)

type videoUseCase struct {
	repo repo.VideoRepository
}

const (
	videosDir   = "./videos"
	playlistDir = "./playlists"
)

func NewVideoUseCase(repo repo.VideoRepository) interfaces.VideUseCase {
	return &videoUseCase{
		repo: repo,
	}
}

// Save Video
func (c *videoUseCase) Save(ctx context.Context, uploadReq request.UploadVideo) (string, error) {
	// open the fileHeader to get file
	file, err := uploadReq.FileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed open file \nerror:%w", err)
	}

	videID := utils.NewVideUniqueID()
	//create a new file with unique id
	destFile, err := c.createNewFile(videID)
	if err != nil {
		return "", fmt.Errorf("failed to create new file \nerror:%w", err)
	}

	// run copy function and database save function on two go routine
	// two functions are independents
	errChan := make(chan error, 2)

	// copy file from uploaded to new created file
	go func(errChan chan error) {
		err = c.copyFiles(destFile, file)
		if err != nil {
			errChan <- fmt.Errorf("failed to copy given file \nerror:%w", err)
		}
		errChan <- nil
	}(errChan)

	// save file details on database
	go func(errChan chan error) {
		err = c.repo.Save(ctx, domain.Video{
			ID:   videID,
			Name: uploadReq.Name,
		})
		if err != nil {
			errChan <- fmt.Errorf("failed to save file details on database \nerror:%w", err)
		}
		errChan <- nil
	}(errChan)

	// wait for the two go routines
	for i := 1; i <= 2; i++ {
		err = <-errChan
		if err != nil {
			return "", err
		}
	}

	return videID, nil
}
func (c *videoUseCase) FindAll(ctx context.Context, pagination request.Pagination) ([]domain.Video, error) {

	video, err := c.repo.FindAll(ctx, pagination)

	return video, err
}

func (c *videoUseCase) createNewFile(fileID string) (*os.File, error) {

	err := os.MkdirAll(videosDir, 0700)
	if err != nil {
		return nil, err
	}

	fileName := videosDir + "/" + fileID + ".mp4"

	file, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (c *videoUseCase) copyFiles(dest *os.File, src multipart.File) error {
	
	_, err := io.Copy(dest, src)
	return err
}
