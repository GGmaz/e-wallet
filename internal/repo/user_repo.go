package repo

import (
	"github.com/GGmaz/wallet-arringo/internal/db/model"
	"gorm.io/gorm"
)

func HasUserAccount(db *gorm.DB, userId int64, accNum string) bool {
	user := &model.User{}
	acc := &model.Account{}
	db.First(&user, userId)
	db.First(&acc, "acc_number = ?", accNum)
	return acc.UserId == user.ID
}
