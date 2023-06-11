package request

import "mime/multipart"

type UploadVideo struct {
	Name       string
	FileHeader *multipart.FileHeader
}
