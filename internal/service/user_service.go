package service

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	UserRepo repo.Repo[model.User]
}

func (s *UserServiceImpl) CreateUser(db *gorm.DB, email string) (*model.User, error) {
	user := &model.User{
		Email: email,
	}
	dbRes := s.UserRepo.Create(db, user)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return user, nil
}
