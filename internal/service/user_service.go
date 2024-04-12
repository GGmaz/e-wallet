package service

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/pkg/models"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	UserRepo repo.Repo[model.User]
}

func (s *UserServiceImpl) CreateUser(db *gorm.DB, email, firstName, lastName, address string) (*model.User, error) {
	user := &model.User{
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
	}
	dbRes := s.UserRepo.Create(db, user)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	return user, nil
}

func (s *UserServiceImpl) GetAccounts(db *gorm.DB, email string) (*models.AccDataRes, error) {
	dbRes, err := repo.GetUserByMail(db, email)
	if err != nil {
		return nil, err
	}

	accData := &models.AccDataRes{
		Accounts: []model.Account{},
	}

	if dbRes.Accounts != nil {
		accData.Accounts = dbRes.Accounts
	}

	return accData, nil
}
