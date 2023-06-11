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
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/response"
)

type videoUseCase struct {
	repo repo.VideoRepository
}

const (
	videosDir = "./videos"
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

	videoID := utils.NewVideUniqueID()
	//create a new file with unique id
	destFile, filePath, err := c.createNewVideoFile(videoID)
	if err != nil {
		return "", fmt.Errorf("failed to create new file \nerror:%w", err)
	}

	// run copy function and database save function on two go routine
	// two functions are independent
	errChan := make(chan error, 2)

	// copy file from uploaded to new created file on one go routine
	go func(errChan chan error) {
		err = c.copyFiles(destFile, file)
		if err != nil {
			errChan <- fmt.Errorf("failed to copy given file \nerror:%w", err)
		}
		errChan <- nil
	}(errChan)

	// save file details on database on another routine
	go func(errChan chan error) {
		err = c.repo.Save(ctx, domain.Video{
			ID:          videoID,
			Name:        uploadReq.Name,
			Url:         filePath,
			Description: uploadReq.Description,
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

	return videoID, nil
}
func (c *videoUseCase) FindAll(ctx context.Context, pagination request.Pagination) (videos []response.VideoDetails, err error) {

	video, err := c.repo.FindAll(ctx, pagination)

	return video, err
}

func (c *videoUseCase) createNewVideoFile(fileID string) (file *os.File, filePath string, err error) {

	err = os.MkdirAll(videosDir, 0700)
	if err != nil {
		return nil, "", err
	}

	filePath = videosDir + "/" + fileID + ".mp4"

	file, err = os.Create(filePath)
	if err != nil {
		return nil, "", err
	}

	return file, filePath, nil
}

func (c *videoUseCase) copyFiles(dest *os.File, src multipart.File) error {

	_, err := io.Copy(dest, src)
	return err
}
