package filemanager

import (
	"time"

	"github.com/google/uuid"
)

type File struct {
	FileID       uuid.UUID `json:"file_id"`
	Label        string    `json:"label"`
	Location     string    `json:"location"`
	LocationType string    `json:"location_type"`
	FileType     string    `json:"file_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
