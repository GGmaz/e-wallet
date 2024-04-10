package service

import (
	"errors"
	"fmt"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/internal/repo"
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

type TransactionServiceImpl struct {
	TransactionRepo repo.Repo[model.Transaction]
	UserRepo        repo.Repo[model.User]
	Conn            gorm.DB
}

func (s *TransactionServiceImpl) CreateTransaction(db *gorm.DB, userId int64, amount, balance float64, txType model.TxType) error {
	tx := &model.Transaction{
		UserID:          userId,
		Amount:          amount,
		Balance:         balance,
		TransactionType: txType,
	}

	dbRes := s.TransactionRepo.Create(db, tx)
	if dbRes.Error != nil {
		return dbRes.Error
	}

	return nil
}

func (s *TransactionServiceImpl) AddMoney(c *gin.Context, id int64, amount float64) (float64, error) {
	db := c.MustGet("transaction").(*gorm.DB)
	tx := db.Begin()
	defer tx.Rollback()

	user := &model.User{}

	res := s.UserRepo.GetById(tx, user, strconv.FormatInt(id, 10))
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	user.Balance += amount

	res = s.UserRepo.Update(tx, user, user.ID)
	if res.Error != nil {
		tx.Rollback()
		return 0, res.Error
	}

	err := s.CreateTransaction(tx, user.ID, amount, user.Balance, enums.CREDIT)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()

	return user.Balance, nil
}

func (s *TransactionServiceImpl) TransferMoney(c *gin.Context, from int64, to int64, amount float64) error {
	db := c.MustGet("transaction").(*gorm.DB)
	tx := db.Begin()
	defer tx.Rollback()

	fromUser := &model.User{}
	toUser := &model.User{}

	res := s.UserRepo.GetById(tx, fromUser, strconv.FormatInt(from, 10))
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	res = s.UserRepo.GetById(tx, toUser, strconv.FormatInt(to, 10))
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	if fromUser.Balance < amount {
		return errors.New(fmt.Sprintf("insufficient balance on user %d account", fromUser.ID))
	}

	fromUser.Balance -= amount
	toUser.Balance += amount

	res = s.UserRepo.Update(tx, fromUser, fromUser.ID)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	res = s.UserRepo.Update(tx, toUser, toUser.ID)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	err := s.CreateTransaction(tx, fromUser.ID, amount, fromUser.Balance, enums.DEBIT)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = s.CreateTransaction(tx, toUser.ID, amount, toUser.Balance, enums.CREDIT)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}
