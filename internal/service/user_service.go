package service

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserServiceImpl struct {
	UserRepo repo.Repo[model.User]
}

func (s *UserServiceImpl) CreateUser(c *gin.Context, email, firstName, lastName, address string) (*model.User, error) {
	db := c.MustGet("transaction").(*gorm.DB)

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

func (s *UserServiceImpl) GetAccounts(c *gin.Context, email string) (*models.AccDataRes, error) {
	db := c.MustGet("transaction").(*gorm.DB)
	user := &model.User{}

	dbRes := s.UserRepo.GetByField(db, user, "email", email, "Accounts")
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	accData := &models.AccDataRes{
		Accounts: []model.Account{},
	}

	if user.Accounts != nil {
		accData.Accounts = user.Accounts
	}

	return accData, nil
}
