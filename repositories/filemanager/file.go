package filemanager

import (
	"context"

	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
)

type FileRepository interface {
	Create(ctx context.Context, file *filemanager.File) error
	GetByID(ctx context.Context, fileID uuid.UUID) (*filemanager.File, error)
	Delete(ctx context.Context, fileID uuid.UUID) error
}

type fileModelPG struct {
	tableName struct{} `pg:"file.files"` // Schema "file", Table "files"
	*filemanager.File
}
