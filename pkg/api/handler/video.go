package handler

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase"
	usecaseInterface "github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/response"
)

type videHandler struct {
	usecase usecaseInterface.VideUseCase
}

func NewVideoHandler(usecase usecaseInterface.VideUseCase) interfaces.VideHandler {
	return &videHandler{
		usecase: usecase,
	}
}

// Upload godoc
// @summary api for upload videos to server
// @tags Video
// @id Upload
// @Param     video   formData     file   true   "Video file to upload"
// @Param     name   formData     string   true   "Video Name"
// @Param     description   formData     string   true   "Video Description"
// @Router /video [post]
// @Success 201 {object} response.Response{} "successfully video uploaded"
// @Failure 400 {object} response.Response{}  "failed get inputs"
// @Failure 500 {object} response.Response{}  "failed to save video"
func (c *videHandler) Upload(ctx *gin.Context) {
	// user given video name
	videoName := ctx.Request.PostFormValue("name")
	videoDescription := ctx.Request.PostFormValue("description")
	if videoName == "" || videoDescription == "" {
		err := errors.New("failed get inputs 'name' or 'description' ")
		response.ErrorResponse(ctx, http.StatusBadRequest, "failed to get inputs", err, nil)
		return
	}
	// get video file from request
	fileHeader, err := ctx.FormFile("video")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "failed to get video from request", err, nil)
		return
	}

	uploadVideo := request.UploadVideo{
		Name:        videoName,
		FileHeader:  fileHeader,
		Description: videoDescription,
	}

	videoID, err := c.usecase.Save(ctx, uploadVideo)

	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save video", err, nil)
		return
	}

	response.SuccessResponse(ctx, http.StatusCreated, "successfully video uploaded", gin.H{
		"video_id": videoID,
	})
}

// FindAll godoc
// @summary api for find all videos on server
// @tags Video
// @id FindAll
// @Param     page_number   query     string   false   "Page Number"
// @Param     count   query     string   false   "Count"
// @Router /video/all [get]
// @Success 201 {object} response.Response{} "successfully found all videos"
// @Failure 500 {object} response.Response{}  "failed to get all videos"
func (c *videHandler) FindAll(ctx *gin.Context) {

	pagination := utils.GetPagination(ctx)

	videos, err := c.usecase.FindAll(ctx, pagination)
	if err != nil {
		response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to get all videos", err, nil)
		return
	}

	if videos == nil {
		response.SuccessResponse(ctx, http.StatusOK, "there is no videos to show")
		return
	}

	response.SuccessResponse(ctx, http.StatusOK, "successfully found all videos", videos)
}

// FindAll godoc
// @summary api for stream video through a single tcp connection
// @tags Video
// @id Stream
// @Param     video_id   path     string   true   "video ID"
// @Router /video/stream/{video_id} [get]
// @Failure 500 {object} response.Response{}  "failed to stream video"
func (c *videHandler) Stream(ctx *gin.Context) {

	videoID := ctx.Param("video_id")

	video, err := c.usecase.FindByID(ctx, videoID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == usecase.ErrInvalidVideoID {
			statusCode = http.StatusBadRequest
		}
		response.ErrorResponse(ctx, statusCode, "failed to get video", err, nil)
		return
	}

	ctx.Header("Content-Type", "video/mp4")
	ctx.Header("Transfer-Encoding", "chunked")
	// ctx.Header("Content-Length", strconv.FormatInt(fileSize, 10))

	fCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	bufferChan := make(chan []byte)
	sendNext := make(chan bool)
	errChan := make(chan error)

	// call the usecase in other go routine to read data concurrently
	go c.usecase.Stream(fCtx, video.Url, bufferChan, sendNext, errChan)

	for {
		select {
		// if data came then write it on writer
		case data := <-bufferChan:
			_, err := ctx.Writer.Write(data)
			if err != nil {
				return
			}
			ctx.Writer.Flush()
			// sendNext to true for read the next part
			sendNext <- true

			// catch error through error channel
		case err = <-errChan:
			// if error is not the end of video then response error and break
			if err != io.EOF {
				response.ErrorResponse(ctx, http.StatusInternalServerError, "failed to stream video", err, nil)
			}
			return
		case <-ctx.Done():
			return
		}
	}

}
