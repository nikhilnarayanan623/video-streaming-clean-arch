package utils

import "github.com/google/uuid"

func NewVideUniqueID() string {
	return uuid.New().String()
}
