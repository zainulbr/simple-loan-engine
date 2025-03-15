package filemanager

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/zainulbr/simple-loan-engine/models/filemanager"
)

type fileRepo struct {
	db *pg.DB
}

func NewFileRepository(db *pg.DB) FileRepository {
	return &fileRepo{db: db}
}

// Create a new file entry
func (r *fileRepo) Create(ctx context.Context, file *filemanager.File) error {
	data := fileModelPG{File: file}
	_, err := r.db.Model(&data).Context(ctx).Insert()
	return err
}

// Get file by ID
func (r *fileRepo) GetByID(ctx context.Context, fileID uuid.UUID) (*filemanager.File, error) {
	file := new(fileModelPG)
	err := r.db.Model(file).Context(ctx).
		Where("file_id = ?", fileID).
		Select()
	if err != nil {
		return nil, err
	}
	return file.File, nil
}

// Delete file by ID
func (r *fileRepo) Delete(ctx context.Context, fileID uuid.UUID) error {
	_, err := r.db.Model((*fileModelPG)(nil)).Context(ctx).
		Where("file_id = ?", fileID).
		Delete()
	return err
}
