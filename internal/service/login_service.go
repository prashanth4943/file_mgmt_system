package service

import (
	"file_mgmt_system/internal/models"
	"file_mgmt_system/internal/storage"
	"log"
)

type LoginService struct {
	db *storage.DB
}

func NewLoginService(db *storage.DB) *LoginService {
	return &LoginService{
		db: db,
	}
}

func (s *LoginService) Login(req *models.Input) (int, bool, error) {

	affectedRows, fileExists, err := s.db.InsertUser(req)
	if err != nil {
		return 0, false, err
	}
	log.Println(affectedRows)
	return affectedRows, fileExists, nil
}
