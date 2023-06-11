package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/api/handler/interfaces"
	usecase "github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/usecase/interfaces"
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

// Upload Video
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

}
func (c *videHandler) Play(ctx *gin.Context) {

}
