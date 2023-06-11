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
// @Router /video [post]
// @Success 201 {object} response.Response{} "successfully video uploaded"
// @Failure 400 {object} response.Response{}  "failed get inputs"
// @Failure 500 {object} response.Response{}  "failed to save video"
func (c *videHandler) Upload(ctx *gin.Context) {
	// get video file from request
	fileHeader, err := ctx.FormFile("video")
	if err != nil {
		response.ErrorResponse(ctx, http.StatusBadRequest, "failed to get video from request", err, nil)
		return
	}
	// user given video name
	videoName := ctx.Request.PostFormValue("name")
	if videoName == "" {
		response.ErrorResponse(ctx, http.StatusBadRequest, "failed to get vide name", errors.New("name not provided"), nil)
		return
	}

	uploadVideo := request.UploadVideo{
		Name:       videoName,
		FileHeader: fileHeader,
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
func (c *videHandler) Play(ctx *gin.Context) {

}
