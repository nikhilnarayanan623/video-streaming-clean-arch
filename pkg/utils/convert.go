package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nikhilnarayanan623/video-streaming-clean-arch/pkg/utils/request"
)

const (
	defaultPageNumber = 1
	defaultPageCount  = 10
)

func GetPagination(ctx *gin.Context) request.Pagination {

	pagination := request.Pagination{
		PageNumber: defaultPageNumber,
		Count:      defaultPageCount,
	}

	num, err := strconv.ParseUint(ctx.Query("page_number"), 10, 64)
	if err == nil {
		pagination.PageNumber = num
	}

	num, err = strconv.ParseUint(ctx.Query("count"), 10, 64)
	if err == nil {
		pagination.Count = num
	}
	return pagination
}
