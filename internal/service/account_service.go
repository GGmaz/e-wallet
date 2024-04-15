package service

import (
	"context"
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"github.com/GGmaz/wallet-arringo/pkg/enums"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"log"
	"time"
)

// Define an interface for the repository
type Repository[T any] interface {
	Create(db *gorm.DB, t *T) *gorm.DB
	GetByField(db *gorm.DB, t *T, field string, fieldVal string, preload ...string) *gorm.DB
	Update(db *gorm.DB, t *T, id int64) *gorm.DB
	GetById(db *gorm.DB, t *T, id int64, preload ...string) *gorm.DB
}

type AccountServiceImpl struct {
	AccountRepo Repository[model.Account]
	UserRepo    Repository[model.User]
}

func (s *AccountServiceImpl) CreateAccount(c *gin.Context, accNum string, userId int64) (*model.Account, error) {
	db := c.MustGet("transaction").(*gorm.DB)
	r := c.MustGet("redis").(*redis.Client)

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

	// KYC - Insert account ID into Redis using account ID as key and current time as value
	currentTime := time.Now().Unix()
	err := r.Set(context.Background(), "accountNum:"+acc.AccNumber, currentTime, 0).Err()
	if err != nil {
		return nil, err
	}

	log.Println("Account was sent for verification")
	return acc, nil
}

func (s *AccountServiceImpl) VerifyAccount(c *gin.Context, accNum string) error {
	db := c.MustGet("transaction").(*gorm.DB)
	acc := &model.Account{}

	dbRes := s.AccountRepo.GetByField(db, acc, "acc_number", accNum)
	if dbRes.Error != nil {
		return dbRes.Error
	}

	acc.Status = enums.VERIFIED

	dbRes = s.AccountRepo.Update(db, acc, acc.ID)
	if dbRes.Error != nil {
		return dbRes.Error
	}

	return nil
}
