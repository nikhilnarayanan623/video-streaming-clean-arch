package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
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
	videosDir  = "./videos"
	BufferSize = 1024 * 1024 // for 1mb
)

var (
	ErrInvalidVideoID = errors.New("invalid video id")
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

// Find one video by id
func (c *videoUseCase) FindByID(ctx context.Context, id string) (domain.Video, error) {

	video, err := c.repo.FindByID(ctx, id)
	if err != nil {
		return domain.Video{}, err
	}

	if video.ID == "" {
		return domain.Video{}, ErrInvalidVideoID
	}

	return video, nil
}

// Stream video
func (c *videoUseCase) Stream(ctx context.Context, filePath string,
	bufferChan chan []byte, sendNext chan bool, errChan chan error) {

	//open the file in the file path
	file, err := os.Open(filePath)
	if err != nil {
		errChan <- fmt.Errorf("failed to open file \nerror:%w", err)
		log.Printf("failed to open file for stream \nerror:%s", err.Error())
		return
	}
	defer file.Close()

	// create a buffer(slice of bytes) in the size of buffer size
	buffer := make([]byte, BufferSize)

	for {

		// read the file to buffer
		n, err := file.Read(buffer)

		if err != nil {
			if err == io.EOF { // if the error is end of file then send it and return
				log.Println("video streaming completed")
				errChan <- err
				return
			}
			errChan <- fmt.Errorf("failed to read file \nerror:%w", err)
			log.Printf("failed to read file \nerror:%s", err.Error())
			return
		}

		// check the data's are read or not
		if n > 0 { // if data read then send the portion of data read
			bufferChan <- buffer[:n]
		} else {
			log.Println("no values found from reading")
			errChan <- fmt.Errorf("no values read from file")
			break
		}

		// wait for the next request for send data or timeout
		select {
		case <-sendNext:
			//  continue the loop for next data send
		case <-ctx.Done():
			log.Println("streaming closed")
			errChan <- ctx.Err()
			return
		}
	}
}

// Find all videos
func (c *videoUseCase) FindAll(ctx context.Context, pagination request.Pagination) (videos []response.VideoDetails, err error) {

	video, err := c.repo.FindAll(ctx, pagination)

	return video, err
}

// create a video file on system
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

// copy the one file source to destination
func (c *videoUseCase) copyFiles(dest *os.File, src multipart.File) error {

	_, err := io.Copy(dest, src)
	return err
}
