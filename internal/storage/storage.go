package storage

import "io"

// Storage defines the methods for file operations
type Storage interface {
	UploadFile(objectName string, content io.Reader, contentLength int64) error
	DownloadFile(objectName string) (io.Reader, error)
	DeleteFile(objectName string) error
}
