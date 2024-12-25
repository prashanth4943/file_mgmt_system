package service

import (
	"file_mgmt_system/internal/models"
	"file_mgmt_system/internal/storage"
)

type LoginService struct {
	db *storage.DB
}

func NewLoginService(db *storage.DB) *LoginService {
	return &LoginService{
		db: db,
	}
}

func (s *LoginService) Login(req *models.Input) (bool, error) {

	affectedRows, err := s.db.InsertUser(req)
	if err != nil {
		return false, err
	}
	return affectedRows > 0, nil
}
