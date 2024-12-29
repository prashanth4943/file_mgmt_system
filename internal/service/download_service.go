package service

import (
	"file_mgmt_system/internal/storage"
	"fmt"
	"io"
)

type DownloadService struct {
	db  *storage.DB
	oci *storage.OCIStorage
}

func NewDownloadService(db *storage.DB, oci *storage.OCIStorage) *DownloadService {
	return &DownloadService{
		db:  db,
		oci: oci,
	}
}

func (d *DownloadService) Download(fileID string) (string, io.Reader, error) {

	metadata, err := d.db.GetOCIFileName(fileID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch file metadata: %w", err)
	}
	fileName := metadata.FileName
	ociFileName := metadata.OCIFileName

	fileStream, err := d.oci.DownloadFile(ociFileName)
	if err != nil {
		return "", nil, err
	}

	return fileName, fileStream, nil

}
