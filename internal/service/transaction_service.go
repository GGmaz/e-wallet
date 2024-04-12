package service

import (
	"errors"
	"fmt"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TransactionServiceImpl struct {
	TransactionRepo repo.Repo[model.Transaction]
	UserRepo        repo.Repo[model.User]
	AccountRepo     repo.Repo[model.Account]
}

func (s *TransactionServiceImpl) CreateTransaction(db *gorm.DB, userId int64, amount, balance float64, txType enums.TxType, accNum string, txStatus enums.TxStatus) error {
	tx := &model.Transaction{
		UserID:          userId,
		Amount:          amount,
		Balance:         balance,
		TransactionType: txType,
		AccountNumber:   accNum,
		Status:          txStatus,
	}

	dbRes := s.TransactionRepo.Create(db, tx)
	if dbRes.Error != nil {
		return dbRes.Error
	}

	return nil
}

func (s *TransactionServiceImpl) AddMoney(c *gin.Context, userId int64, amount float64, accNum string) (float64, error) {
	db := c.MustGet("transaction").(*gorm.DB)

	if !repo.HasUserAccount(db, userId, accNum) {
		return 0, errors.New("user does not have account with provided account number")
	}

	tx := db.Begin()
	defer tx.Rollback()

	acc := &model.Account{}
	res := repo.GetAccountByNum(tx, acc, accNum)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	acc.Balance += amount

	res = s.AccountRepo.Update(tx, acc, acc.ID)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	err := s.CreateTransaction(tx, userId, amount, acc.Balance, enums.CREDIT, accNum, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return acc.Balance, nil
}

func (s *TransactionServiceImpl) Withdraw(c *gin.Context, userId int64, amount float64, accNum string) (float64, error) {
	db := c.MustGet("transaction").(*gorm.DB)

	if !repo.HasUserAccount(db, userId, accNum) {
		return 0, errors.New("user does not have account with provided account number")
	}

	tx := db.Begin()
	defer tx.Rollback()

	acc := &model.Account{}
	res := repo.GetAccountByNum(tx, acc, accNum)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	if acc.Balance < amount {
		return 0, errors.New("insufficient balance on account")
	}

	acc.Balance -= amount

	res = s.AccountRepo.Update(tx, acc, acc.ID)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	err := s.CreateTransaction(tx, userId, amount, acc.Balance, enums.DEBIT, accNum, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return acc.Balance, nil
}

func (s *TransactionServiceImpl) TransferMoney(c *gin.Context, from, to string, amount float64) error {
	db := c.MustGet("transaction").(*gorm.DB)
	tx := db.Begin()
	defer tx.Rollback()

	fromAcc := &model.Account{}
	toAcc := &model.Account{}

	res := repo.GetAccountByNum(tx, fromAcc, from)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	res = repo.GetAccountByNum(tx, toAcc, to)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if fromAcc.Balance < amount {
		return errors.New(fmt.Sprintf("insufficient balance on account %s", fromAcc.AccNumber))
	}

	if fromAcc.UserId != toAcc.UserId {
		err := s.CreateTransaction(tx, fromAcc.UserId, amount, fromAcc.Balance, enums.DEBIT, fromAcc.AccNumber, enums.RESERVED)
		if err != nil {
			tx.Rollback()
			return err
		}
		//TODO: add to redis for KYT
		return nil
	}

	err := s.CreateTransaction(tx, fromAcc.UserId, amount, fromAcc.Balance, enums.DEBIT, fromAcc.AccNumber, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return err
	}

	//TODO: to posle iz redisa sve ide (MOZE SE IZVUCI U POSEBNU FUNKCIJU)
	fromAcc.Balance -= amount
	toAcc.Balance += amount

	res = s.AccountRepo.Update(tx, fromAcc, fromAcc.ID)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	res = s.AccountRepo.Update(tx, toAcc, toAcc.ID)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	err = s.CreateTransaction(tx, toAcc.UserId, amount, toAcc.Balance, enums.CREDIT, toAcc.AccNumber, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
