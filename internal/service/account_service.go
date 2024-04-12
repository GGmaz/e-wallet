package service

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"gorm.io/gorm"
)

type AccountServiceImpl struct {
	AccountRepo repo.Repo[model.Account]
	UserRepo    repo.Repo[model.User]
}

func (s *AccountServiceImpl) CreateAccount(db *gorm.DB, accNum string, userId int64) (*model.Account, error) {
	//TODO: mozda ubaciti sve u sklopu jedne transakcije
	user := &model.User{}

	dbRes := s.UserRepo.GetById(db, user, userId)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	acc := &model.Account{
		AccNumber: accNum,
		Balance:   0,
		Status:    enums.AccountStatus(enums.UNVERIFIED),
		UserId:    userId,
	}

	dbRes = s.AccountRepo.Create(db, acc)
	if dbRes.Error != nil {
		return nil, dbRes.Error
	}

	//TODO: insert in redis for KYC

	return acc, nil
}
