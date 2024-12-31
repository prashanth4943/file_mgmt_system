package service

import (
	"bytes"
	"file_mgmt_system/internal/storage"
	"fmt"
	"io"
	"log"
)

type DownloadService struct {
	db          *storage.DB
	oci         *storage.OCIStorage
	redisClient *storage.RedisClient
}

func NewDownloadService(db *storage.DB, oci *storage.OCIStorage, redisClient *storage.RedisClient) *DownloadService {
	return &DownloadService{
		db:          db,
		oci:         oci,
		redisClient: redisClient,
	}
}

func (d *DownloadService) Download(fileID string) (string, io.Reader, error) {

	metadata, err := d.db.GetOCIFileName(fileID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to fetch file metadata: %w", err)
	}
	fileName := metadata.FileName
	ociFileName := metadata.OCIFileName

	fileBytes, err := d.redisClient.GetFile(fileID)
	if err == nil && len(fileBytes) > 0 {
		log.Printf("File %s fetched from Redis", fileID)
		return fileName, bytes.NewReader(fileBytes), nil
	}

	fileStream, err := d.oci.DownloadFile(ociFileName)
	if err != nil {
		return "", nil, err
	}

	fileBytes, err = io.ReadAll(fileStream)
	if err != nil {
		return "", nil, fmt.Errorf("failed to read file stream: %w", err)
	}

	err = d.redisClient.SetFile(fileID, fileBytes)
	if err != nil {
		log.Printf("Failed to cache file %s in Redis: %v", fileID, err)
	}

	return fileName, fileStream, nil

}
