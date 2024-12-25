package service

import (
	"file_mgmt_system/internal/models"
	"file_mgmt_system/internal/storage"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"
)

type UploadService struct {
	db  *storage.DB
	oci *storage.OCIStorage
}

func NewUploadService(db *storage.DB, oci *storage.OCIStorage) *UploadService {
	return &UploadService{
		db:  db,
		oci: oci,
	}
}
func (u *UploadService) Upload(file multipart.File, header *multipart.FileHeader, req models.UploadRequest) (models.FileMetadata, error) {
	FileID, uniqueFileName := generateUniqueName(header.Filename)
	log.Println(FileID)
	ociReference, err := u.oci.UploadFile(uniqueFileName, file, header.Size)
	if err != nil {
		return models.FileMetadata{}, err
	}
	fileMetadata := models.FileMetadata{
		FileName:     req.FileName,
		UniqueName:   uniqueFileName,
		FileType:     req.FileType,
		FileSize:     req.FileSize,
		Email:        req.Email,
		UploadTime:   time.Now(),
		OCIReference: ociReference,
		FileID:       FileID,
	}

	// Save the metadata in the database
	err = u.db.SaveFileMetadata(fileMetadata)
	if err != nil {
		return models.FileMetadata{}, err
	}
	return fileMetadata, nil
}

func generateUniqueName(originalName string) (string, string) {
	// Extract the file extension
	parts := strings.Split(originalName, ".")
	extension := ""
	if len(parts) > 1 {
		extension = parts[len(parts)-1]
	}
	FileID := uuid.New().String()

	// Generate UUID and timestamp-based unique name
	uniqueName := fmt.Sprintf("%s_%d.%s", FileID, time.Now().Unix(), extension)
	return FileID, uniqueName
}
