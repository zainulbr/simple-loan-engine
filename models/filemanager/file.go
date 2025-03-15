package filemanager

import (
	"time"

	"github.com/google/uuid"
)

type LocationType string

const LocationTypeS3 LocationType = "s3"
const LocationTypeLocal LocationType = "local"

type File struct {
	FileID       uuid.UUID    `json:"file_id"`
	Label        string       `json:"label"`
	Location     string       `json:"location"`
	LocationType LocationType `json:"location_type"`
	FileType     string       `json:"file_type"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}
