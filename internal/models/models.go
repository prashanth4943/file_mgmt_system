package models

import "time"

type Input struct {
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type UploadRequest struct {
	FileName string `json:"fileName"`
	FileType string `json:"fileType"`
	FileSize int64  `json:"fileSize"`
	Email    string `json:"email"`
}
type FileMetadata struct {
	FileName     string
	UniqueName   string
	FileType     string
	FileSize     int64
	Email        string
	UploadTime   time.Time
	OCIReference string
	FileID       string
}

type FileName struct {
	FileName    string
	OCIFileName string
}
