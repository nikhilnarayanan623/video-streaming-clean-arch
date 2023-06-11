package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler/interfaces"
	usecase "github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase/interfaces"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/response"
)

type videHandler struct {
	usecase usecase.VideUseCase
}

func NewVideoHandler(usecase usecase.VideUseCase) interfaces.VideHandler {
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
func (c *videHandler) Stream(ctx *gin.Context) {

}
