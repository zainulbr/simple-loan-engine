package filemanager

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
)

// FileService defines the file use case interface
type FileService interface {
	UploadFile(ctx context.Context, file multipart.File, fileHeader *multipart.FileHeader, locationType filemanager.LocationType) (*filemanager.File, error)
	GetFileByID(ctx context.Context, fileID uuid.UUID) (*filemanager.File, error)
	DeleteFile(ctx context.Context, fileID uuid.UUID) error
	ValidateFileFormat(file multipart.File, fileHeader *multipart.FileHeader) error
}

// Allowed file extensions and MIME types
var allowedExtensions = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".pdf":  "application/pdf",
}
