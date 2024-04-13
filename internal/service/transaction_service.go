package service

import (
	"errors"
	"fmt"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type TransactionServiceImpl struct {
	TransactionRepo repo.Repo[model.Transaction]
	UserRepo        repo.Repo[model.User]
	AccountRepo     repo.Repo[model.Account]
}

func (s *TransactionServiceImpl) CreateTransaction(db *gorm.DB, userId int64, amount, balance float64, txType enums.TxType, accNum string, txStatus enums.TxStatus) (*model.Transaction, error) {
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
		return nil, dbRes.Error
	}

	return tx, nil
}

func (s *TransactionServiceImpl) AddMoney(c *gin.Context, userId int64, amount float64, accNum string) (float64, error) {
	db := c.MustGet("transaction").(*gorm.DB)

	if !repo.HasUserAccount(db, userId, accNum) || !repo.IsAccountVerified(db, accNum) {
		return 0, errors.New("user does not have account with provided account number or account is not verified")
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

	_, err := s.CreateTransaction(tx, userId, amount, acc.Balance, enums.CREDIT, accNum, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return acc.Balance, nil
}

func (s *TransactionServiceImpl) Withdraw(c *gin.Context, userId int64, amount float64, accNum string) (float64, error) {
	db := c.MustGet("transaction").(*gorm.DB)

	if !repo.HasUserAccount(db, userId, accNum) || !repo.IsAccountVerified(db, accNum) {
		return 0, errors.New("user does not have account with provided account number or account is not verified")
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
		tx.Rollback()
		return 0, errors.New("insufficient balance on account")
	}

	acc.Balance -= amount

	res = s.AccountRepo.Update(tx, acc, acc.ID)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	_, err := s.CreateTransaction(tx, userId, amount, acc.Balance, enums.DEBIT, accNum, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return acc.Balance, nil
}

func (s *TransactionServiceImpl) TransferMoney(c *gin.Context, from, to string, amount float64) error {
	db := c.MustGet("transaction").(*gorm.DB)
	r := c.MustGet("redis").(*redis.Client)

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

	if fromAcc.Status != enums.VERIFIED || toAcc.Status != enums.VERIFIED {
		tx.Rollback()
		return errors.New("account is not verified")
	}

	if fromAcc.Balance < amount {
		tx.Rollback()
		return errors.New(fmt.Sprintf("insufficient balance on account %s", fromAcc.AccNumber))
	}

	if fromAcc.UserId != toAcc.UserId {
		transaction, err := s.CreateTransaction(tx, fromAcc.UserId, amount, fromAcc.Balance, enums.DEBIT, fromAcc.AccNumber, enums.RESERVED)
		if err != nil {
			tx.Rollback()
			return err
		}

		// KYT - Insert transaction ID into Redis using transaction ID as key and current time as value
		currentTime := time.Now().Unix()
		err = r.Set(c, "transactionId:"+strconv.FormatInt(transaction.ID, 10), currentTime, 0).Err()
		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()
		return nil
	}

	_, err := s.CreateTransaction(tx, fromAcc.UserId, amount, fromAcc.Balance, enums.DEBIT, fromAcc.AccNumber, enums.SUCCESS)
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

	_, err = s.CreateTransaction(tx, toAcc.UserId, amount, toAcc.Balance, enums.CREDIT, toAcc.AccNumber, enums.SUCCESS)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (s *TransactionServiceImpl) GetTransactionById(c *gin.Context, txId string) (*model.Transaction, error) {
	// Convert txId to int64
	txIdInt, err := strconv.ParseInt(txId, 10, 64)
	if err != nil {
		return nil, err
	}

	db := c.MustGet("transaction").(*gorm.DB)
	tx := &model.Transaction{}

	res := s.TransactionRepo.GetById(db, tx, txIdInt)
	if res.Error != nil {
		return nil, res.Error
	}

	return tx, nil
}
