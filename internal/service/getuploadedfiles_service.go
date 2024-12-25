package service

import (
	"file_mgmt_system/internal/models"
	"file_mgmt_system/internal/storage"
)

type GetFilesService struct {
	db *storage.DB
}

func NewGetFilesService(db *storage.DB) *GetFilesService {
	return &GetFilesService{
		db: db,
	}
}

func (service *GetFilesService) GetUploadedFilesByEmail(email string) ([]models.FileMetadata, error) {
	files, err := service.db.GetFileList(email)
	if err != nil {
		return nil, err
	}

	// Business logic could be added here if needed
	return files, nil
}
