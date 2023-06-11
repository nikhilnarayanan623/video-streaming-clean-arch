package domain

import "time"

type Video struct {
	ID          string    `json:"id" gorm:"primaryKey;not null"`
	Name        string    `json:"name" gorm:"not null"`
	Url         string    `json:"url" gorm:"not null"`
	Description string    `json:"description" gorm:"not null"`
	UploadedAt  time.Time `json:"uploaded_at" gorm:"not null"`
}
