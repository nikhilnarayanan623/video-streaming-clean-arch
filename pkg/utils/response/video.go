package response

import "time"

type VideoDetails struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" `
	Description string    `json:"description" `
	UploadedAt  time.Time `json:"uploaded_at" `
}
