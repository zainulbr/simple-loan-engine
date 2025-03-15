package filemanager

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
	repoFile "github.com/zainulbr/simple-loan-engine/repositories/filemanager"
)

type fileService struct {
	fileRepo repoFile.FileRepository
	basePath string // Directory untuk menyimpan file
}

// NewFileService creates a new instance of FileService
func NewFileService(repo repoFile.FileRepository, basePath string) FileService {
	return &fileService{fileRepo: repo, basePath: basePath}
}

// Validate file extension and MIME type
func (s *fileService) ValidateFileFormat(file multipart.File, fileHeader *multipart.FileHeader) error {
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))

	// Check if the extension is allowed
	mimeType, allowed := allowedExtensions[ext]
	if !allowed {
		return errors.New("invalid file type: only JPG, PNG, PDF are allowed")
	}

	// Validate MIME type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return errors.New("failed to read file header")
	}

	// Reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return errors.New("failed to reset file pointer")
	}

	detectedMimeType := http.DetectContentType(buffer)
	if detectedMimeType != mimeType {
		return fmt.Errorf("file type mismatch: expected %s, got %s", mimeType, detectedMimeType)
	}

	return nil
}

// UploadFile handles file upload and storage
func (s *fileService) UploadFile(ctx context.Context, file multipart.File,
	fileHeader *multipart.FileHeader,
	locationType filemanager.LocationType) (*filemanager.File, error) {

	// Generate unique filename
	ext := filepath.Ext(fileHeader.Filename)
	newFileName := uuid.New().String() + ext
	filePath := filepath.Join(s.basePath, newFileName)

	// Save file to disk
	outFile, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}
	defer outFile.Close()

	_, err = outFile.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// Create file entry in DB
	newFile := &filemanager.File{
		FileID:       uuid.New(),
		Label:        fileHeader.Filename,
		Location:     filePath,
		FileType:     ext,
		LocationType: filemanager.LocationTypeLocal,
	}

	err = s.fileRepo.Create(ctx, newFile)
	if err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}

	return newFile, nil
}

// GetFileByID retrieves a file metadata by ID
func (s *fileService) GetFileByID(ctx context.Context, fileID uuid.UUID) (*filemanager.File, error) {
	return s.fileRepo.GetByID(ctx, fileID)
}

// DeleteFile deletes a file entry and removes the actual file
func (s *fileService) DeleteFile(ctx context.Context, fileID uuid.UUID) error {
	// Fetch file info
	file, err := s.fileRepo.GetByID(ctx, fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}

	// Remove file from disk
	if err := os.Remove(file.Location); err != nil {
		return fmt.Errorf("failed to delete file from disk: %w", err)
	}

	// Delete from DB
	return s.fileRepo.Delete(ctx, fileID)
}
