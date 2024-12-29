package service

import (
	"file_mgmt_system/internal/storage"
)

type DeleteService struct {
	db  *storage.DB
	oci *storage.OCIStorage
}

func NewDeleteService(db *storage.DB, oci *storage.OCIStorage) *DeleteService {
	return &DeleteService{
		db:  db,
		oci: oci,
	}
}

func (d *DeleteService) Delete(fileID string) (string, error) {

	fileName, uniqueName, err := d.db.DeleteFile(fileID)
	if err != nil {
		return "", err
	}
	err = d.oci.DeleteFile(uniqueName)
	if err != nil {
		return "", err
	}

	return fileName, nil

}
