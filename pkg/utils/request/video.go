package request

import "mime/multipart"

type UploadVideo struct {
	Name        string
	Description string
	FileHeader  *multipart.FileHeader
}
